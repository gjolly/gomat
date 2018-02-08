#!/bin/bash

buildDir="build"
targetDir="Daemon/main"
targetName="gomatDaemon"

function clean {
    # rm -R ../$buildDir
    rm /tmp/gomat.sock
}

function build {
    cd ..
    if [ ! -d $buildDir ]; then
        mkdir $buildDir
    fi
    cd $buildDir

    go build -o $targetName ../$targetDir
}

function run {
    ../build/$targetName
}

build
run
clean