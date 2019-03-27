SOLAR System
============

<div class="subtitle"><b>S</b>imple <b>O</b>rchestration of the <b>L</b>ifecycle <b>A</b>utomation of <b>R</b>esources - a reference implementation</div>

---

SOLAR-System is a reference implementation of the concepts presented in the book [Orchestration (Lifecycle Automation)](https://github.com/BernardTsai/book/raw/master/book.pdf) . It is deliberately intended to be kept as simple as possible and in its current state is not intended to serve for production purposes.

It has been implemented in [GO](https://golang.org/) and the source code together with sample data can be found on [https://github.com/BernardTsai/solar](https://github.com/BernardTsai/solar).

This readme serves as a quickstart introduction to the SOLAR system covering which prerequisites will need to be met, the installation procedure, and how to configure and make use of the reference implementation.

Prerequisites
-------------

SOLAR basically only requires:

- access to a **BASH** environment with internet access,
- a tool e.g. **wget** to retrieve files from the internet,
- **tar** - a tool to handle archived information and  
- a **golang** runtime environment. Instructions of how to install GO onto the preferred operating system can be found on the golang website: [Getting Started](https://golang.org/doc/install).

Installation
------------

The installation of SOLAR is described here for a Linux environment with a BASH command line interface (but could be applied with only slight modifications to other operating systems) and requires following steps:

1. Create a root directory <SOLAR> and change into this directory.

    ```
    > mkdir <SOLAR>    # create the working directory    
    > cd <SOLAR>       # change into the working directory    
    ```

2. Retrieve the solar src files from SOLAR GitHub repository:

    ```
    > wget https://github.com/BernardTsai/solar/archive/master.zip
    ```

3. Unpack the archive (omitting the top-level folder):

    ```
    > tar xvf master.zip --strip 1
    ```

4. Source the setup script:

    ```
    > source setup.sh
    ```
    The script will create the proper go environment for SOLAR by creating the required directories, downloading the required dependencies and installing the SOLAR binaries.

Configuration
-------------

The configuration file of solar has the name "**solar-conf.yaml**". It has following structure:

```
MSG:
  Events: events
  Status: status
```

It currently lists the topics which the message bus interface should use in order to publish task event and status related information.

The solar binary looks for the configuration file in the current working directory and will reflect the information it finds there in the course of its initialisation.

Usage
-----

Simply invoke the **solar** binary on the command line and the application will respond with a command prompt waiting for input:

```
> solar
SOLAR Version 1.0.0
>>>
```

The available subcommands can be identified by simply entering help:

```
>>> help
Commands:
  architecture      architecture commands
  clear             clear the screen
  comment           comment
  component         component commands
  domain            domain commands
  exit              exit the program
  help              display help
  model             model commands
  output            output commands
  solution          solution commands
  task              task commands
  usage             usage command
>>>
```

The options for each subcommand can be retrieved by adding a question mark to the subcommand, e.g.:

```
>>> solution ?
usage:
  solution list <domain> <solution>
           set <domain> <filename>
           get <domain> <solution> <version>
           delete <domain> <solution> <version>
>>>
```

Solar will terminate after it has received the "exit" command and return control to BASH:

```
>>> exit
>
```

**Web-UI**

SOLAR will in addition expose a REST interface and a web-based UI which can be accessed via the URL: http://localhost/solar/index.html.
