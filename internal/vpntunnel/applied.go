package vpntunnel

import (
	"context"

	"github.com/blackdark-org/travo-greenfield/internal/device"
	"github.com/blackdark-org/travo-greenfield/internal/vpnprofiles"
)

// ConfigFinder looks up stored profiles (vpnprofiles.Service).
type ConfigFinder interface {
	FindConfig(match func(vpnprofiles.Config) bool) (vpnprofiles.Config, error)
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
// edited outside Travo), and errUnavailable when network cannot be read.
func (a Applied) AppliedConfig(ctx context.Context) (vpnprofiles.Config, error) {
	current, err := a.UCI.Show(ctx, "network")
	if err != nil {
		return vpnprofiles.Config{}, errUnavailable
	}
	split, err := currentSplit(current)
	if err != nil {
		return vpnprofiles.Config{}, err
	}
	return a.Profiles.FindConfig(func(cfg vpnprofiles.Config) bool {
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
}
