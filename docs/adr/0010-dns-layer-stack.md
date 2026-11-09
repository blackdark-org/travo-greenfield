# 0010 One DNS layer stack

Status: accepted (property-tested, container-tested; device behaviour via spike S5)

Features that touch dnsmasq `server`/`noresolv`/rebind push and pop layers on one ordered record (`captive`, `vpn`, `adguard`, reserved `tailscale`). Base is captured before the first layer. Reads never mutate. Built from day one because per-feature snapshots of the same options are how the old bugs happened. See [invariants §8](../invariants.md).

Implementation (`internal/dnslayer`):
- Precedence is fixed by layer name, not push order: captive > adguard > vpn. The highest replace layer sets `server`/`noresolv`; additive layers (tailscale) append domain-scoped servers. Any rebind-disabling layer turns rebind protection off. Effective state depends only on the set of active layers (property test).
- Push saves the record (intent) before staging; the guard stays on any failure after that, so a retry re-applies idempotently. The last pop restores base and deletes the record only after apply succeeded.
- Commit goes through an `Applier` (the transaction module); the stack never commits itself.
