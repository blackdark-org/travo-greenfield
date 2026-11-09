package dnslayer

import (
	"errors"
	"math/rand/v2"
	"reflect"
	"testing"
)

var base = State{
	Server:           List{Set: true, Values: []string{"9.9.9.9"}},
	RebindProtection: Opt{Set: true, Value: "1"},
	// noresolv absent: restore must delete it again
}

var layers = map[Name]Layer{
	Captive:   {Name: Captive, Servers: []string{"10.0.0.1"}, DisableRebind: true},
	AdGuard:   {Name: AdGuard, Servers: []string{"127.0.0.1#5353"}, NoResolv: true},
	VPN:       {Name: VPN, Servers: []string{"10.2.0.1"}, NoResolv: true},
	Tailscale: {Name: Tailscale, Servers: []string{"/ts.net/100.100.100.100"}, Additive: true},
}

func push(t *testing.T, r *Record, n Name) *Record {
	t.Helper()
	out, err := Push(r, State{Server: List{Set: true, Values: []string{"LIVE-NOT-BASE"}}}, layers[n])
	if err != nil {
		t.Fatal(err)
	}
	if r == nil { // first push captures what was live
		out.Base = base
	}
	return &out
}

func TestPrecedenceIsFixedNotPushOrder(t *testing.T) {
	a := push(t, push(t, nil, VPN), AdGuard)
	b := push(t, push(t, nil, AdGuard), VPN)
	if !reflect.DeepEqual(Effective(*a), Effective(*b)) {
		t.Fatalf("order-dependent:\n%+v\n%+v", Effective(*a), Effective(*b))
	}
	if got := Effective(*a).Server.Values; !reflect.DeepEqual(got, []string{"127.0.0.1#5353"}) {
		t.Fatalf("adguard must win over vpn, got %v", got)
	}
}

func TestCaptiveWinsAndDisablesRebind(t *testing.T) {
	r := push(t, push(t, push(t, nil, AdGuard), VPN), Captive)
	st := Effective(*r)
	if !reflect.DeepEqual(st.Server.Values, []string{"10.0.0.1"}) || st.NoResolv.Value != "0" || st.RebindProtection.Value != "0" {
		t.Fatalf("%+v", st)
	}
}

func TestAdditiveAppendsToWinner(t *testing.T) {
	r := push(t, push(t, nil, VPN), Tailscale)
	if got := Effective(*r).Server.Values; !reflect.DeepEqual(got, []string{"10.2.0.1", "/ts.net/100.100.100.100"}) {
		t.Fatal(got)
	}
	only := push(t, nil, Tailscale)
	if got := Effective(*only).Server.Values; !reflect.DeepEqual(got, []string{"9.9.9.9", "/ts.net/100.100.100.100"}) {
		t.Fatalf("additive on base: %v", got)
	}
}

func TestPopNonTopRemovesOnlyItself(t *testing.T) {
	r := push(t, push(t, nil, VPN), Captive)
	r2, err := Pop(r, VPN)
	if err != nil || r2 == nil || !Active(r2, Captive) || Active(r2, VPN) {
		t.Fatalf("%+v %v", r2, err)
	}
	if !reflect.DeepEqual(r2.Base, base) {
		t.Fatal("base changed by pop")
	}
}

func TestPushIsIdempotentAndBaseCapturedOnce(t *testing.T) {
	r := push(t, nil, VPN)
	r2 := push(t, r, VPN)
	if !reflect.DeepEqual(r, r2) {
		t.Fatalf("not idempotent:\n%+v\n%+v", r, r2)
	}
}

func TestPopOnEmptyIsNoOp(t *testing.T) {
	if r, err := Pop(nil, VPN); r != nil || err != nil {
		t.Fatal(r, err)
	}
}

func TestRejectsBadInputAndCorruptRecords(t *testing.T) {
	if _, err := Push(nil, base, Layer{Name: "evil", Servers: []string{"x"}}); !errors.Is(err, ErrUnknownLayer) {
		t.Error(err)
	}
	if _, err := Push(nil, base, Layer{Name: VPN}); !errors.Is(err, ErrNoServers) {
		t.Error(err)
	}
	bad := []Record{
		{Version: 2, Base: base, Layers: []Layer{layers[VPN]}},
		{Version: 1, Base: base},
		{Version: 1, Base: base, Layers: []Layer{layers[VPN], layers[VPN]}},
	}
	for _, r := range bad {
		if _, err := Push(&r, base, layers[AdGuard]); !errors.Is(err, ErrBadRecord) {
			t.Errorf("push accepted %+v: %v", r, err)
		}
		if _, err := Pop(&r, VPN); !errors.Is(err, ErrBadRecord) {
			t.Errorf("pop accepted %+v: %v", r, err)
		}
	}
}

// Property: for any sequence of pushes and pops, the effective state equals the state
// computed from the set of active layers alone, and emptying the stack yields no record.
func TestPropertyRandomSequences(t *testing.T) {
	names := []Name{Captive, AdGuard, VPN, Tailscale}
	rng := rand.New(rand.NewPCG(1, 2))
	for iter := 0; iter < 2000; iter++ {
		var r *Record
		active := map[Name]bool{}
		for step := 0; step < 12; step++ {
			n := names[rng.IntN(len(names))]
			if rng.IntN(2) == 0 {
				r = push(t, r, n)
				active[n] = true
			} else {
				var err error
				if r, err = Pop(r, n); err != nil {
					t.Fatal(err)
				}
				delete(active, n)
			}
			if (r == nil) != (len(active) == 0) {
				t.Fatalf("record presence %v vs active %v", r != nil, active)
			}
			if r == nil {
				continue
			}
			var fresh *Record
			for _, nn := range names {
				if active[nn] {
					fresh = push(t, fresh, nn)
				}
			}
			if !reflect.DeepEqual(Effective(*r), Effective(*fresh)) {
				t.Fatalf("history-dependent state at iter %d: %+v vs %+v", iter, Effective(*r), Effective(*fresh))
			}
			if !reflect.DeepEqual(r.Base, base) {
				t.Fatal("base drifted")
			}
		}
	}
}
