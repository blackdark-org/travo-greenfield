#!/bin/sh
# Installs the release tarball into a booted OpenWrt container (procd as PID 1), then checks
# health, listeners, reboot survival and uninstall, plus the signature cases of docs/adr/0007.
# Signs with throwaway minisign keys only (never the release key). Needs Docker and minisign.
# Run: mise run test:package
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

# sign_to <secret key> <out>: legacy-format signature of dist/checksums.txt as the router's
# usign reads it (the first two lines of the minisign file).
sign_to() {
  minisign -S -l -s "$1" -m dist/checksums.txt -x "$tmp/full.minisig" >/dev/null
  head -n 2 "$tmp/full.minisig" >"$2"
}

# build_pkg <public key> <dir>: package with that key shipped, extracted into <dir>.
build_pkg() {
  PUBKEY=$1 mise run package
  mkdir "$2"
  tar -xzf "$tarball" -C "$2"
}

step "build the tarball"
tmp=$(mktemp -d)
trap 'cleanup; rm -rf "$tmp"' EXIT
tarball=dist/travo-greenfield-$VERSION-aarch64.tar.gz
minisign -G -W -p "$tmp/k1.pub" -s "$tmp/k1.key" >/dev/null
minisign -G -W -p "$tmp/k2.pub" -s "$tmp/k2.key" >/dev/null
build_pkg "$tmp/k1.pub" "$tmp/pkg"
[ -f "$tarball" ] || fail "$tarball missing"
cmp dist/checksums.txt "$tmp/pkg/checksums.txt" || fail "dist/checksums.txt differs from the one in the tarball"
cmp "$tmp/k1.pub" "$tmp/pkg/minisign.pub" || fail "tarball minisign.pub is not the key given to PUBKEY"
grep -q ' \./minisign.pub$' dist/checksums.txt || fail "minisign.pub not listed in checksums.txt"
# before copying the .sig into pkg: variants that share pkg's files
cp -R "$tmp/pkg" "$tmp/pkg-nosig"
sign_to "$tmp/k1.key" "$tmp/pkg/checksums.txt.sig"
cp -R "$tmp/pkg" "$tmp/pkg-tampered"
echo "0000000000000000000000000000000000000000000000000000000000000000  ./tampered" >>"$tmp/pkg-tampered/checksums.txt"
cp -R "$tmp/pkg" "$tmp/pkg-wrongsig"
sign_to "$tmp/k2.key" "$tmp/pkg-wrongsig/checksums.txt.sig"
# validly signed, plus a file under root/ that checksums.txt does not list
cp -R "$tmp/pkg" "$tmp/pkg-extra"
mkdir -p "$tmp/pkg-extra/root/etc/rc.d"
printf '#!/bin/sh\n' >"$tmp/pkg-extra/root/etc/rc.d/S00evil"
# validly signed, but the listed binary replaced by a symlink to an identical copy outside root/
cp -R "$tmp/pkg" "$tmp/pkg-link"
mv "$tmp/pkg-link/root/usr/bin/travo-greenfield" "$tmp/pkg-link/binary-copy"
ln -s /tmp/tg-link/binary-copy "$tmp/pkg-link/root/usr/bin/travo-greenfield"

step "tarball entries equal root/ plus the manifest"
(cd "$tmp/pkg" && find root -type f | sed 's|^root||' | LC_ALL=C sort) >"$tmp/on-disk.txt"
tar -tzf "$tarball" | grep -v '/$' | grep '^root/' | sed 's|^root||' | LC_ALL=C sort >"$tmp/listed.txt"
cmp "$tmp/on-disk.txt" "$tmp/listed.txt" || fail "tar listing differs from extracted files"
cmp "$tmp/listed.txt" "$tmp/pkg/manifest.txt" || fail "manifest.txt differs from the files under root/"
tar -tzf "$tarball" | LC_ALL=C sort >"$tmp/all.txt"
for f in install.sh VERSION manifest.txt checksums.txt minisign.pub; do
  grep -qx "$f" "$tmp/all.txt" || fail "$f missing from the tarball"
done
echo "ok: $(wc -l <"$tmp/pkg/manifest.txt" | tr -d ' ') manifest paths"

step "boot the container"
docker rm -f "$NAME" >/dev/null 2>&1 || true
docker run -d --name "$NAME" --platform "$PLATFORM" "$IMAGE" /sbin/init >/dev/null
wait_boot
for d in tg:pkg tg-nosig:pkg-nosig tg-tampered:pkg-tampered tg-wrongsig:pkg-wrongsig tg-extra:pkg-extra tg-link:pkg-link; do
  dx "mkdir -p /tmp/${d%%:*}"
  docker cp "$tmp/${d##*:}/." "$NAME:/tmp/${d%%:*}/"
done

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

# Signature cases 1 and 2 (3 is the install below). Each runs install.sh and prints rc=<exit code>.
step "signature: no .sig refuses"
out=$(dx 'sh /tmp/tg-nosig/install.sh 2>&1; echo rc=$?')
echo "$out" | grep -qx 'rc=1' || fail "install without .sig not refused with exit 1: $out"
echo "$out" | grep -q -e '--insecure' || fail "refusal does not mention --insecure: $out"
dx '[ ! -e /usr/bin/travo-greenfield ]' || fail "a refused install left files behind"
echo "ok: refused: $(echo "$out" | head -n 1)"

