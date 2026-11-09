package dnslayer

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/BlackDark/travo-greenfield/internal/device/mock"
	"github.com/BlackDark/travo-greenfield/internal/guard"
)

type commitApplier struct {
	d    *mock.Device
	fail bool
	n    int
}

func (a *commitApplier) Apply(ctx context.Context, configs []string) error {
	a.n++
	if a.fail {
		return errors.New("apply failed")
	}
	for _, c := range configs {
		if err := a.d.Commit(ctx, c); err != nil {
			return err
		}
	}
	return nil
}

var ctx = context.Background()

func setup(t *testing.T) (*Service, *mock.Device, *commitApplier) {
	t.Helper()
	d := mock.New()
	// dnsmasq section shape from internal/device/cli/testdata/uci-show-dhcp.txt
	_ = d.SetType(ctx, "dhcp", "@dnsmasq[0]", "dnsmasq")
	_ = d.Set(ctx, "dhcp", "@dnsmasq[0]", "rebind_protection", "1")
	_ = d.AddToList(ctx, "dhcp", "@dnsmasq[0]", "server", "9.9.9.9")
	_ = d.SetType(ctx, "dhcp", "lan", "dhcp")
	_ = d.Commit(ctx, "dhcp")
	dir := t.TempDir()
	a := &commitApplier{d: d}
	return &Service{
		Device:     d,
		Guards:     guard.Dir{Path: filepath.Join(dir, "guards"), Now: time.Now},
		RecordPath: filepath.Join(dir, "dns-layers.json"),
		Apply:      a,
	}, d, a
}

func live(t *testing.T, d *mock.Device) State {
	t.Helper()
	_, st, err := ReadState(ctx, d)
	if err != nil {
		t.Fatal(err)
	}
	return st
}

func TestPushThenPopRestoresBaseExactly(t *testing.T) {
	s, d, _ := setup(t)
	before := live(t, d)
	if err := s.Push(ctx, layers[AdGuard]); err != nil {
		t.Fatal(err)
	}
	if err := s.Push(ctx, layers[Captive]); err != nil {
		t.Fatal(err)
	}
	st := live(t, d)
	if !reflect.DeepEqual(st.Server.Values, []string{"10.0.0.1"}) || st.RebindProtection.Value != "0" {
		t.Fatalf("captive not effective: %+v", st)
	}
	if err := s.Pop(ctx, Captive); err != nil {
		t.Fatal(err)
	}
	if got := live(t, d).Server.Values; !reflect.DeepEqual(got, []string{"127.0.0.1#5353"}) {
		t.Fatalf("adguard not back: %v", got)
	}
	if err := s.Pop(ctx, AdGuard); err != nil {
		t.Fatal(err)
	}
	if after := live(t, d); !reflect.DeepEqual(before, after) {
		t.Fatalf("base not restored:\n%+v\n%+v", before, after)
	}
	if _, err := os.Stat(s.RecordPath); !errors.Is(err, os.ErrNotExist) {
		t.Fatal("record not deleted after last pop")
	}
	if st, _ := s.Guards.Stale(); len(st) != 0 {
		t.Fatalf("guard left: %v", st)
	}
}

func TestFailedApplyKeepsGuardAndRecord(t *testing.T) {
	s, _, a := setup(t)
	a.fail = true
	if err := s.Push(ctx, layers[VPN]); err == nil {
		t.Fatal("want error")
	}
	if st, _ := s.Guards.Stale(); len(st) != 1 {
		t.Fatal("guard removed after failed apply")
	}
	if rec, err := s.Load(); err != nil || !Active(rec, VPN) {
		t.Fatalf("intent record lost: %v %v", rec, err)
	}
	if err := s.Push(ctx, layers[VPN]); !errors.Is(err, guard.ErrActive) {
		t.Fatalf("feature must skip while guarded: %v", err)
	}
}

func TestLastPopFailureNeverReportsSuccess(t *testing.T) {
	s, _, a := setup(t)
	_ = s.Push(ctx, layers[Captive])
	a.fail = true
	if err := s.Pop(ctx, Captive); err == nil {
		t.Fatal("restore reported success while bypass may be in force")
	}
	if _, err := os.Stat(s.RecordPath); err != nil {
		t.Fatal("record deleted before base was restored")
	}
}

func TestCorruptRecordRefusesAndKeepsGuard(t *testing.T) {
	s, _, a := setup(t)
	_ = os.WriteFile(s.RecordPath, []byte("{"), 0o600)
	if err := s.Pop(ctx, VPN); err == nil {
		t.Fatal("accepted corrupt record")
	}
	if a.n != 0 {
		t.Fatal("applied with corrupt record")
	}
	if st, _ := s.Guards.Stale(); len(st) != 1 {
		t.Fatal("guard not kept")
	}
}

func TestInvalidLayerWritesNothing(t *testing.T) {
	s, _, a := setup(t)
	if err := s.Push(ctx, Layer{Name: "evil", Servers: []string{"1.1.1.1"}}); !errors.Is(err, ErrUnknownLayer) {
		t.Fatal(err)
	}
	if a.n != 0 {
		t.Fatal("applied")
	}
	if st, _ := s.Guards.Stale(); len(st) != 0 {
		t.Fatal("guard left after validation error")
	}
}

func TestLoadIsReadOnly(t *testing.T) {
	s, _, a := setup(t)
	_ = s.Push(ctx, layers[VPN])
	n := a.n
	for i := 0; i < 3; i++ {
		if _, err := s.Load(); err != nil {
			t.Fatal(err)
		}
	}
	if a.n != n {
		t.Fatal("read mutated DNS")
	}
}
