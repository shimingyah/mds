#!/bin/bash

# This script generates Go representations of Protobuf protocols.
# It will generate Go code in the local dir use gogo protocols.

set -ex

function init() {
    PROGRAM=$(basename "$0")
    if [ -z $(go env GOPATH) ]; then
        printf "Error: the environment variable GOPATH is not set, please set it before running %s\n" $PROGRAM > /dev/stderr
        exit 1
    fi

    GOPATH=$(go env GOPATH)
    export PATH=$GOPATH/bin:$PATH

    brew install protobuf

    # google protobuf(slow)
    go get github.com/golang/protobuf/proto
    go get github.com/golang/protobuf/protoc-gen-go

    # gogo protobuf(fast)
    go get github.com/gogo/protobuf/proto
    go get github.com/gogo/protobuf/protoc-gen-gofast
}

function gen() {
    # protoc --go_out=plugins=grpc:. mds.proto
    protoc --gofast_out=plugins=grpc:. mds.proto
}

gen