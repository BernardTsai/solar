Docker
======

In order to build the corresponding docker images invoke the commands:

```
> docker/solar/build.sh
> docker/controller/build.sh
```

A docker compose file is available to start SOLAR and a couple of corresponding services.

```
> docker/start.sh
```

This will setup and install five containers:

* **portainer**: a container management tool  (https://www.portainer.io/)
* **kafka**: the message broker (https://github.com/wurstmeister/kafka-docker)
* **zookeeper**: a backend for kafka (https://github.com/wurstmeister/zookeeper-docker)
* **solar**: the orchestrator
* **gRPC**: the reference gRPC controller

Portainer is accessible via port 9000 of the host and solar on port 80.

To stop the containers invoke:

```
> docker/stop.sh
```
