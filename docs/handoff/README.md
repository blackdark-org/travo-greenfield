# Handoff (temporary)

Working state for the next agent or session. **Temporary by design**: when an item is done, move anything durable into the long-term docs (spec, invariants, device-facts, testing, ADRs, roadmap) and delete it here. Delete this directory when it is empty.

Long-term docs stay the source of truth: [../spec.md](../spec.md), [../roadmap.md](../roadmap.md), [../invariants.md](../invariants.md), [../device-facts.md](../device-facts.md), [../testing.md](../testing.md), [../adr/](../adr/).

## Status (2026-11-09)

| Slice (roadmap) | State |
|---|---|
| 1 tooling | merged |
| 2 daemon skeleton, ports | merged (#5) |
| 5 DNS layer stack | merged (#6) |
| 4 transactions, watchdog, lockout, device parsers, hardware spikes S1–S5 | PR #7 |
| 3 auth + first-run | next, see [slice-3-auth.md](slice-3-auth.md) |

Open release-please PR #4 proposes `1.0.0`. Do not merge it; the first release must be `0.1.0` (set `release-as` in `release-please-config.json` before 0.1.0, then remove it).

## Device access

- Test router: GL-AXT1800, OpenWrt 25.12.6, `ssh root@192.168.1.1` (key auth) when the operator has it on LAN. Old Travo 0.10.2 runs on it (`:80`; LuCI `:8080`/`:8443`; AdGuard `:3000`/`:5353`).
- Backup taken before the first spikes: `~/travo-router-backups/20261109-171110-axt1800/` on the operator's machine (not in the repo), with `RESTORE.md`.
- Before any live test, read the "Hardware layer" section in [../testing.md](../testing.md): stop old Travo and pause its per-minute cron (`/etc/travo/wifi-reconnect.sh`), restore both after, and verify configs against the backup.
- Never `wifi` / `wifi up/down/reload`. GET-only first.

## Nothing left that needs the device right now

All spikes done: [../../testdata/device/axt1800-25.12.6/SPIKES.md](../../testdata/device/axt1800-25.12.6/SPIKES.md). Device needed again only for:
- each feature slice's own playbook item (e.g. lockout refusal from a phone, captive portal in a real hotel network, AdGuard migration to loopback),
- the 0.1.0 release playbook in [../testing.md](../testing.md).

## Known issues found on the device (not ours to fix there)

- Old Travo left AdGuard's admin API open to the LAN (`0.0.0.0:3000`, `users: []`). Greenfield's AdGuard slice must migrate it to `127.0.0.1` under a transaction ([../adr/0006-adguard-package-and-topology.md](../adr/0006-adguard-package-and-topology.md)).
