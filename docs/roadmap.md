# Roadmap

Nothing here may bypass the `Device` port, the transaction module or the DNS layer stack. Each feature lists its constraints from [invariants.md](invariants.md).

## 0.1.0

See [spec.md](spec.md). Order of work (vertical slices, each shippable in mock mode first):

1. (done) Repo tooling: `mise.toml`, `mise run check`, CI from the agent-presets snippets, release-please, size-limit.
2. Daemon skeleton: lifecycle, Huma, health/spec, `Device` port + mock, guards, atomic files.
3. Auth + first-run + LAN-only listener.
4. Transaction module + watchdog + lockout classifier (spike on device first).
5. DNS layer stack (property tests).
6. WiFi uplink/AP/guest.
7. Captive assist.
8. WireGuard.
9. AdGuard (package, topology, proxy, fail-open).
10. Dashboard, system, backup, self-update.

### Spikes before the dependent slice

| Spike | Unblocks | Question |
|---|---|---|
| S1 AGH package | 9 | `adguardhome` package name, config path, conffile, sysupgrade behaviour, bind to `127.0.0.1:5353` with no users |
| S2 AGH proxy | 9 | AGH behind the Travo-gated `:3001` proxy with auth disabled (AGH#7985); first-run setup API; fallback = sync one AGH user |
| S3 Caller payloads | 4 | `hostapd.*` `get_clients`, neighbour table, `network.interface dump` captured verbatim as fixtures |
| S4 Watchdog | 4 | Restore path survives daemon kill; `uci reload_config` + `network reload` is enough |
| S5 DNS on device | 5, 9 | dnsmasq `:53` → `127.0.0.1#5353` composes with WireGuard DNS and captive bypass |

## 0.2

- USB tethering (RNDIS/NCM, iOS via usbmuxd) as an uplink
- Tailscale (`--accept-dns=false`, resolver as a dnsmasq layer, exit node, SSH option)
- API tokens for automation (hashed at rest)
- Auto-accept captive portals (step walker)

## Later (all wanted eventually)

- **Failover (mwan3, IPv4)**, per invariants §9
- **Clients**: hostname/MAC/IP, block/kick, alias, reservations, leases, blocked list
- **Network**: WAN (DHCP/static/PPPoE), LAN, DHCP server, local DNS entries, port forwards, firewall zones view, IPv6 settings, DoH, DDNS, WoL, diagnostics, data usage + budget, uptime log, traffic history, interface up/down
- **System**: hostname, timezone, NTP, LEDs (+schedule), hardware buttons (action mapping), alerts, logs, firmware upgrade (own image, warn about dropped packages), factory reset, SSH keys, update checks (OpenWrt versions, package upgrades)
- **Services catalog** (install/remove/autostart/status): SQM, cloudflared, vnstat, watchcat, speedtest
- WiFi schedule
- Persistent stats history (needs a store and a retention rule)
- HTTPS (self-signed), expert DNS topology (AGH on `:53`), real `.apk`, offline AGH bundle
- Not wanted: OpenVPN, Bluetooth tether, Mesh/WDS, AmneziaWG (revisit on request)
