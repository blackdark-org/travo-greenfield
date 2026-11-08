Captured from `openwrt/rootfs:aarch64_generic-25.12.6` (OpenWrt 25.12.6 r33311-981c978eff) with the `uci` CLI on 2026-11-09.
Container, not hardware: proves format and quoting, not wireless/netifd behaviour.
`uci-show-dhcp.txt` has two list entries on `server` and a value containing an apostrophe, added with `uci add_list` / `uci set` before `uci show`.
