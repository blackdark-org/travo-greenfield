#!/bin/sh
# Installs the release tarball into a booted OpenWrt container (procd as PID 1), then checks
# health, listeners, reboot survival and uninstall. Needs Docker. Run: mise run test:package
set -eu
cd "$(dirname "$0")/.."

IMAGE=openwrt/rootfs:aarch64_generic-25.12.6@sha256:22f13387da208c0fb4150c41f3ec12e34cc5994a0717a49e50aa9f9ae6cb5e7b
PLATFORM=linux/aarch64_generic
NAME=${CONTAINER_NAME:-travo-pkgtest-$$}
VERSION=${VERSION:-0.0.0-dev}
export VERSION

fail() {
  echo "FAIL: $*" >&2
  exit 1
}
step() { echo "== $*"; }
dx() { docker exec "$NAME" sh -c "$1"; }

listeners() { dx "netstat -ltn | awk 'NR>2 {print \$4}' | sort"; }

cleanup() { docker rm -f "$NAME" >/dev/null 2>&1 || true; }
trap cleanup EXIT

wait_boot() {
  i=0
  until dx 'ubus call service list "{\"name\":\"dnsmasq\"}" | grep -q running' 2>/dev/null; do
    i=$((i + 1))
    [ "$i" -lt 60 ] || fail "container did not finish booting"
    sleep 1
  done
}

wait_health() { # label
  i=0
  until dx 'wget -qO- http://127.0.0.1:8090/api/health' 2>/dev/null | grep -qF "\"version\":\"$VERSION\""; do
    i=$((i + 1))
    [ "$i" -lt 30 ] || fail "$1: /api/health did not report version $VERSION"
    sleep 1
  done
  echo "ok: $1: health reports $VERSION"
}

step "build the tarball"
mise run package
tarball=dist/travo-greenfield-$VERSION-aarch64.tar.gz
[ -f "$tarball" ] || fail "$tarball missing"

step "tarball entries equal root/ plus the manifest"
tmp=$(mktemp -d)
trap 'cleanup; rm -rf "$tmp"' EXIT
mkdir "$tmp/pkg"
tar -xzf "$tarball" -C "$tmp/pkg"
(cd "$tmp/pkg" && find root -type f | sed 's|^root||' | LC_ALL=C sort) >"$tmp/on-disk.txt"
tar -tzf "$tarball" | grep -v '/$' | grep '^root/' | sed 's|^root||' | LC_ALL=C sort >"$tmp/listed.txt"
cmp "$tmp/on-disk.txt" "$tmp/listed.txt" || fail "tar listing differs from extracted files"
cmp "$tmp/listed.txt" "$tmp/pkg/manifest.txt" || fail "manifest.txt differs from the files under root/"
tar -tzf "$tarball" | LC_ALL=C sort >"$tmp/all.txt"
for f in install.sh VERSION manifest.txt checksums.txt; do
  grep -qx "$f" "$tmp/all.txt" || fail "$f missing from the tarball"
done
echo "ok: $(wc -l <"$tmp/pkg/manifest.txt" | tr -d ' ') manifest paths"

step "boot the container"
docker rm -f "$NAME" >/dev/null 2>&1 || true
docker run -d --name "$NAME" --platform "$PLATFORM" "$IMAGE" /sbin/init >/dev/null
wait_boot
dx 'mkdir -p /tmp/tg'
docker cp "$tmp/pkg/." "$NAME:/tmp/tg/"

step "refusals before the real install"
dx 'cp /etc/openwrt_release /tmp/openwrt_release.real'
out=$(dx "sed -i \"s/^DISTRIB_RELEASE=.*/DISTRIB_RELEASE='24.10.8'/\" /etc/openwrt_release; sh /tmp/tg/install.sh 2>&1; echo rc=\$?")
dx 'cp /tmp/openwrt_release.real /etc/openwrt_release'
echo "$out" | grep -qx 'rc=1' || fail "old release not refused with exit 1: $out"
echo "$out" | grep -q '24.10.8' || fail "refusal does not name the release: $out"
echo "ok: 24.10.8 refused: $(echo "$out" | head -n 1)"
out=$(docker exec -u 65534 "$NAME" sh -c 'sh /tmp/tg/install.sh 2>&1; echo rc=$?')
echo "$out" | grep -qx 'rc=1' || fail "non-root install not refused with exit 1: $out"
echo "$out" | grep -q 'root' || fail "non-root refusal has no reason: $out"
echo "ok: non-root refused: $(echo "$out" | head -n 1)"
dx '[ ! -e /usr/bin/travo-greenfield ]' || fail "a refused install left files behind"

step "install"
listeners >"$tmp/listeners-before.txt"
dx 'sh /tmp/tg/install.sh'
wait_health "after install"
listeners >"$tmp/listeners-after.txt"
while IFS= read -r l; do
  grep -qxF "$l" "$tmp/listeners-after.txt" || fail "listener $l disappeared after install"
done <"$tmp/listeners-before.txt"
grep -q ':8090$' "$tmp/listeners-after.txt" || fail "nothing listens on :8090: $(cat "$tmp/listeners-after.txt")"
echo "ok: all $(wc -l <"$tmp/listeners-before.txt" | tr -d ' ') pre-install listeners still present"
dx '[ -L /etc/rc.d/S99travo-greenfield ] && [ -L /etc/rc.d/K10travo-greenfield ]' || fail "boot links missing"
dx '[ ! -e /etc/uci-defaults/90-travo-greenfield ]' || fail "uci-defaults script not consumed"
dx 'sysupgrade -l 2>/dev/null' | grep -q travo-greenfield || fail "sysupgrade -l lists no travo-greenfield path"

step "restart the container (boot link)"
docker restart "$NAME" >/dev/null
wait_boot
wait_health "after reboot"

step "uninstall keeps config and ignores foreign manifest paths"
# A tampered manifest must not reach files the package does not own.
dx 'printf "%s\n" /etc/config/network /etc /www/travo-greenfield/../../etc/config/firewall >> /etc/travo-greenfield/manifest.txt'
dx 'sh /tmp/tg/install.sh --uninstall'
dx '[ -f /etc/config/network ] && [ -f /etc/config/firewall ] && [ -d /etc ]' || fail "uninstall followed a foreign manifest path"
while IFS= read -r p; do
  # the manifest copy lives under /etc/travo-greenfield/ and is not in the manifest
  dx "[ ! -e '$p' ]" || fail "$p still exists after uninstall"
done <"$tmp/pkg/manifest.txt"
dx '[ ! -e /etc/rc.d/S99travo-greenfield ] && [ ! -e /etc/rc.d/K10travo-greenfield ]' || fail "boot links remain"
dx '[ -e /etc/config/travo-greenfield ]' || fail "uninstall removed /etc/config/travo-greenfield"
dx '[ ! -e /www/travo-greenfield ]' || fail "/www/travo-greenfield remains"
echo "ok: uninstalled, config kept"

step "uninstall --purge"
dx 'sh /tmp/tg/install.sh --uninstall --purge'
dx '[ ! -e /etc/config/travo-greenfield ] && [ ! -e /etc/travo-greenfield ]' || fail "purge left config or state"
echo "ok: purge removed config and state"

step "package test passed"