step "signature: --insecure installs without writing the key"
dx 'sh /tmp/tg-nosig/install.sh --insecure'
dx '[ -x /usr/bin/travo-greenfield ]' || fail "--insecure install did not install"
dx '[ ! -e /etc/travo-greenfield/minisign.pub ]' || fail "--insecure install wrote the trusted key"
dx 'sh /tmp/tg-nosig/install.sh --uninstall --purge'
dx '[ ! -e /usr/bin/travo-greenfield ] && [ ! -e /etc/travo-greenfield ]' || fail "uninstall/purge left files"
echo "ok: --insecure installed, no key written"

step "install"
listeners >"$tmp/listeners-before.txt"
dx 'sh /tmp/tg/install.sh'
wait_health "after install"
dx 'cat /etc/travo-greenfield/minisign.pub' | cmp - "$tmp/k1.pub" || fail "signed install did not store the tarball key"
echo "ok: signed install stored the key"
listeners >"$tmp/listeners-after.txt"
while IFS= read -r l; do
  grep -qxF "$l" "$tmp/listeners-after.txt" || fail "listener $l disappeared after install"
done <"$tmp/listeners-before.txt"
grep -q ':8090$' "$tmp/listeners-after.txt" || fail "nothing listens on :8090: $(cat "$tmp/listeners-after.txt")"
echo "ok: all $(wc -l <"$tmp/listeners-before.txt" | tr -d ' ') pre-install listeners still present"
dx '[ -L /etc/rc.d/S99travo-greenfield ] && [ -L /etc/rc.d/K10travo-greenfield ]' || fail "boot links missing"
dx '[ ! -e /etc/uci-defaults/90-travo-greenfield ]' || fail "uci-defaults script not consumed"
# BusyBox tr has no [:upper:] class: it once turned OpenWrt into wenrt.lan
# shellcheck disable=SC2018,SC2019 # same ASCII-only mapping as the script under test
want=$(dx 'uci -q get system.@system[0].hostname' | tr 'A-Z' 'a-z').lan
got=$(dx 'uci -q get travo-greenfield.main.host_allow')
[ "$got" = "$want" ] || fail "host_allow is '$got', want '$want'"
dx 'sysupgrade -l 2>/dev/null' | grep -q travo-greenfield || fail "sysupgrade -l lists no travo-greenfield path"

step "restart the container (boot link)"
docker restart "$NAME" >/dev/null
wait_boot
wait_health "after reboot"

# Signature case 4: the installed key (k1) is trusted; a tarball with another key needs --rotate-key.
step "signature: key rotation"
build_pkg "$tmp/k2.pub" "$tmp/pkg2"
sign_to "$tmp/k2.key" "$tmp/pkg2/checksums.txt.sig"
dx 'mkdir -p /tmp/tg2'
docker cp "$tmp/pkg2/." "$NAME:/tmp/tg2/"
out=$(dx 'sh /tmp/tg2/install.sh 2>&1; echo rc=$?')
echo "$out" | grep -qx 'rc=1' || fail "different key not refused with exit 1: $out"
echo "$out" | grep -q 'differs from the installed key' || fail "refusal does not say the key differs: $out"
dx 'cat /etc/travo-greenfield/minisign.pub' | cmp - "$tmp/k1.pub" || fail "refused rotation changed the installed key"
dx 'sh /tmp/tg2/install.sh --rotate-key'
dx 'cat /etc/travo-greenfield/minisign.pub' | cmp - "$tmp/k2.pub" || fail "--rotate-key did not store the new key"
wait_health "after key rotation"
echo "ok: different key refused, --rotate-key stored it"

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

# Signature cases 5 and 6 run on the clean container left by the purge: the tarball key is used.
step "signature: .sig from another key refuses"
out=$(dx 'sh /tmp/tg-wrongsig/install.sh 2>&1; echo rc=$?')
echo "$out" | grep -qx 'rc=1' || fail "wrong-key .sig not refused with exit 1: $out"
echo "$out" | grep -q 'signature check failed' || fail "refusal is not 'signature check failed': $out"
dx '[ ! -e /usr/bin/travo-greenfield ]' || fail "a refused install left files behind"
echo "ok: refused: $(echo "$out" | head -n 1)"

step "signature: tampered checksums.txt refuses"
out=$(dx 'sh /tmp/tg-tampered/install.sh 2>&1; echo rc=$?')
echo "$out" | grep -qx 'rc=1' || fail "tampered checksums.txt not refused with exit 1: $out"
echo "$out" | grep -q 'signature check failed' || fail "refusal is not 'signature check failed': $out"
dx '[ ! -e /usr/bin/travo-greenfield ] && [ ! -e /etc/travo-greenfield ]' || fail "a refused install left files behind"
echo "ok: refused: $(echo "$out" | head -n 1)"

step "signature: unlisted file under root/ refuses"
out=$(dx 'sh /tmp/tg-extra/install.sh 2>&1; echo rc=$?')
echo "$out" | grep -qx 'rc=1' || fail "unlisted file not refused with exit 1: $out"
echo "$out" | grep -q 'not all listed' || fail "refusal does not name the unlisted files: $out"
dx '[ ! -e /etc/rc.d/S00evil ] && [ ! -e /usr/bin/travo-greenfield ]' || fail "a refused install left files behind"
echo "ok: refused: $(echo "$out" | head -n 1)"

step "signature: symlink under root/ refuses"
out=$(dx 'sh /tmp/tg-link/install.sh 2>&1; echo rc=$?')
echo "$out" | grep -qx 'rc=1' || fail "symlink not refused with exit 1: $out"
echo "$out" | grep -q 'other than regular files' || fail "refusal does not name the symlink: $out"
dx '[ ! -e /usr/bin/travo-greenfield ]' || fail "a refused install left files behind"
echo "ok: refused: $(echo "$out" | head -n 1)"

step "package test passed"
