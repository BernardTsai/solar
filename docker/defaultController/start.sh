#!/bin/bash

# start container
#   options:
#     --restart on-failure                                            automatically restart
#     --name solar-default-controller                                 set user-friendly name for container
#     --itd                                                           interactive / attach TTY / daemon
#     -p 10000:10000                                                  map port 10000 to host port 10000
#     -l tsai.eu.solar.controller.image=tsai/solar-default-controller tag the image name
#     -l tsai.eu.solar.controller.version=V1.0.0                      tag the image version

docker run                                                    \
  --restart on-failure                                        \
  --name solar-default-controller                             \
  -itd                                                        \
  -p 10000:10000                                              \
  -l tsai.eu.solar.controller.image=tsai/solar-k8s-controller \
  -l tsai.eu.solar.controller.version=V1.0.0                  \
  tsai/solar-default-controller:V1.0.0
