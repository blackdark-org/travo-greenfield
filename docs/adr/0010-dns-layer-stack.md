# 0010 One DNS layer stack

Status: accepted (property-tested; device behaviour via spike S5)

Features that touch dnsmasq `server`/`noresolv`/rebind push and pop layers on one ordered record (`captive`, `vpn`, `adguard`, reserved `tailscale`). Base is captured before the first layer. Reads never mutate. Built from day one because per-feature snapshots of the same options are how the old bugs happened. See [invariants §8](../invariants.md).
