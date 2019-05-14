Instructions
============

This directory contains the scripts and docker definitions to build and run a k8s REST controller.

The assumption is that the sources have already been cloned from https://github.com/BernardTsai/solar.

Building the docker image
-------------------------

Invoke the script "build.sh" and verify that an image with the tag "tsai/solar-k8s-controller" has been created .

```
> build.sh
> docker images | grep tsai/solar-k8s-controller
tsai/solar-k8s-controller                  V1.0.0              fe0322749ce0        10 seconds ago      14.2MB
>
```

 Starting k8s controller as a daemon
----------------------------------------

Invoke the script "start.sh"

```
> start.sh
```

Stopping the k8s controller
---------------------------

Invoke the script "stop.sh". After executing the script the container will be removed automatically.

```
> stop.sh
```
