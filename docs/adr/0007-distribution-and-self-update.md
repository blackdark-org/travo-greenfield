# 0007 Release tarball, install script, signed self-update

Status: accepted

Not a native package: a self-updating binary and apk/opkg would both own the same files and desync the package database. A tarball holds binary, UI, init script, `uci-defaults` and a `/lib/upgrade/keep.d/` entry (incl. `S99travo-greenfield` boot link). The installer targets a fresh OpenWrt installation and verifies the service answers before moving LuCI to :8080/:8443 (reversible on uninstall). Only `install.sh` moves the uhttpd ports. `internal/txn` has no uhttpd path. The install and uninstall port changes use a guard and a uhttpd snapshot. Restoration preserves each section’s listen values, their order, and whether they were an option or list. Settings live in UCI so backup/sysupgrade cover them.

Update: GitHub release assets `<name>-<arch>.tar.gz`, `checksums.txt`, minisign signature (public key embedded). Download to /tmp, verify, keep `.old`, atomic swap, procd restart, post-restart health check, automatic rollback. No background polling: check on demand and at most daily when the UI opens. Release source configurable; disabled while the repo is private. Version via release-please, stamped with build time. A real `.apk` can come later.

## Signing key

- One minisign key pair signs `checksums.txt` of every release, in the legacy format (`minisign -S -l`) so the router's `usign` verifies it (OpenWrt ships `usign`, not minisign). `usign` reads only the first two lines of the `.minisig` (untrusted comment and signature); the release ships them as `checksums.txt.sig`. The trusted comment is not checked by `usign` and must not be relied on.
- The secret key lives only in the GitHub secrets `MINISIGN_SECRET_KEY` and `MINISIGN_PASSWORD`, plus the operator's offline backup. The public key is `packaging/minisign.pub`, shipped in the tarball as `minisign.pub`. Without that file builds are unsigned and `release.yml` fails before uploading anything.
- Trust model: a key inside the tarball proves integrity only. `install.sh` verifies against `/etc/travo-greenfield/minisign.pub` when it exists (kept by sysupgrade through the `/etc/travo-greenfield/` keep.d entry). On the first install (no key on the router) it verifies against the tarball's key: trust on first use, the operator fetched both over HTTPS from the GitHub release. A tarball with a different key is refused unless `--rotate-key` is given. `--insecure` installs without a signature check and never writes the key. Self-update embeds the key in the binary instead.
- Rotation: generate a new key pair, commit the new `packaging/minisign.pub`, update the secrets; the release notes tell users to pass `--rotate-key` once.

