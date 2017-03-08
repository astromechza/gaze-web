#!/usr/bin/env bash

set -e

# change these!
GO_PROJECT=github.com/AstromechZA/gaze-web
SHORT_NAME=gaze-web
VERSION=$(cat VERSION)


function buildbinary {
    goos=$1
    goarch=$2

    echo "Building official $goos $goarch binary"

    LONG_NAME="${SHORT_NAME}-${VERSION}_${goos}_${goarch}"
    outputfolder="build/$LONG_NAME/$SHORT_NAME"
    echo "Output Folder $outputfolder"
    mkdir -pv $outputfolder

    export GOOS=$goos
    export GOARCH=$goarch

    govvv build -i -v -o "$outputfolder/$SHORT_NAME" $GO_PROJECT

    cp -r static/ $outputfolder/static
    cp -r templates/ $outputfolder/templates

    tar -czvf "build/$LONG_NAME.tar.gz" -C "build/$LONG_NAME" "$SHORT_NAME"
    echo
}

# clear build dir
rm -rfv build/*

# build local
unset GOOS
unset GOARCH
govvv build $GO_PROJECT

# build for mac
buildbinary darwin amd64

# build for linux
buildbinary linux amd64
