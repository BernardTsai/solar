#!/bin/bash

# determine the location of the script
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# change to the root of the files
cd $DIR

# start containers
docker-compose -p solar down
