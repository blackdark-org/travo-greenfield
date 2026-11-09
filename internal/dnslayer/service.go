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

// Applier commits staged UCI configs and reloads services. The transaction module
// implements it; nothing here commits directly.
type Applier interface {
	Apply(ctx context.Context, configs []string) error
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

// Push activates a layer. The record (intent) is saved before the commit, so a crash
// in between leaves the guard and a record that a retry re-applies idempotently.
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
	sec, cur, err := ReadState(ctx, s.Device)
	if err != nil {
		return fmt.Errorf("dns layers: %w (guard kept)", err)
	}
	next, err := f(rec, cur)
	if err != nil {
		_ = s.Guards.Release(guardName) // validation error: nothing was written
		return err
	}
	var target State
	switch {
	case next != nil:
		if err := atomicfile.WriteJSON(s.RecordPath, next, 0o600); err != nil {
			return fmt.Errorf("dns layers: save record: %w (guard kept)", err)
		}
		target = Effective(*next)
	case rec != nil:
		target = rec.Base
	default:
		return s.Guards.Release(guardName) // nothing active, nothing to do
	}
	if err := Stage(ctx, s.Device, sec, target); err != nil {
		_ = s.Device.Revert(ctx, "dhcp")
		return fmt.Errorf("dns layers: stage: %w (guard kept)", err)
	}
	if err := s.Apply.Apply(ctx, []string{"dhcp"}); err != nil {
		return fmt.Errorf("dns layers: apply: %w (guard kept)", err)
	}
	if next == nil {
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
		st.Server = List{Set: true, Values: []string{v}}
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
	if st.Server.Set {
		for _, v := range st.Server.Values {
			if err := d.AddToList(ctx, "dhcp", sec, "server", v); err != nil {
				return err
			}
		}
	}
	for opt, o := range map[string]Opt{"noresolv": st.NoResolv, "rebind_protection": st.RebindProtection} {
		var err error
		if o.Set {
			err = d.Set(ctx, "dhcp", sec, opt, o.Value)
		} else {
			err = d.DeleteOption(ctx, "dhcp", sec, opt)
		}
		if err != nil && !errors.Is(err, device.ErrNotFound) {
			return err
		}
	}
	return nil
}
