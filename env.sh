#!/usr/bin/env bash

export APPROOT=$(pwd)

# adjust GOPATH
case ":$GOPATH:" in
    *":$APPROOT:"*) :;;
    *) GOPATH=$APPROOT:$GOPATH;;
esac
export GOPATH


# adjust PATH
if [ -n "$ZSH_VERSION" ]; then
    readopts="rA"
else
    readopts="ra"
fi
while IFS=':' read -$readopts ARR; do
    for i in "${ARR[@]}"; do
        case ":$PATH:" in
            *":$i/bin:"*) :;;
            *) PATH=$i/bin:$PATH
        esac
    done
done <<< "$GOPATH"
export PATH


# mock development && test envs
if [ ! -d "$APPROOT/src/github.com/dolab/bench" ]; then
    mkdir -p "$APPROOT/src/github.com/dolab"
    ln -s "$APPROOT/bench/" "$APPROOT/src/github.com/dolab/bench"
fi
