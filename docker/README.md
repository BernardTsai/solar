Quick Start
===========

The docker-compose instructions will start 5 containers:

1) **portainer:**  management
2) **zookeeper:**  data store
3) **kafka:**      messaging
4) **solar:**      lifecycle automation
5) **controller:** example controller

Prerequisites
-------------

Make sure following prerequisites are met:

- bash shell with internet access
- docker & docker-compose
- golang 

Installation
------------

First build solar and the controller:

```
> docker/solar/build.sh
> docker/controller/build.sh
```

after that start the containers:

```
docker/start.sh
```

Usage:
------

The container management can be accessed by opening the URL: "http://localhost:9000" and making use of the default credentials: "admin/admin".

The SOLAR application is accessible via the URL: "http://localhost".

Deinstallation
--------------

Stop the containers with the command:

```
docker/stop.sh
```
