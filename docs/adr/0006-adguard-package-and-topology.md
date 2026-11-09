# 0006 AdGuard Home: package, topology, proxy, fail-open

Status: accepted; S1, S2, S5 verified on GL-AXT1800 25.12.6 (testdata/device/axt1800-25.12.6/SPIKES.md)

- **Install source**: OpenWrt `adguardhome` package via the PackageManager port. OpenWrt owns binary and conffile (survives sysupgrade); Travo supervises and configures. Rejected: upstream binary in `/opt` (no conffile tracking, schema downgrade crash loop).
- **Topology**: dnsmasq stays on :53 and forwards to AGH on `127.0.0.1:5353`. dnsmasq remains the single DNS authority for router self-resolution, DHCP names, hotel DNS, VPN DNS and Tailscale (`--accept-dns=false`, resolver as a layer). Cost: AGH sees queries from 127.0.0.1 (no per-client stats). Expert option later: AGH on :53 via the layer stack.
- **Access**: AGH binds loopback, no accounts. Travo serves a gated reverse proxy on `:3001` (own port: no base-path rewriting). Fallback if the proxy fails: sync one AGH user at login/password change.
- **Travo UI**: on/off and "default DNS" toggle only; everything else in AGH's UI.
- **First boot**: no uplink means no package. Install after the first successful uplink as a guarded package op with a flash check; DNS stays at base until the liveness probe passes.
- **Failure**: liveness probe resolves a fresh name via 127.0.0.1:5353; after M failures pop the `adguard` layer (fail open), show a banner, re-enable only on user confirm or after a hold-down of passing probes. Refuse AGH downgrades over a newer config.

Verified facts: package 0.107.79-r1, config `/etc/adguardhome/adguardhome.yaml` (conffile, kept by sysupgrade), runs in ujail. A root-path reverse proxy on its own port works with `users: []`. Required change from the package/old-Travo state: `http.address` must become `127.0.0.1:3000` (it was `0.0.0.0:3000`, admin API open to the LAN) and DNS `bind_hosts` should be `127.0.0.1` only. Migration from an existing install must do both under a transaction.
