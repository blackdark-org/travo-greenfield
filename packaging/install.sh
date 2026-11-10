#!/bin/sh
# Installs or removes travo-greenfield on OpenWrt. Run as root from the extracted tarball:
#   tar -xzf travo-greenfield-<version>-aarch64.tar.gz -C /tmp/tg && sh /tmp/tg/install.sh
# Flags: (none) install | --uninstall | --purge (with --uninstall: also config and state)
#        | --force (skip the OS and architecture checks only)
#        | --insecure (install without a signature check) | --rotate-key (accept a new signing key).
# Before installing it verifies $here/checksums.txt.sig over checksums.txt with usign (docs/adr/0007).
# Touches only the paths in manifest.txt and the ones named in uninstall() (invariants §15).
set -eu

NAME=travo-greenfield
STATE_DIR=/etc/travo-greenfield
MIN_FREE_KIB=20480

die() {
  echo "install.sh: $*" >&2
  exit 1
}

usage() {
  echo "usage: install.sh [--force] [--insecure] [--rotate-key] | --uninstall [--purge]" >&2
  exit 2
}

# Major.minor of DISTRIB_RELEASE must be >= 25.12; SNAPSHOT builds pass.
check_os() {
  release=$(sed -n "s/^DISTRIB_RELEASE=['\"]\{0,1\}\([^'\"]*\)['\"]\{0,1\}\$/\1/p" /etc/openwrt_release 2>/dev/null) ||
    release=""
  [ -n "$release" ] || die "cannot read DISTRIB_RELEASE from /etc/openwrt_release (not OpenWrt? use --force)"
  case "$release" in *SNAPSHOT*) return 0 ;; esac
  major=${release%%.*}
  rest=${release#*.}
  minor=${rest%%.*}
  minor=${minor%%[!0-9]*}
  case "$major$minor" in
    "" | *[!0-9]*) die "unrecognised OpenWrt release '$release' (use --force)" ;;
  esac
  if [ "$major" -lt 25 ] || { [ "$major" -eq 25 ] && [ "${minor:-0}" -lt 12 ]; }; then
    die "OpenWrt $release is not supported, need 25.12 or newer (use --force)"
  fi
}

check_arch() {
  arch=$(uname -m)
  [ "$arch" = aarch64 ] || die "architecture $arch is not supported, need aarch64 (use --force)"
}

check_space() {
  target=/overlay
  [ -d "$target" ] || target=/
  free=$(df -k "$target" | tail -n 1 | awk '{print $4}')
  case "$free" in "" | *[!0-9]*) die "cannot read free space of $target" ;; esac
  [ "$free" -ge "$MIN_FREE_KIB" ] || die "only $free KiB free on $target, need $MIN_FREE_KIB"
}

