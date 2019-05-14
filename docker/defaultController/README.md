Instructions
============

This directory contains the scripts and docker definitions to build and run a dummy REST controller.

The assumption is that the sources have already been cloned from https://github.com/BernardTsai/solar.

Building the docker image
-------------------------

Invoke the script "build.sh" and verify that an image with the tag "tsai/solar-default-controller" has been created .

```
> build.sh
> docker images | grep tsai/solar-default-controller
tsai/solar-default-controller              V1.0.0              aecee5613cb2        4 seconds ago      12.2MB
>
```

Starting default controller as a daemon
---------------------------------------

Invoke the script "start.sh"

```
> start.sh
```

Stopping the default controller
-------------------------------

Invoke the script "stop.sh". After executing the script the container will be removed automatically.

```
> stop.sh
```
