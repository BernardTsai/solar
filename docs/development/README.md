Development
===========

This document gives a short introduction of how to setup a development environment for SOLAR.

Prerequisites
-------------

* BASH command line interface
* internet access
* git (https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
* golang (https://golang.org/doc/install)
* optional: docker and docker-compose (https://docs.docker.com/install/ and https://docs.docker.com/compose/install/)

Installation
------------

1. Clone the SOLAR repository

```
> git clone https://github.com/BernardTsai/solar.git
```

2. Change to the "solar" directory

```
> cd solar
```

3. Setup the environment and source dependencies

```
> source setup.sh
```

The setup script will install the required dependencies and update the GOPATH.
It will create the following two binaries as well:

* **solar**: the orchestrator
* **gRPC**: a reference gRPC controller

4. Configuration

The configuration file of SOLAR is written in yaml and is named: "solar-conf.yaml" and currently needs to be located in the current working directory (this will fixed at some point in time).

```
MSG:
  Notifications:  notifications
  Monitoring:     monitoring
  Address:        127.0.0.1:9092
CORE:
  IDENTIFIER: solar
  LOGLEVEL:   debug
CTRL:
  dummy: 127.0.0.1:10000
```

In the MSG section it defines where to find the Kafka message broker (if solar can't find the message broker it will stop attempting to send and receive messages) and which topics to use for receiving monitoring information and publishing notifications.

The CORE section defines an identifier for the SOLAR node and the log level to use.

The CTRL section lists all gRPC controllers which need to be used with a help of a simple map. This map states the type of the controller as the key and the service endpoint of the controller as a "hostname:port" string. If a controller can not be found it is ignored and replaced with the internal default controller.

5. Starting

The gRPC controller should be started first and then solar. In the case solar can not connect to the gRPC controller it will forward the corresponding requests to the internal default controller.

The SOLAR orchestrator offers a simple command line interface and web interface.
Typing "help" into the command line lists all available commands.
The web interface is accessible at port 80 of localhost.

6. Stoping

The SOLAR orchestrator can be terminated by either entering "exit" into the command line interface or pressing "ctrl-c".

The gRPC controller can be stopped by pressing "ctrl-c".
