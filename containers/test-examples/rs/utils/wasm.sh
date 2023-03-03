#/bin/bash

. /utils/common.sh

# USAGE:
# build [filename]
build() {
    (
        cd "${SRC}"
        filename="$1"
        [ "${filename}" != "" ] || filename="./"

        echo "Building ${filename} with as-v1"
        mv *.rs lib.rs
    
        cargo wasi build --release

        cp $(echo target/wasm32-wasi/release/*.wasi.wasm | head -1) ${OUT}/artifact.wasm

        #find . -name "*.wasi.wasm" -exec cp {} ${OUT}/artifact.wasm \;

        cd -

        exit $ret
    )
    return $?
}
