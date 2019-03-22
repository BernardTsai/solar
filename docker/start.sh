#!/bin/bash

# start container
#   options:
#     --rm             remove after termination
#     --name solar     set user-friendly name for container
#     --itd            interactive / attach TTY / daemon
#     -p 80:80         map port 80 to host port 80
docker run --rm --name solar -itd -p 80:80 tsai/solar
