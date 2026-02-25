#!/bin/bash
. /utils/common.sh

export CGO_ENABLED=1
export PATH=$PATH:$GOPATH/bin
# Avoid "GOPROXY list is not the empty string, but contains no entries" when env has a broken value
export GOPROXY=https://proxy.golang.org,direct

# TinyGo 0.40+ supports -wasm-abi=generic; older images (e.g. go-wasi v3.x) may not.
# Detect support so the same script works in both cases.
wasm_abi_flag() {
    if tinygo build -h 2>&1 | grep -q '\-wasm-abi'; then
        echo "-wasm-abi=generic"
    else
        echo ""
    fi
}

# USAGE:
# build [filename]
# Reactor-style WASM: -buildmode=c-shared so TinyGo emits _initialize (not _start). Host must call _initialize then //export functions.
# Build runs in a temp copy of SRC so /src is never modified; rewrite-package converts package lib (or other) to main in the copy.
build() {
    (
        out="${OUT}"
        filename="$1"
        [ "${filename}" != "" ] || filename="./"
        wasm_abi=$(wasm_abi_flag)

        tmp=$(mktemp -d)
        trap "rm -rf ${tmp}" EXIT
        cp -a "${SRC}/." "${tmp}/"

        if ! /utils/rewrite-package "${tmp}"; then
            echo " * Rewrite package [FAILED]"
            exit 1
        fi

        cd "${tmp}"
        echo "Building ${filename} with taubyte/go-wasi (reactor-style WASM)"
        mv .git .git.mv 2>/dev/null || true

        go mod tidy

        timeout 300 tinygo build -buildmode=c-shared -o "${out}/_artifact.wasm" -panic=trap --no-debug ${wasm_abi} -scheduler=none -target=wasi -tags=wasi -gc=leaking . 2>&1
        ret=$?

        echo -n " * Compile "
        if [ $ret -ne 0 ]; then
            echo "[FAILED]"
            exit $ret
        else
            echo "[DONE]"
        fi

        timeout 300 wasm-opt -c -O "${out}/_artifact.wasm" -o "${out}/artifact.wasm"
        ret=$?

        echo -n " * Optimize "
        if [ $ret -ne 0 ]; then
            echo "[FAILED]"
            exit $ret
        else
            echo "[DONE]"
        fi

        mv .git.mv .git 2>/dev/null || true
        exit $ret
    )
    return $?
}

# USAGE:
# debug_build [opt] [filename]
debug_build() {
    (
        out="${OUT}"
        opt="${1:-1}"
        filename="$2"
        [ "${filename}" != "" ] || filename="./"

        tmp=$(mktemp -d)
        trap "rm -rf ${tmp}" EXIT
        cp -a "${SRC}/." "${tmp}/"

        if ! /utils/rewrite-package "${tmp}"; then
            echo " * Rewrite package [FAILED]"
            exit 1
        fi

        cd "${tmp}"
        echo "Building ${filename} with taubyte/go-wasi (reactor-style WASM, c-shared, debug)"
        mv .git .git.mv 2>/dev/null || true

        go mod tidy

        wasm_abi=$(wasm_abi_flag)
        timeout 300 tinygo build -buildmode=c-shared -o "${out}/artifact.wasm" -opt="${opt}" -panic=trap ${wasm_abi} -scheduler=none -target=wasi -tags=wasi -gc=leaking . 2>&1
        ret=$?

        echo -n " * Compile "
        if [ $ret -ne 0 ]; then
            echo "[FAILED]"
            exit $ret
        else
            echo "[DONE]"
        fi

        mv .git.mv .git 2>/dev/null || true
        exit $ret
    )
    return $?
}
