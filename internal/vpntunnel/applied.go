package vpntunnel

import (
	"context"
	"errors"
	"net/netip"
	"slices"

	"github.com/blackdark-org/travo-greenfield/internal/device"
	"github.com/blackdark-org/travo-greenfield/internal/vpnprofiles"
)

// ErrAmbiguous: several stored profiles match wg0 but their DNS lines differ, so the
// resolvers of the running tunnel are unknown.
var ErrAmbiguous = errors.New("several VPN profiles match the configured tunnel")

// ConfigFinder looks up stored profiles (vpnprofiles.Service).
type ConfigFinder interface {
	MatchingConfigs(match func(vpnprofiles.Config) bool) ([]vpnprofiles.Config, error)
}

// Applied finds the stored profile that the configured wg0 was applied from (#289), so the
// `vpn` DNS layer follows the running tunnel, not a later selection.
type Applied struct {
	UCI interface {
		Show(context.Context, string) ([]device.Section, error)
	}
	Profiles ConfigFinder
}

// AppliedConfig returns vpnprofiles.ErrNotFound when no profile matches wg0 (absent, or
// edited outside Travo), ErrAmbiguous when matches disagree on DNS, and errUnavailable
// when network cannot be read.
func (a Applied) AppliedConfig(ctx context.Context) (vpnprofiles.Config, error) {
	current, err := a.UCI.Show(ctx, "network")
	if err != nil {
		return vpnprofiles.Config{}, errUnavailable
	}
	return appliedFrom(current, a.Profiles)
}

func appliedFrom(current []device.Section, profiles ConfigFinder) (vpnprofiles.Config, error) {
	split, err := currentSplit(current)
	if err != nil {
		return vpnprofiles.Config{}, err
	}
	found, err := profiles.MatchingConfigs(func(cfg vpnprofiles.Config) bool {
		desired, err := desiredSections(cfg)
		if err != nil {
			return false
		}
		if len(split) > 0 { // as Set stages it: split routes turn route_allowed_ips off
			desired = withRouteAll(desired, "0")
		}
		// enabled=true: compare the configured sections whatever wg0's disabled value is
		return matches(current, desired, true)
	})
	if err != nil {
		return vpnprofiles.Config{}, err
	}
	if len(found) == 0 {
		return vpnprofiles.Config{}, vpnprofiles.ErrNotFound
	}
	// wg0 does not carry DNS, so profiles differing only there are indistinguishable
	for _, cfg := range found[1:] {
		if !slices.EqualFunc(cfg.DNS, found[0].DNS, func(a, b netip.Addr) bool { return a == b }) {
			return vpnprofiles.Config{}, ErrAmbiguous
		}
	}
	return found[0], nil
}

// tunnelProfile is the profile that split routes are built for: the applied one, else the
// selected one (wg0 edited outside Travo, or matches that disagree on DNS).
func (s *Service) tunnelProfile(current []device.Section) (vpnprofiles.Config, error) {
	if cfg, err := appliedFrom(current, s.profiles); err == nil {
		return cfg, nil
	}
	return s.profiles.SelectedConfig()
}
