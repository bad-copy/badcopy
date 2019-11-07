#!/bin/bash
VERSION=2019a

function build() {
    export GOOS=$1
    export GOARCH=$2

    if [ "$2" = "386" ]; then
        BITS=32
    else
        BITS=64
    fi
    
    if [ "$1" = "windows" ]; then
        OUTPUT="badcopy-$1$BITS-$VERSION.exe"
    else
        OUTPUT="badcopy-$1$BITS-$VERSION"
    fi 

    echo Compiling dist/$OUTPUT
    go build -ldflags "-s -w" -o dist/$OUTPUT badcopy.go
    # UPX 3.95 will crash on DARWIN
    # UPX 3.94 works
    ~/software/upx_3.94_org -9 --lzma dist/$OUTPUT
}

unset CGO_ENABLED

build linux 386
build windows 386

build darwin amd64
build linux amd64
build windows amd64

unset GOOS
unset GOARCH
