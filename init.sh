#!/bin/bash
# Instructions - source this file to setup the proper environment:
# > source init.sh

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

# change to root directory
cd $ROOTDIR
