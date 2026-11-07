#!/bin/bash
# Fill GOPATH/pkg/mod and the TinyGo/Go build caches by compiling a stub that imports every go-sdk package.
set -e
. /utils/common.sh
mkdir -p /tmp/stub && cd /tmp/stub
printf 'module function\n\ngo 1.19\n' > go.mod
go get github.com/taubyte/go-sdk@latest
sdk=$(go list -m -f '{{.Dir}}' github.com/taubyte/go-sdk)
{
  echo 'package lib'; echo 'import ('
  (cd "$sdk" && find . -name '*.go' ! -name '*_test.go' ! -path '*/internal/*' ! -path './examples/*' ! -path './vendor/*' -printf '%h\n' | sort -u | while read -r d; do grep -qsl "^package main" "$d"/*.go || echo "$d"; done | sed 's|^\.$||;s|^\./||') \
    | while read -r p; do echo "  _ \"github.com/taubyte/go-sdk${p:+/$p}\""; done
  echo ')'; echo '//export warm'; echo 'func warm() uint32 { return 0 }'
} > /src/fn.go
cp go.mod /src/go.mod
mkdir -p /src/.git; cd /src
. /utils/wasm.sh
set +e
build ""
test -s /out/artifact.wasm