# Polls /api/health once a second (BusyBox has no timeout or curl) until it reports $1.
wait_healthy() {
  listen=$(uci -q get "$NAME.main.listen") || die "no $NAME.main.listen in UCI"
  port=${listen##*:}
  i=0
  while [ "$i" -lt 30 ]; do
    if wget -q -O - "http://127.0.0.1:$port/api/health" 2>/dev/null | grep -qF "\"version\":\"$1\""; then
      return 0
    fi
    i=$((i + 1))
    sleep 1
  done
  return 1
}

# Verifies checksums.txt against checksums.txt.sig with a key the router already trusts: the
# installed one, or on the first install the one in the tarball (trust on first use).
# Sets new_key=1 when the verified key is the tarball's and must be saved after the install.
verify_signature() {
  new_key=0
  key=$here/minisign.pub
  installed=$STATE_DIR/minisign.pub
  if [ -f "$installed" ]; then
    key=$installed
    if [ -f "$here/minisign.pub" ] && ! cmp -s "$installed" "$here/minisign.pub"; then
      [ "$rotate" = 1 ] || die "tarball key differs from the installed key; re-run with --rotate-key if this rotation is intended"
      key=$here/minisign.pub
      new_key=1
    fi
  elif [ -f "$key" ]; then
    new_key=1
  fi
  if [ -f "$key" ] && [ -f "$here/checksums.txt.sig" ]; then
    command -v usign >/dev/null 2>&1 || die "usign not found"
    usign -V -q -p "$key" -m "$here/checksums.txt" -x "$here/checksums.txt.sig" || die "signature check failed"
    return 0
  fi
  new_key=0
  [ "$insecure" = 1 ] || die "unsigned: no checksums.txt.sig or minisign.pub; re-run with --insecure to install without a signature check"
  if [ -f "$installed" ]; then
    echo "install.sh: WARNING: --insecure although a signing key is installed: this tarball is NOT verified against it" >&2
  else
    echo "install.sh: WARNING: --insecure: installing without a signature check" >&2
  fi
}

install() {
  [ "$force" = 1 ] || { check_os; check_arch; }
  check_space
  cd "$here"
  [ -f checksums.txt ] && [ -d root ] || die "$here is not an extracted release tarball"
  verify_signature
  sha256sum -c checksums.txt >/dev/null || die "checksum mismatch in $here (run sha256sum -c checksums.txt)"
  # Only regular files and directories: a symlink passes sha256sum (it follows it) but cp -a
  # would install the link itself.
  if [ -n "$(find root ! -type f ! -type d | head -n 1)" ]; then
    die "$here/root contains something other than regular files and directories"
  fi
  # sha256sum -c checks only listed files; cp copies all of root/. Refuse anything unlisted.
  find root ! -type d | sed 's|^|./|' | LC_ALL=C sort >"$here/.on-disk"
  sed -n 's|^[0-9a-f]\{64\}  \(\./root/.*\)$|\1|p' checksums.txt | LC_ALL=C sort >"$here/.listed"
  if ! cmp -s "$here/.on-disk" "$here/.listed"; then
    rm -f "$here/.on-disk" "$here/.listed"
    die "files under $here/root are not all listed in checksums.txt"
  fi
  rm -f "$here/.on-disk" "$here/.listed"
  version=$(cat VERSION)

  cp -a root/. / || die "copying files to / failed"
  # OpenWrt runs uci-defaults only at boot; run it now. It deletes itself on success.
  defaults=/etc/uci-defaults/90-$NAME
  if [ -e "$defaults" ]; then
    sh "$defaults" || die "$defaults failed"
    rm -f "$defaults"
  fi
  mkdir -p "$STATE_DIR"
  cp manifest.txt "$STATE_DIR/manifest.txt" || die "cannot save the manifest to $STATE_DIR"
  # only a verified install may establish or rotate the trusted key
  if [ "$new_key" = 1 ]; then
    cp "$here/minisign.pub" "$STATE_DIR/minisign.pub" || die "cannot save the signing key to $STATE_DIR"
  fi

  # restart, not start: an older binary may be running
  /etc/init.d/$NAME restart || die "/etc/init.d/$NAME restart failed"
  if ! wait_healthy "$version"; then
    echo "install.sh: $NAME $version did not answer on /api/health within 30 s; files stay installed. Last log lines:" >&2
    logread -e "$NAME" | tail -20 >&2
    exit 1
  fi
  echo "$NAME $version installed and healthy on port $port"
}

uninstall() {
  manifest=$STATE_DIR/manifest.txt
  [ -f "$manifest" ] || die "$manifest is missing: nothing to uninstall (or it was purged)"
  # after a plain uninstall the init script is gone; --purge then only cleans config and state
  if [ -x /etc/init.d/$NAME ]; then
    /etc/init.d/$NAME stop || true
    /etc/init.d/$NAME disable || true
  fi
  while IFS= read -r path; do
    # Only paths this package can own, and only files: a tampered manifest cannot reach others.
    case "$path" in
      *..*)
        echo "install.sh: skipping $path: contains .." >&2
        continue
        ;;
      "/usr/bin/$NAME" | "/etc/init.d/$NAME" | "/etc/uci-defaults/90-$NAME" | "/lib/upgrade/keep.d/$NAME" | "/www/$NAME/"?*) ;;
      *)
        echo "install.sh: skipping $path: not a $NAME path" >&2
        continue
        ;;
    esac
    if [ -f "$path" ] || [ -L "$path" ]; then
      rm -f "$path"
    elif [ -e "$path" ]; then
      echo "install.sh: skipping $path: not a file" >&2
    fi
    # empty directories below the web root go too, innermost first
    dir=$(dirname "$path")
    while :; do
      case "$dir" in /www/$NAME/*) ;; *) break ;; esac
      rmdir "$dir" 2>/dev/null || break
      dir=$(dirname "$dir")
    done
  done <"$manifest"
  rm -f /etc/rc.d/S99$NAME /etc/rc.d/K10$NAME
  rmdir /www/$NAME 2>/dev/null || true
  if [ "$purge" = 1 ]; then
    rm -f /etc/config/$NAME
    rm -rf "$STATE_DIR" /etc/travo-greenfield-guards
  fi
  echo "$NAME removed"
}

force=0 purge=0 insecure=0 rotate=0 new_key=0 mode=install
for arg in "$@"; do
  case "$arg" in
    --uninstall) mode=uninstall ;;
    --purge) purge=1 ;;
    --force) force=1 ;;
    --insecure) insecure=1 ;;
    --rotate-key) rotate=1 ;;
    *) usage ;;
  esac
done
[ "$purge" = 0 ] || [ "$mode" = uninstall ] || usage
[ "$mode" = install ] || { [ "$insecure" = 0 ] && [ "$rotate" = 0 ]; } || usage

[ "$(id -u)" = 0 ] || die "must run as root"
here=$(cd "$(dirname "$0")" && pwd)
"$mode"
