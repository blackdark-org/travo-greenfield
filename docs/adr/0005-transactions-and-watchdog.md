# 0005 Self-owned transactions with an out-of-process watchdog

Status: accepted (watchdog behaviour to be proven by spike S4)

rpcd `uci apply` rollback snapshots at apply time and needs a session id and stored root login ([device-facts §4](../device-facts.md)). Travo does not use it.

Flow: guard → snapshot affected `/etc/config/*` to `/tmp` → one commit/apply → a tiny procd-run script restores the snapshot (`uci reload_config` + `network reload`, never `wifi`) unless a confirm file appears within N seconds → per-interface carrier and kernel-route probes → only the user's confirm writes the file. A daemon crash or kill cannot strand the router because the watchdog is a separate process. Guard is removed only after full success.

Consequence: the UI promise is that a bad change reverts itself; "rolled back" is shown only when true.
