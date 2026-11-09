# Testing

Unit-green does not mean it works. Several real bugs shipped green because a fixture invented a field the device never emits.

## Layers

1. **Unit** with fixtures: verbatim captured device payloads, each with provenance (device, firmware version, date) in `testdata/`. Never hand-write a payload field.
2. **Container integration**: the official `openwrt/rootfs` image (25.12; 24.10 for the opkg adapter) with real `uci`, `ubus`, `rpcd`. Catches mocks that diverge from real semantics.
3. **Hardware playbook** per release on a real device. Some behaviour is provable only there: wireless apply/confirm, lockout, DNS composition with VPN + captive + AdGuard, sysupgrade survival, AdGuard reinstall, watchdog restore.

## Rules

- Mocks reproduce real semantics: `uci set` upserts (a mock that errors on an existing section is a false failure); a mock that ignores the device name is a false pass.
- A constant asserted against itself proves nothing. Check timing budgets against measured latency.
- Test the fail-closed direction explicitly: a broken AP reads DOWN; an unknown caller is refused.
- Property tests for the DNS layer stack (any push/pop order restores base).
- Table-driven tests for caller classification using real payload shapes.
- No destructive probes (e.g. reboot POST) against an unknown build: GET-only reachability first.
- Verify the running binary is the one under test (path + version) before measuring.

## Contract tests

Scripts parse; guard paths under `/etc/travo-guards`; every route is in the spec; auth covers every route (walk the live route table); strict-PUT coverage; spec field names match models; no feature writes UCI outside the transaction module; generated files (OpenAPI, TS client) have no diff.

## Mock mode

`--mock` runs the whole daemon against the mock `Device`, scripted from fixtures. Frontend development needs no device.

## Running the container layer

`mise run test:integration` cross-compiles the `integration`-tagged tests and runs them against real `uci`, `apk` (`openwrt/rootfs:aarch64_generic-25.12.6`) and `opkg` (`openwrt/rootfs:x86-64-24.10.8`, amd64: no arm64 24.10 manifest exists) (needs Docker). A second step boots the image with `/sbin/init` (procd as PID 1, real ubusd/rpcd/netifd) for tests that need procd. Not in `mise run check` yet; CI needs an x86_64 image or QEMU first. Container fixtures prove format and quoting only, never wireless or netifd behaviour; mark them as such in `testdata/PROVENANCE.md`.

## Hardware layer

`hardware`-tagged tests (`internal/txn/hardware_test.go`, `internal/dnslayer/hardware_test.go`) run ON the router from Ethernet: cross-compile with `GOOS=linux GOARCH=arm64 go test -tags hardware -c`, copy with `scp -O`, run with `-test.timeout` (no `timeout` in BusyBox). They change live wireless/DNS state and rely on the watchdog/stack to restore it; take a sysupgrade backup first and pause other writers (old Travo, its cron). The reboot test needs a temporary boot hook that calls `Recover`; remove it after. Results of the first run: `testdata/device/axt1800-25.12.6/SPIKES.md`.

## Hardware playbook (run before 0.1.0; console or Ethernet fallback ready)

GET-only captures first (spike S3), saved verbatim as fixtures with provenance:
- `ubus call network.interface dump`, `ubus call network.device status`, `ip -4 neigh`, `ip -4 route show table main`, `ubus list hostapd.*` and `ubus call hostapd.<if> get_clients`, `uci show wireless`, `uci show network`, `uci show dhcp`
- once on Ethernet, once from a phone on the AP, once with an uplink STA up

Live checks (spike S4), each from Ethernet:
1. Begin a wireless change, do not confirm: WiFi returns after the window via `reload_config` alone. Measure time to AP up and STA carrier; set the confirm window from it.
2. Begin, then `kill -9` the daemon: the watchdog restores.
3. Begin, then reboot within the window: startup recovery restores.
4. From a phone on the AP, switch to client mode: refused with 409 `wifi_lockout_risk`.
5. Probe with a wrong STA password: carrier stays false, rollback happens, UI says "rolled back" only after restore.

## Web

`pnpm -C web check`: Biome, svelte-check, Vitest (pure logic and server-rendered components), build, size-limit. Visual and keyboard behaviour (theme, drawer, focus) is checked in headless Chrome at 390 px and 1280 px before a UI PR; attach screenshots to the PR.
