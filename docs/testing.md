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
