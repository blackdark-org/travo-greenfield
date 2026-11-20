package vpntunnel

import (
	"context"
	"encoding/base64"
	"errors"
	"net/netip"
	"path/filepath"
	"testing"

	"github.com/blackdark-org/travo-greenfield/internal/vpnprofiles"
)

func testKey(b byte) string {
	var raw [32]byte
	raw[0] = b
	return base64.StdEncoding.EncodeToString(raw[:])
}

// profileText is a WireGuard profile with its own key and resolver; peer 2 is shared.
func profileText(key byte, dns string) string {
	return "[Interface]\nPrivateKey=" + testKey(key) + "\nAddress=100.64.41.140/32\nDNS=" + dns +
		"\n[Peer]\nPublicKey=" + testKey(2) + "\nAllowedIPs=0.0.0.0/0\nEndpoint=203.0.113.8:51820\n"
}

// #289: the resolvers follow the profile wg0 was applied from, not a later selection.
func TestAppliedConfigFollowsTheAppliedProfile(t *testing.T) {
	ctx := context.Background()
	store := vpnprofiles.New(filepath.Join(t.TempDir(), "profiles.json"))
	first, err := store.Import("first", profileText(1, "198.18.0.1"))
	if err != nil {
		t.Fatal(err)
	}
	second, err := store.Import("second", profileText(3, "10.8.0.1"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := store.Select(first.ID); err != nil {
		t.Fatal(err)
	}
	s, d, _ := splitService(t, "disabled")
	s.profiles = store
	if _, err := s.Set(ctx, true); err != nil {
		t.Fatal(err)
	}
	applied := Applied{UCI: d, Profiles: store}
	want := []netip.Addr{netip.MustParseAddr("198.18.0.1")}
	check := func(when string) {
		t.Helper()
		cfg, err := applied.AppliedConfig(ctx)
		if err != nil || len(cfg.DNS) != 1 || cfg.DNS[0] != want[0] {
			t.Fatalf("%s: %v %v, want %v", when, cfg.DNS, err, want)
		}
	}
	check("after apply")
	if _, err := store.Select(second.ID); err != nil {
		t.Fatal(err)
	}
	check("after selecting another profile")
	if _, err := s.SetSplit(ctx, prefixes4(t, "10.9.0.0/16")); err != nil {
		t.Fatal(err)
	}
	check("with a split selection")
	if _, err := s.Set(ctx, true); err != nil { // re-apply: the selection now wins
		t.Fatal(err)
	}
	want = []netip.Addr{netip.MustParseAddr("10.8.0.1")}
	check("after re-applying")
}

// The captured wg0 was written with plain uci, not from a stored profile: no match.
func TestAppliedConfigWithoutMatchingProfile(t *testing.T) {
	store := vpnprofiles.New(filepath.Join(t.TempDir(), "profiles.json"))
	if _, err := store.Import("other", profileText(1, "198.18.0.1")); err != nil {
		t.Fatal(err)
	}
	for _, state := range []string{"connected", "disabled"} {
		applied := Applied{UCI: capturedNetwork(t, state), Profiles: store}
		if _, err := applied.AppliedConfig(context.Background()); !errors.Is(err, vpnprofiles.ErrNotFound) {
			t.Fatalf("%s: %v", state, err)
		}
	}
}
