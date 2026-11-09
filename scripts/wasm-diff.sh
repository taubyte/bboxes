#!/bin/bash
# Usage: wasm-diff.sh IMAGE_A IMAGE_B [RUNS=1]
# Builds every corpus dir with both images (same build.sh as the tau Go template:
# /src rw, /out), compares artifact.wasm sha256 across runs and images, prints
# median build seconds per image. Every build is a fresh container (cold caches).
set -u
[ $# -ge 2 ] || { echo "usage: wasm-diff.sh IMAGE_A IMAGE_B [RUNS=1]" >&2; exit 2; }
command -v bc >/dev/null || { echo "wasm-diff.sh needs bc" >&2; exit 2; }
A=$1; B=$2; RUNS=${3:-1}
here=$(cd "$(dirname "$0")" && pwd)
work=$(mktemp -d); logs=$(mktemp -d)
# the containers run as root and leave root-owned files in $work: clean up from inside one
trap 'docker run --rm -v "$work:/w" "$A" bash -c "rm -rf /w/* /w/.[!.]*" >/dev/null 2>&1; rm -rf "$work"' EXIT
fail=0

build() { # image corpus run label -> "sha seconds"
  local d=$work/$4-$2-$3; rm -rf "$d"; mkdir -p "$d/src/.taubyte" "$d/out"
  cp -a "$here/corpus/$2/." "$d/src/"
  cp "$here/build.sh.tmpl" "$d/src/.taubyte/build.sh"
  local t0 t1; t0=$(date +%s.%N)
  docker run --rm -v "$d/src:/src" -v "$d/out:/out" "$1" bash /src/.taubyte/build.sh >"$d/log" 2>&1
  t1=$(date +%s.%N)
  if [ -f "$d/out/artifact.wasm" ]; then
    echo "$(sha256sum "$d/out/artifact.wasm" | cut -c1-16) $(echo "$t1 - $t0" | bc)"
  else
    cp "$d/log" "$logs/$4-$2-$3.log"; echo "build failed, log: $logs/$4-$2-$3.log" >&2; echo "FAILED 0"
  fi
}
median() { sort -n | awk '{a[NR]=$1} END{print a[int((NR+1)/2)]}'; }

printf "%-12s %-17s %-17s %s\n" corpus shaA shaB result
for c in $(ls "$here/corpus"); do
  ref=""; ok=1; ta=(); tb=()
  for r in $(seq "$RUNS"); do
    read -r sa t1 <<<"$(build "$A" "$c" "$r" A)"; read -r sb t2 <<<"$(build "$B" "$c" "$r" B)"
    [ -z "$ref" ] && ref=$sa
    [ "$sa" != FAILED ] && [ "$sa" = "$ref" ] && [ "$sb" = "$ref" ] || ok=0
    ta+=("$t1"); tb+=("$t2")
  done
  printf "%-12s %-17s %-17s %s  A=%ss B=%ss\n" "$c" "$ref" "$sb" "$([ $ok = 1 ] && echo SAME || echo DIFF)" \
    "$(printf '%s\n' "${ta[@]}" | median)" "$(printf '%s\n' "${tb[@]}" | median)"
  [ $ok = 1 ] || fail=1
done
[ $fail = 0 ] && echo PASS || { echo FAIL; exit 1; }
