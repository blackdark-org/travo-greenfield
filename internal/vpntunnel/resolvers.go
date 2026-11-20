package vpntunnel

import (
	"net/netip"
	"strings"

	"github.com/blackdark-org/travo-greenfield/internal/device"
	"github.com/blackdark-org/travo-greenfield/internal/vpnprofiles"
)

// ResolversOption is the Travo-owned option on network.wg0 holding the applied profile's
// `DNS =` resolvers, space-separated, or noResolvers (#289). Enable stages it in the tunnel
// transaction, so rollback restores it with wg0 and selecting another profile changes
// nothing until it is applied. netifd ignores the unknown option (procd container, 25.12.6).
const ResolversOption = "travo_dns"

const noResolvers = "none"

func resolversValue(dns []netip.Addr) string {
	if len(dns) == 0 {
		return noResolvers
	}
	out := make([]string, 0, len(dns))
	for _, a := range dns {
		out = append(out, a.String())
	}
	return strings.Join(out, " ")
}

// Recorded is what network.wg0 says about the applied resolvers.
type Recorded int

const (
	NoTunnel   Recorded = iota // no wg0 section
	Unrecorded                 // wg0 applied before #289: no ResolversOption yet
	Applied                    // Resolvers holds what Enable staged
)

// AppliedResolvers reads the resolvers the configured wg0 was applied with. A malformed
// value or a foreign wg0 section is errUnavailable.
func AppliedResolvers(network []device.Section) ([]netip.Addr, Recorded, error) {
	for _, s := range network {
		if s.Name != iface {
			continue
		}
		if s.Type != "interface" || s.Options["proto"] != "wireguard" || len(s.Lists[ResolversOption]) != 0 {
			return nil, NoTunnel, errUnavailable
		}
		raw, ok := s.Options[ResolversOption]
		switch {
		case !ok:
			return nil, Unrecorded, nil
		case raw == noResolvers:
			return nil, Applied, nil
		}
		var out []netip.Addr
		for _, v := range strings.Fields(raw) {
			a, err := netip.ParseAddr(v)
			if err != nil || a.Zone() != "" {
				return nil, NoTunnel, errUnavailable
			}
			out = append(out, a)
		}
		if len(out) == 0 {
			return nil, NoTunnel, errUnavailable
		}
		return out, Applied, nil
	}
	return nil, NoTunnel, nil
}

// ResolverSource reads the selected profile, for tunnels applied before #289.
type ResolverSource interface {
	SelectedConfig() (vpnprofiles.Config, error)
}

// TunnelResolvers are the `DNS =` resolvers the configured wg0 runs with: the recorded ones,
// else the selected profile's for a wg0 applied before #289. vpnprofiles.ErrNotFound means
// neither exists.
func TunnelResolvers(network []device.Section, selected ResolverSource) ([]netip.Addr, error) {
	dns, rec, err := AppliedResolvers(network)
	if err != nil || rec == Applied {
		return dns, err
	}
	cfg, err := selected.SelectedConfig()
	return cfg.DNS, err
}
