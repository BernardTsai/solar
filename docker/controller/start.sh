#!/bin/bash

# start container
#   options:
#     --restart on-failure                             automatically restart
#     --name solar                                     set user-friendly name for container
#     --itd                                            interactive / attach TTY / daemon
#     -p 80:80                                         map port 80 to host port 80
#     -v /var/run/docker.sock:/var/run/docker.sock     expose the docker API socket to the container
docker run                                     \
  --restart on-failure                         \
  --name solar-grpc-controller                 \
  -itd                                         \
  -p 10000:10000                               \
  tsai/solar-grpc-controller:V1.0.0
