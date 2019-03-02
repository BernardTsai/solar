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
go get github.com/pkg/errors
go get gopkg.in/yaml.v2
go get gopkg.in/abiosoft/ishell.v2
go get github.com/spf13/viper
go get github.com/spf13/viper
go get github.com/gorilla/mux

# install binaries
cd $SRCDIR
go install tsai.eu/solar/cmd/solar

# change to root directory
cd $ROOTDIR
