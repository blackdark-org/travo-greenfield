# 0007 Release tarball, install script, signed self-update

Status: accepted

Not a native package: a self-updating binary and apk/opkg would both own the same files and desync the package database. A tarball holds binary, UI, init script, `uci-defaults` and a `/lib/upgrade/keep.d/` entry (incl. `S99travo` boot link). An install script does first install and verifies the service answers before moving LuCI to :8080/:8443 (reversible on uninstall). Settings live in UCI so backup/sysupgrade cover them.

Update: GitHub release assets `<name>-<arch>.tar.gz`, `checksums.txt`, minisign signature (public key embedded). Download to /tmp, verify, keep `.old`, atomic swap, procd restart, post-restart health check, automatic rollback. No background polling: check on demand and at most daily when the UI opens. Release source configurable; disabled while the repo is private. Version via release-please, stamped with build time. A real `.apk` can come later.
