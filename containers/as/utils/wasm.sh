#!/bin/bash

. /utils/common.sh

# USAGE:
# build [filename]
build() {
    (
        cd "${SRC}"
        out="${OUT}"
        filename="$1"
        [ "${filename}" != "" ] || filename="./"

        echo "Building ${filename} with as-v1"
    
        mkdir temp
        mv *.ts temp

        cp -r /utils/as-dir .
        mkdir as-dir/assembly/lib
        mv temp/*.ts as-dir/assembly/lib

        # write an index.ts which has a line for each item in as-dir/assembly/lib
        # export * from $ <- without the extention
        echo "" > as-dir/assembly/index.ts
        for filename in as-dir/assembly/lib/*.ts; do
            echo "export * from \"./lib/$(basename $filename .ts)\";" >> as-dir/assembly/index.ts
        done

        cd -


        cd ${SRC}/as-dir
        
        npm i
        npm update sdk
        npm run build

        cp build/artifact.wasm ${OUT}/artifact.wasm

        cd -

        exit $ret
    )
    return $?
}