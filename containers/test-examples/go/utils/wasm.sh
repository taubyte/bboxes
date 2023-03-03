#!/bin/bash

# Uncomment to see commands
# set -x

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

        echo "Building ${filename} with taubyte/go-wasi:test-examples"
        mv .git .git.mv

        mkdir lib
        mv $(find . -type f -not -name "go.mod" -not -name "go.sum" -not -path "./.taubyte/*" -not -path "./taubyte/*" -not -path "./lib/*") lib/

        # TODO INJECT
        # go: finding module for package github.com/mailru/easyjson
        # go: finding module for package github.com/mailru/easyjson/jlexer
        # go: finding module for package github.com/mailru/easyjson/jwriter

        # Grab module name from go.mod
        MODNAME="$(awk '/^module/ { print $2}' go.mod)"
        sed "s/@pkg@/${MODNAME}/g" /utils/_lib_main.go > main.go
        go mod tidy

        (
            echo "Replacing Local SDKs"
            echo "replace bitbucket.org/taubyte/go-ray => ${SDK}/go-ray" >> go.mod
            echo "replace github.com/taubyte/go-sdk-symbols => ${SDK}/go-sdk-symbols" >> go.mod
            echo "replace github.com/taubyte/go-sdk-smartops => ${SDK}/go-sdk-smartops" >> go.mod
            echo "replace github.com/taubyte/go-sdk => ${SDK}/go-sdk" >> go.mod
        )

        # Generate .s files,  need to confirm working
        go run /utils
    
        go generate ./...
        go mod tidy


               
        #ref: https://github.com/tinygo-org/tinygo/issues/1450
        timeout 300 tinygo build -o "${out}/_artifact.wasm" -panic=trap --no-debug -wasm-abi=generic -scheduler=none -target=wasi -tags=wasi -gc=leaking . 2>&1
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
        # Uncomment to display and save wat file
        #  wasm2wat --enable-all  -o "${out}/artifact.wat"  "${out}/artifact.wasm" && (
        #         echo "====== WAT ====="
        #         echo
        #         cat "${out}/artifact.wat"
        #         echo
        #         echo "================"
        #     )
            echo "[DONE]"
        fi


        cd -
        exit $ret
    )
    return $?
}