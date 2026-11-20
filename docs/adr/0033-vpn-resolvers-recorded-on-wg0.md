# 0033 The tunnel transaction records the applied resolvers on wg0

Status: accepted (#289); changes the resolver source of the `vpn` layer in [0010](0010-dns-layer-stack.md)

## Context

The `vpn` layer pushed the SELECTED profile's `DNS =` resolvers. Selecting another profile without re-applying switched DNS to its resolvers while traffic still used the old tunnel. `wg0` carries no DNS, so the applied profile cannot be derived from it: matching stored profiles against `wg0` raced with staging and could not tell apart profiles that differ only in DNS (#295 review).

## Decision

- Enable stages the applied profile's resolvers as `network.wg0.travo_dns` (space-separated, or `none`) in the WireGuard transaction, so rollback restores them with `wg0`. If `wg0` already runs as the profile, only that option is staged, without a restart or handshake proof.
- The reconciler reads them under the `network` write lock (`txn.Manager.ShowSettled`), checks the transaction guard before and after that read, and reads the tunnel state last, so it never acts on a staged delta or an unconfirmed apply. An unreadable value changes nothing, except that a tunnel that is off still pops.
- A `wg0` without the option (applied before #289) uses the selected profile until the next Enable records it.
- Split routes use the same resolvers; when Enable changes them, the resolver `/32`s move with them and `wg0` restarts (#296). With unchanged DNS the split stays as it is. A user route equal to an old resolver address cannot be told apart from the automatic `/32`, so it moves too when the DNS changes.

## Consequences

The tunnel transaction writes DNS metadata but still never touches dnsmasq. netifd ignores the option: in a booted `aarch64_generic-25.12.6` container `wg0` came up with it and kept its uptime and ifindex when only the option changed (`reload_config`). The hardware check is in the "VPN DNS layer" playbook step.
