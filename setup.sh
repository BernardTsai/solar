#!/bin/bash
# Instructions - source this file to setup the proper environment:
# > source setup.sh

# determine working directory
ROOTDIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
BINDIR=$ROOTDIR/bin
SRCDIR=$ROOTDIR/src

# create required directories
mkdir -p $ROOTDIR/src
mkdir -p $ROOTDIR/bin
mkdir -p $ROOTDIR/pkg

# export GOPATH
export GOPATH=$ROOTDIR

# update PATH
if [ -d "$BINDIR" ] && [[ ":$PATH:" != *":$BINDIR:"* ]]; then
    PATH="${PATH:+"$PATH:"}$BINDIR"
fi

# import required libraries
go get github.com/google/uuid
go get gopkg.in/yaml.v2
go get gopkg.in/abiosoft/ishell.v2
go get github.com/spf13/viper
go get github.com/gorilla/mux
go get github.com/segmentio/kafka-go
go get google.golang.org/grpc
go get github.com/golang/protobuf/protoc-gen-go
go get github.com/rs/zerolog/log
go get bou.ke/monkey
go get github.com/cbroglie/mustache

# test
go test -cover                   \
  tsai.eu/solar/util             \
  tsai.eu/solar/model            \
  tsai.eu/solar/msg              \
  tsai.eu/solar/controller/dummy \
  tsai.eu/solar/controller/gRPC  \
  tsai.eu/solar/controller       \
  tsai.eu/solar/engine           \
  tsai.eu/solar/monitor          \
  tsai.eu/solar/cli              \
  tsai.eu/solar/api

# install binaries
# cd $SRCDIR
go install tsai.eu/solar/cmd/solar
go install tsai.eu/solar/cmd/controller

# change to root directory
cd $ROOTDIR
