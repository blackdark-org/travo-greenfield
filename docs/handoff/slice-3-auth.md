# Slice 3: auth and first-run (temporary plan)

Delete this file when slice 3 is merged; durable parts belong in ADR 0008, invariants §3 and spec.

Decisions already made: [../adr/0008-auth-and-trust-boundary.md](../adr/0008-auth-and-trust-boundary.md), [../invariants.md](../invariants.md) §3.

## Device facts for this slice

Captures: [../../testdata/device/axt1800-25.12.6/auth/](../../testdata/device/axt1800-25.12.6/auth/).
- `ubus call session login {"username","password","timeout"}` is the password check. Wrong password AND unknown user both fail identically: `Permission denied`, exit 250. Treat both as "invalid credentials"; never distinguish (user enumeration).
- rpcd `login` section maps `root` to `$p$root` (password from `/etc/shadow`), `read='*'`, `write='*'`. A successful login creates an rpcd session we do not need: call `session destroy` right after, or pass a short `timeout`. Travo keeps its own cookie session.
- `ubus call session login` works from the root shell without a session, so the daemon (root) can verify passwords through the `Device` port; no HTTP/JSON-RPC to uhttpd needed.
- A successful-login reply shape is not captured from hardware (needs the real password): capture it in the booted container (`/sbin/init`, set a known root password with `passwd`), mark provenance as container.
- No-password state: `/etc/shadow` root field empty means rpcd allows any password. First-run detection must read shadow (root-only file) rather than try a login.
- Listeners today: old Travo `0.0.0.0:80`, uhttpd `0.0.0.0:8080`/`8443`, AdGuard `:::3000`. Firewall zones in `firewall-zones.txt`.

## Build order

1. Password verifier over the `Device` port (`ubus call session login` + `destroy`), container test with a known password, uniform error.
2. Session store: memory, monotonic clock, 24 h, relative `expires_in`; epoch file in `/etc/<name>/` (atomic) bumped on password change.
3. Middleware: default-deny route table walk test, case-sensitive paths, CSRF header for unsafe methods, `SameSite=Strict` HttpOnly cookie, per-IP login rate limit with periodic sweep.
4. Trust boundary: daemon refuses requests whose source is not in a `br-lan` subnet (parsed via `facts.ParseInterfaceDump`), in addition to firewall rules.
5. First-run: no root password → only "set password" route, LAN-only; sets password via `passwd` path in the `Device` port (stdin, never argv), then normal login.
6. Wire `txn.Manager` Confirm/Rollback/Get endpoints behind auth; call `Recover` at daemon start before serving.
