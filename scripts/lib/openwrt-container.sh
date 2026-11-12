#!/bin/sh
# Shared OpenWrt integration-test helpers.

fail() {
  echo "FAIL: $*" >&2
  exit 1
}
step() { echo "== $*"; }
dx() { docker exec "$NAME" sh -c "$1"; }

listeners() { dx "netstat -ltn | awk 'NR>2 {print \$4}' | sort"; }

cleanup() { docker rm -f "$NAME" >/dev/null 2>&1 || true; }

wait_boot() {
  i=0
  until dx 'ubus call service list "{\"name\":\"dnsmasq\"}" | grep -q running' 2>/dev/null; do
    i=$((i + 1))
    [ "$i" -lt 60 ] || fail "container did not finish booting"
    sleep 1
  done
}

wait_health() { # label port
  i=0
  until dx "wget -qO- http://127.0.0.1:$2/api/health" 2>/dev/null | grep -qF "\"version\":\"$VERSION\""; do
    i=$((i + 1))
    [ "$i" -lt 30 ] || fail "$1: /api/health did not report version $VERSION"
    sleep 1
  done
  echo "ok: $1: health reports $VERSION"
}
