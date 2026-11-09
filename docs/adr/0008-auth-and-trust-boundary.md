# 0008 Auth, sessions, trust boundary

Status: accepted

Login validates the root password through rpcd `session.login`. Cookie session (HttpOnly, SameSite=Strict) + CSRF header; SSE for live push on the same cookie (fixes the old token-in-query-string weakness). Memory-only sessions on the monotonic clock; a persisted epoch file kills all older sessions on password change. API tokens for automation come in 0.2.

Admin and first-run are served on `br-lan` (wired + AP WiFi) only, never from uplink or guest zones, enforced by firewall rules AND a daemon source/bind check. First-run with no root password shows only "set password"; the installer also requires `--password` or leaves that path. No default password.

Implementation (`internal/auth`, `internal/trust`, `internal/api`):
- Password check: `ubus call session login` with `timeout:1` (rpcd denies `session destroy`); wrong password and unknown user are the same failure. Password only in the JSON argument (argv-safe) or on stdin for `passwd`, never in logs or responses. NUL in a password is rejected before rpcd (it would truncate).
- First-run: `/etc/shadow` root field empty = not set. Unreadable shadow fails closed (503, setup stays shut). Setup re-checks shadow under a lock; exactly one request wins.
- Sessions: 32-byte token, stored as sha256 only; monotonic expiry; epoch file persisted before the password-change reply; persist failure revokes all and returns 503.
- Trust: outer wrapper before routing, `RemoteAddr` only (no `X-Forwarded-For`), `lan` IPv4 subnets plus loopback; before the first interface read only loopback. **Loopback is trusted**, so a local reverse proxy in front of the daemon would make every client look trusted: packaging must not put a proxy in front, or must drop the loopback allowance.
- Default-deny: every route registers through one helper with an explicit `public` flag; a test walks the spec. Errors `{error, code}`; bad input is 400 `invalid_request` (unknown field, null, trailing data, case-variant keys).
