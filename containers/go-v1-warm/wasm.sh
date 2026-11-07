#!/bin/bash
. /utils/common.sh

export CGO_ENABLED=1
export PATH=$PATH:$GOPATH/bin


# USAGE:
# build [filename]
build() {
    (
        cd "${SRC}"
        out="${OUT}"
        filename="$1"
        [ "${filename}" != "" ] || filename="./"

        echo "Building ${filename} with taubyte/go-wasm"
        mv .git .git.mv

        mkdir lib
        mv $(find . -type f -not -name "go.mod" -not -name "go.sum" -not -path "./.taubyte/*" -not -path "./taubyte/*" -not -path "./lib/*") lib/

        MODNAME="$(awk '/^module/ { print $2}' go.mod)"

        sed "s/@pkg@/${MODNAME}/g" /utils/_lib_main.go_ > main.go
        go mod tidy

        if grep -rqs --include=*.go "^//go:generate" .; then
            go generate ./...
            go mod tidy
        fi
        
               
        #ref: https://github.com/tinygo-org/tinygo/issues/1450
        timeout 300 tinygo build -o "${out}/_artifact.wasm" -panic=trap --no-debug -scheduler=none -target=wasi -tags=wasi -gc=leaking . 2>&1
        ret=$?

        echo -n " * Compile "
        if [ $ret -ne 0 ]
        then 
            echo "[FAILED]"
            exit $ret
        else
            echo "[DONE]"
        fi

        timeout 300 wasm-opt -c -O "${out}/_artifact.wasm" -o "${out}/artifact.wasm"
        ret=$?

        echo -n " * Optimize "
        if [ $ret -ne 0 ]
        then 
            echo "[FAILED]"
            exit $ret
        else
            echo "[DONE]"
        fi


        cd -
        exit $ret
    )
    return $?
}

# USAGE:
# debug_build [opt] [filename]
debug_build() {
    (
        cd "${SRC}"
        out="${OUT}"
        opt="${1:-1}"
        filename="$2"
        [ "${filename}" != "" ] || filename="./"

        echo "Building ${filename} with taubyte/go-wasm"
        mv .git .git.mv

        mkdir lib
        mv $(find . -type f -not -name "go.mod" -not -name "go.sum" -not -path "./.taubyte/*" -not -path "./taubyte/*" -not -path "./lib/*") lib/

        MODNAME="$(awk '/^module/ { print $2}' go.mod)"

        sed "s/@pkg@/${MODNAME}/g" /utils/_lib_main.go_ > main.go
        go mod tidy

        if grep -rqs --include=*.go "^//go:generate" .; then
            go generate ./...
            go mod tidy
        fi
        
               
        #ref: https://github.com/tinygo-org/tinygo/issues/1450
        timeout 300 tinygo build -o "${out}/artifact.wasm" -opt=${opt} -panic=trap -wasm-abi=generic -scheduler=none -target=wasi -tags=wasi -gc=leaking . 2>&1
        ret=$?

        echo -n " * Compile "
        if [ $ret -ne 0 ]
        then 
            echo "[FAILED]"
            exit $ret
        else
            echo "[DONE]"
        fi

        cd -
        exit $ret
    )
    return $?
}