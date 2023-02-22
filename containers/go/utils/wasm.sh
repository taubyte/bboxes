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

        echo "Building ${filename} with taubyte/go-was:v0"
        mv .git .git.mv

        mkdir lib
        mv *.go lib/

        MODNAME="$(awk '/^module/ { print $2}' go.mod)"

        sed "s/@pkg@/${MODNAME}/g" /utils/_lib_main.go > main.go
        go mod tidy

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
            echo "[DONE]"
        fi


        cd -
        exit $ret
    )
    return $?
}