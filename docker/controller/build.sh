#!/bin/bash

# determine the location of the script
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# change to the root of the files
cd $DIR/../..

# ensure that the "dist" directory exists
mkdir -p dist

# cleanup dist directory
rm -rf dist/*

# Compile solar binary to the target environment (place binary in dist)
# following instructions as in:
# https://medium.com/travis-on-docker/how-to-dockerize-your-go-golang-app-542af15c27a2
# options (explained):
#   --rm                 remove instance after termination
#   -v "$PWD":"/go/"     map current directory to "/go" in the container
#   -w "/go/src"         working directory is "/go/src" in the container
#   go build -o ...      command to build solar binary and store to "/dist/solar"
docker run --rm -v "$PWD":"/go/" -w "/go/src" iron/go:dev go build -o ../dist/solar-grpc-controller tsai.eu/solar/controller/gRPC

# copy Dockerfile
cp docker/controller/Dockerfile dist

# change to "dist" directory
cd dist

# build docker image
docker build -q -t tsai/solar-grpc-controller:V1.0.0 .
