Container captures, not hardware (format only).
- `apk-query-installed.json`: `apk query --installed --fields name,version --format json "*"` in `openwrt/rootfs:aarch64_generic-25.12.6` (apk-tools 3.0.5), 2026-11-09.
- `opkg-list-installed.txt`: `opkg list-installed` in `openwrt/rootfs:x86-64-24.10.8` (opkg 38eccbb, 2024-10-16), 2026-11-09. Image is amd64 because no arm64 24.10 manifest exists for the host.
