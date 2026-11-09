package dnslayer

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"sync"

	"github.com/BlackDark/travo-greenfield/internal/atomicfile"
	"github.com/BlackDark/travo-greenfield/internal/device"
	"github.com/BlackDark/travo-greenfield/internal/guard"
)

const guardName = "dns-layer-in-progress"

// Applier stages under the config locks, commits and reloads (txn.Manager.Apply).
// Nothing here commits directly; staging inside the lock closes the foreign-writer race.
type Applier interface {
	Apply(ctx context.Context, configs []string, stage func(ctx context.Context) error) error
}

// Service is the only writer of the DNS options the stack owns.
type Service struct {
	Device     device.Device
	Guards     guard.Dir
	RecordPath string
	Apply      Applier
	mu         sync.Mutex
}

// Load reads the record. Absent = no layers. Corrupt = error (fail closed, never "empty").
func (s *Service) Load() (*Record, error) {
	var r Record
	err := atomicfile.ReadJSON(s.RecordPath, &r)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if err := Check(r); err != nil {
		return nil, err
	}
	return &r, nil
}

// ErrForeignStaged means dhcp has uncommitted changes from someone else; staging on top
// would commit or revert them. The transaction module's single staging area prevents this.
var ErrForeignStaged = errors.New("dhcp has foreign staged changes")

// Push activates a layer. The record (intent) is saved before staging, so after a crash
// or failed apply the guard stays and Resume re-applies the same record.
func (s *Service) Push(ctx context.Context, l Layer) error {
	return s.mutate(ctx, func(rec *Record, cur State) (*Record, error) {
		out, err := Push(rec, cur, l)
		return &out, err
	})
}

// Pop deactivates a layer. The last pop restores base, and only after that
// succeeds is the record deleted.
func (s *Service) Pop(ctx context.Context, name Name) error {
	return s.mutate(ctx, func(rec *Record, _ State) (*Record, error) { return Pop(rec, name) })
}

// Status is the read-only view. Pending means a guarded change did not finish:
// the record is intent, not proof the device has it.
type Status struct {
	Layers  []Name `json:"layers"`
	Pending bool   `json:"pending"`
	// Corrupt is set when the record cannot be read; Pending still reflects the guard.
	Corrupt string `json:"corrupt,omitempty"`
}

// Status never mutates. A corrupt record is reported in the value, not as an error,
// so the UI can still show "pending" when it matters most.
func (s *Service) Status() (Status, error) {
	st := Status{Layers: []Name{}}
	g, err := s.Guards.Stale()
	if err != nil {
		return st, err
	}
	for _, x := range g {
		st.Pending = st.Pending || x.Name == guardName
	}
	rec, err := s.Load()
	if err != nil {
		st.Corrupt = err.Error()
		return st, nil
	}
	if rec != nil {
		for _, l := range rec.Layers {
			st.Layers = append(st.Layers, l.Name)
		}
	}
	return st, nil
}

// Resume is the explicit retry path after a failure kept the guard: it finishes the
// interrupted operation (re-apply layers, or restore base if a last pop was in flight)
// and releases the guard on success. It refuses on a corrupt record.
func (s *Service) Resume(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	active, err := s.guardActive()
	if err != nil {
		return err
	}
	if !active {
		return nil
	}
	rec, err := s.Load()
	if err != nil {
		return fmt.Errorf("dns layers: resume: %w (guard kept)", err)
	}
	if rec == nil {
		// record already deleted: the last restore applied, only the release was lost
		return s.Guards.Release(guardName)
	}
	// Staged dhcp changes may be someone else's: never commit or revert them. A crash
	// mid-Stage also lands here; `uci revert dhcp` or a reboot (staging is in /tmp) clears it.
	if pending, err := s.Device.Pending(ctx, "dhcp"); err != nil {
		return fmt.Errorf("dns layers: resume: %w (guard kept)", err)
	} else if pending {
		return fmt.Errorf("dns layers: resume: %w (guard kept)", ErrForeignStaged)
	}
	if rec.Restoring {
		return s.applyAndFinish(ctx, rec.Base, true)
	}
	return s.applyAndFinish(ctx, Effective(*rec), false)
}

func (s *Service) guardActive() (bool, error) {
	g, err := s.Guards.Stale()
	if err != nil {
		return false, err
	}
	for _, x := range g {
		if x.Name == guardName {
			return true, nil
		}
	}
	return false, nil
}

