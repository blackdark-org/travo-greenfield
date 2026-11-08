# 0006 AdGuard Home: package, topology, proxy, fail-open

Status: accepted; package facts and proxy behaviour pending spikes S1, S2, S5

- **Install source**: OpenWrt `adguardhome` package via the PackageManager port. OpenWrt owns binary and conffile (survives sysupgrade); Travo supervises and configures. Rejected: upstream binary in `/opt` (no conffile tracking, schema downgrade crash loop).
- **Topology**: dnsmasq stays on :53 and forwards to AGH on `127.0.0.1:5353`. dnsmasq remains the single DNS authority for router self-resolution, DHCP names, hotel DNS, VPN DNS and Tailscale (`--accept-dns=false`, resolver as a layer). Cost: AGH sees queries from 127.0.0.1 (no per-client stats). Expert option later: AGH on :53 via the layer stack.
- **Access**: AGH binds loopback, no accounts. Travo serves a gated reverse proxy on `:3001` (own port: no base-path rewriting). Fallback if the proxy fails: sync one AGH user at login/password change.
- **Travo UI**: on/off and "default DNS" toggle only; everything else in AGH's UI.
- **First boot**: no uplink means no package. Install after the first successful uplink as a guarded package op with a flash check; DNS stays at base until the liveness probe passes.
- **Failure**: liveness probe resolves a fresh name via 127.0.0.1:5353; after M failures pop the `adguard` layer (fail open), show a banner, re-enable only on user confirm or after a hold-down of passing probes. Refuse AGH downgrades over a newer config.
