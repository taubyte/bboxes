#/bin/bash

. /utils/common.sh

# USAGE:
# build [filename]
build() {
    (
        cd "${SRC}"
        filename="$1"
        [ "${filename}" != "" ] || filename="./"

        echo "Building ${filename} with rs (wasm32-wasip1)"
        cargo build --target wasm32-wasip1 --release
        ret=$?
        if [ $ret -ne 0 ]; then
            cd -
            exit $ret
        fi

        cp $(echo target/wasm32-wasip1/release/*.wasm | head -1) ${OUT}/_artifact.wasm

        cd -

        echo -n " * Optimize "
        timeout 300 wasm-opt -c -O "${OUT}/_artifact.wasm" -o "${OUT}/artifact.wasm"
        ret=$?
        if [ $ret -ne 0 ]; then
            echo "[FAILED]"
            exit $ret
        else
            echo "[DONE]"
        fi

        exit $ret
    )
    return $?
}
