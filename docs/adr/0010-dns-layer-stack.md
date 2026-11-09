# 0010 One DNS layer stack

Status: accepted (property-tested, container-tested; device behaviour via spike S5)

Features that touch dnsmasq `server`/`noresolv`/rebind push and pop layers on one ordered record (`captive`, `vpn`, `adguard`, reserved `tailscale`). Base is captured before the first layer. Reads never mutate. Built from day one because per-feature snapshots of the same options are how the old bugs happened. See [invariants §8](../invariants.md).

Implementation (`internal/dnslayer`):
- Precedence is fixed by layer name, not push order: captive > adguard > vpn. The highest replace layer sets `server`/`noresolv`; tailscale is additive by name and appends domain-scoped servers. Any rebind-disabling layer turns rebind protection off. Effective state depends only on the active layers (property test with random payloads).
- Base keeps the `option` vs `list` shape of `server` and is restored in that shape.
- Push saves the record (intent) before staging. Last pop first marks the record `restoring`, then restores base, and deletes the record only after apply succeeded. Any failure after the first write reverts staging and keeps the guard; `Resume` finishes the interrupted operation. Read failures before any write release the guard.
- A change is refused while dhcp has foreign staged changes (`uci changes`): the stack never commits or reverts someone else's edits.
- `Status` is read-only and reports `pending` (guard held) and `corrupt` (unreadable record).
- Commit goes through an `Applier` (the transaction module); the stack never commits itself.
