# 0008 Auth, sessions, trust boundary

Status: accepted

Login validates the root password through rpcd `session.login`. Cookie session (HttpOnly, SameSite=Strict) + CSRF header; SSE for live push on the same cookie (fixes the old token-in-query-string weakness). Memory-only sessions on the monotonic clock; a persisted epoch file kills all older sessions on password change. API tokens for automation come in 0.2.

Admin and first-run are served on `br-lan` (wired + AP WiFi) only, never from uplink or guest zones, enforced by firewall rules AND a daemon source/bind check. First-run with no root password shows only "set password"; the installer also requires `--password` or leaves that path. No default password.
