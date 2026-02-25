#!/bin/bash
. /utils/common.sh

export CGO_ENABLED=1
export PATH=$PATH:$GOPATH/bin


# USAGE:
# build prod|production|debug|dev|optimized-debug [filename]
build() {
    BUILD_OPTIONS="-panic=trap -wasm-abi=generic -scheduler=none -target=wasi -tags=wasi -gc=leaking ${EXTRA_BUILD_OPTIONS}"
    OPTIMIZE=0
    case "${1}" in
        prod|production)
        BUILD_OPTIONS="${BUILD_OPTIONS} --no-debug -opt=2"
        OPTIMIZE=1
        ;;
        debug|dev)
        BUILD_OPTIONS="${BUILD_OPTIONS} -opt=1"
        ;;
        optimized-debug)
        BUILD_OPTIONS="${BUILD_OPTIONS} -opt=2"
        ;;
        *)
        echo  "Failed to parse options"
        return 1
        ;;
    esac

    cd "${SRC}"
    out="${OUT}"
    filename="$2"
    [ "${filename}" != "" ] || filename="./"

    echo "Building ${filename} with taubyte/go-wasm-lib"
            
    #ref: https://github.com/tinygo-org/tinygo/issues/1450
    timeout 300 tinygo build -o "${out}/_artifact.wasm" ${BUILD_OPTIONS} . 2>&1
    ret=$?

    echo -n " * Compile "
    if [ $ret -ne 0 ]
    then 
        echo "[FAILED]"
        exit $ret
    else
        echo "[DONE]"
    fi


    if [ $OPTIMIZE -eq 1 ]
    then
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
    else
        mv "${out}/_artifact.wasm"  "${out}/artifact.wasm"
    fi

    cd -
    return $ret
}
