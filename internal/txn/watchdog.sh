#!/bin/sh
# Transaction watchdog: restores the snapshot unless the transaction is confirmed in time.
# Runs detached from the daemon so a daemon crash or kill cannot strand the router.
# Usage: watchdog.sh <txn-dir> <config-dir> <window-s> <arm-limit-s> <uci> <reload...>
#   The window starts when the daemon writes "armed" (after commit, reload, probe), so the
#   user gets the full window. If "armed" never appears within arm-limit (daemon died
#   mid-apply), restore. window 0 = restore now (explicit rollback, startup recovery).
# Files in <txn-dir>: configs (one name per line), snap/<cfg> (pre-change copy; absent = config
# did not exist), armed, claim (first writer decides: "confirm" or "restore"), result.
# Never calls `wifi` (ath11k crash, see docs/device-facts.md §3).
set -f

dir=$1 cfgdir=$2 secs=$3 armlimit=$4 uci=$5
shift 5 2>/dev/null || exit 2
[ -n "$dir" ] && [ -d "$dir" ] && [ -n "$cfgdir" ] && [ -n "$uci" ] && [ "$#" -gt 0 ] || exit 2

if [ "$secs" -gt 0 ]; then
	i=0
	while [ ! -e "$dir/armed" ] && [ "$i" -lt "$armlimit" ]; do
		[ -e "$dir/claim" ] && exit 0
		sleep 1
		i=$((i + 1))
	done
	i=0
	while [ -e "$dir/armed" ] && [ "$i" -lt "$secs" ]; do
		[ -e "$dir/claim" ] && exit 0
		sleep 1
		i=$((i + 1))
	done
fi

# noclobber makes claiming atomic: the first claimant (confirm or another restorer) wins,
# so exactly one restore ever runs. Losers exit; the daemon waits for "result".
( set -C; echo restore > "$dir/claim" ) 2>/dev/null || exit 0

status=restored
while IFS= read -r cfg; do
	case $cfg in '' | *[!a-z0-9_-]*) status=failed; continue ;; esac
	"$uci" revert "$cfg" >/dev/null 2>&1
	if [ -e "$dir/snap/$cfg" ]; then
		cp "$dir/snap/$cfg" "$cfgdir/.$cfg.restore" && mv -f "$cfgdir/.$cfg.restore" "$cfgdir/$cfg" || status=failed
	else
		rm -f "$cfgdir/$cfg" || status=failed
	fi
done < "$dir/configs"
sync
"$@" >/dev/null 2>&1 || status=failed
echo "$status" > "$dir/result.tmp" && mv -f "$dir/result.tmp" "$dir/result"
sync
[ "$status" = restored ]