func (s *Service) mutate(ctx context.Context, f func(*Record, State) (*Record, error)) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.Guards.Acquire(guardName); err != nil {
		return fmt.Errorf("dns layers: %w", err)
	}
	rec, err := s.Load()
	if err != nil {
		return fmt.Errorf("dns layers: %w (guard kept)", err)
	}
	// Nothing has been written yet: transient read errors release the guard.
	release := func(err error) error {
		if rerr := s.Guards.Release(guardName); rerr != nil {
			return errors.Join(err, rerr)
		}
		return err
	}
	if pending, err := s.Device.Pending(ctx, "dhcp"); err != nil {
		return release(fmt.Errorf("dns layers: %w", err))
	} else if pending {
		return release(fmt.Errorf("dns layers: %w", ErrForeignStaged))
	}
	_, cur, err := ReadState(ctx, s.Device)
	if err != nil {
		return release(fmt.Errorf("dns layers: %w", err))
	}
	next, err := f(rec, cur)
	if err != nil {
		return release(err)
	}
	switch {
	case next != nil:
		if err := atomicfile.WriteJSON(s.RecordPath, next, 0o600); err != nil {
			return release(fmt.Errorf("dns layers: save record: %w", err))
		}
		return s.applyAndFinish(ctx, Effective(*next), false)
	case rec != nil:
		marked := *rec
		marked.Layers, marked.Restoring = nil, true
		if err := atomicfile.WriteJSON(s.RecordPath, marked, 0o600); err != nil {
			return release(fmt.Errorf("dns layers: save record: %w", err))
		}
		return s.applyAndFinish(ctx, rec.Base, true)
	default:
		return release(nil) // nothing active, nothing to do
	}
}

// applyAndFinish stages target and applies through the Applier (which reverts staging on
// failure), then deletes the record (if restoring base) and releases the guard.
// Any failure keeps the guard.
func (s *Service) applyAndFinish(ctx context.Context, target State, restoreBase bool) error {
	err := s.Apply.Apply(ctx, []string{"dhcp"}, func(lctx context.Context) error {
		sec, _, err := ReadState(lctx, s.Device)
		if err != nil {
			return err
		}
		return Stage(lctx, s.Device, sec, target)
	})
	if err != nil {
		return fmt.Errorf("dns layers: apply: %w (guard kept)", err)
	}
	if restoreBase {
		if err := os.Remove(s.RecordPath); err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("dns layers: delete record: %w (guard kept)", err)
		}
	}
	return s.Guards.Release(guardName)
}

// ReadState finds the dnsmasq section (it may be anonymous, see device-facts) and reads the owned options.
func ReadState(ctx context.Context, d device.Device) (string, State, error) {
	secs, err := d.Show(ctx, "dhcp")
	if err != nil {
		return "", State{}, err
	}
	var found []device.Section
	for _, s := range secs {
		if s.Type == "dnsmasq" {
			found = append(found, s)
		}
	}
	if len(found) != 1 {
		return "", State{}, fmt.Errorf("want exactly one dnsmasq section, found %d", len(found))
	}
	s := found[0]
	st := State{}
	if v, ok := s.Lists["server"]; ok {
		st.Server = List{Set: true, Values: slices.Clone(v)}
	} else if v, ok := s.Options["server"]; ok {
		st.Server = List{Set: true, Values: []string{v}, Scalar: true}
	}
	if v, ok := s.Options["noresolv"]; ok {
		st.NoResolv = Opt{Set: true, Value: v}
	}
	if v, ok := s.Options["rebind_protection"]; ok {
		st.RebindProtection = Opt{Set: true, Value: v}
	}
	return s.Name, st, nil
}

// Stage writes st into the section without committing. Absent options are deleted.
func Stage(ctx context.Context, d device.Device, sec string, st State) error {
	if err := d.DeleteOption(ctx, "dhcp", sec, "server"); err != nil && !errors.Is(err, device.ErrNotFound) {
		return err
	}
	switch {
	case st.Server.Set && st.Server.Scalar && len(st.Server.Values) == 1:
		if err := d.Set(ctx, "dhcp", sec, "server", st.Server.Values[0]); err != nil {
			return err
		}
	case st.Server.Set:
		for _, v := range st.Server.Values {
			if err := d.AddToList(ctx, "dhcp", sec, "server", v); err != nil {
				return err
			}
		}
	}
	for _, o := range []struct {
		name string
		opt  Opt
	}{{"noresolv", st.NoResolv}, {"rebind_protection", st.RebindProtection}} {
		var err error
		if o.opt.Set {
			err = d.Set(ctx, "dhcp", sec, o.name, o.opt.Value)
		} else {
			err = d.DeleteOption(ctx, "dhcp", sec, o.name)
		}
		if err != nil && !errors.Is(err, device.ErrNotFound) {
			return err
		}
	}
	return nil
}
