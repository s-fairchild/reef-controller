#!/bin/bash

set -o pipefail -o nounset

main() {
	if [[ ! -d build ]]; then
		mkdir build
	fi

    if [[ -z $force ]]; then
        if goFilesChangeCheck "$output_file"; then
            build "$output_file"
        fi
    else
        build "$output_file"
    fi
}

build() {
	tinygo build -target=pico -serial=uart -o "$1"
}

goFilesChangeCheck() {
    readarray -d '' goFiles < <(find . -name "*.go" -print0)

    for f in "${goFiles[@]}"; do
        if [[ $f -nt $1 ]]; then
            build="true"
            break 2
        fi
    done
    echo "${build:-false}"
}

usage() {
    echo "${0} [-f] [-h] -o output_binary

        -f          skip checking for source file changes
        -o          name of output binary file, file is placed under $(pwd)/build
        -h          print this help message"
    exit "$1"
}

force='false'
while getopts 'hfo:' flag; do
    case "${flag}" in
        f)
            force='true'
            ;;
        o)
            output_file="build/${OPTARG}"
            ;;
        h)
            usage 0
            ;;
        *)
            usage 1
            ;;
    esac
done

main "$@"
