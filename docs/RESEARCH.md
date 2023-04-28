# Research

The following document contains various topics that were researched during the course of the project that could be useful to future developers.

# Table of Contents

[PANDA](#panda)

[Docker](#docker)

[Go](#golang)

[gRPC](#grpc)

[ProtoBuf](#protobuf)

[QEMU](#qemu)

# PANDA

(Platform for Architecture Neutral Dynamic Analysis) [PANDA](https://panda.re) is a useful tool for software reverse engineering where researchers can record and replay the execution of instructions on various architectures and operating systems using QEMU.

PANDA can be controlled through the command line, a Python package called PyPANDA, or a Jupyter notebook. The objective of this project was to create a web interface for creating and storing PANDA recordings.

## Getting Started

PANDA has a [Slack channel](https://panda.re/invite.php), which can be a good resource for assistance. A great way to get started with PANDA is by following the [class in the pandare repo](https://github.com/panda-re/panda_class). The class guides users through the basic launching of PANDA, creating a recording, and analyzing recordings using plugins. There is also a [manual for PANDA](https://github.com/panda-re/panda/blob/dev/panda/docs/manual.md) and [documentation for PyPANDA](https://docs.panda.re/).

## PyPANDA

PyPANDA was highly utilized in the backend to interact with PANDA to produce the desired results to the user. It is a python interface with PANDA that allows programmers to quickly develop plugins to analyze behavior of a running system, record a system and analyze replays, or do nearly anything you can do using PANDA's C/C++ APIs.

A [YouTube video by user AdaLogics](https://www.youtube.com/watch?v=2HQqZZNC4P8) provides a good overview of a [paper published by MIT LL](https://www.ndss-symposium.org/wp-content/uploads/bar2021_23001_paper.pdf) that goes over the features and an example of PyPANDA.

# Docker

[Docker](https://www.docker.com/) is a key service used in PANDA Studio that allows for code and dependencies to be packaged up in containers. For example, PANDA has a Docker container which is used in this system.

## Docker Image

Docker images are a template with instructions for creating a Docker container. They use a Dockerfile to specify how the image is built. Docker’s website contains a useful [guide](https://docs.docker.com/build/building/packaging/) to creating an image using a Dockerfile. The Dockerfiles for this project can be found in the [docker](../docker/) directory.

## Docker Container

A Docker container is a runnable instance of a Docker image. Each container is relatively isolated from its host but there are ways to communicate with the host and other containers. More information can be found on [Docker's website](https://www.docker.com/resources/what-container/) This project uses gRPC and ProtoBuf, which have their respective subsections in this document.

Multiple containers can be run at once using the same image or different images. The theoretical limit is the system’s resources. This project uses Docker Compose, which is a Docker tool that helps define and share multi-container applications through the use of a `YAML` file.

A Docker container can execute docker commands on the host by adding the argument `-v /var/run/docker.sock:/var/run/docker.sock` when running the container. This is used by the executors because they launch the Docker container that houses the agent that controls PANDA.

## PANDA and Docker

The following steps will help a user learn about PANDA and Docker. There is also a [PANDA class](https://github.com/panda-re/panda_class) that is highly recommended to get users familiar with PANDA and Docker.

1. Install PANDA
To install PANDA using docker, run the following command:

```
docker pull pandare/panda
```
It is also recommended to create a directory for storing qcows to give to the docker container.

2. Run Container
To run the PANDA docker container, use the following command:
```
docker run -it -v ~/<dir>:/<dir> pandare/panda bash
```
Where `dir` is the directory of the images. From there, a docker shell will start.

3. Run PANDA
Next is running a PANDA system which can be done using the command below:
```
panda-system-i386 -nographic /<dir>/<qcow>
```
Where `qcow` is the image to be loaded into PANDA. The `i386` can be changed into any architecture supported by PANDA like `x86_64`, for example.

Note that this will not boot up a GUI so this will only work with some images. To run it with a display, do the following:
```
docker run -it -p 5900:5900 -v ~/<dir>:/<dir> pandare/panda bash
```
The `-p` option will export the default VNC port outside of the container for your PC to access.
Next is running a PANDA system which can be done using the command below:
```
panda-system-i386 -vnc 0.0.0.0:0 /<dir>/<qcow>
```
Once that has started, connect to the display by opening VNC viewer and connecting to `127.0.0.1:5900`.

# Golang

[Go](https://go.dev/) is an open-source programming language supported by Google that is useful for cloud and networking services, command-line interfaces, web development, and development operations and site reliability engineering.

In this project, Go is used for the API and the gRPC and ProtoBuf framework for the frontend to communicate with the backend.

## Backend Testing

Testing and validation are essential to developing a product. Tests ensure the functionality of the system and can be used in tandem with continuous integration to prevent new features from breaking the functionality of the system. Also, tests demonstrate the expected input and output of the system so other developers can get a better understanding of the code.

Go has a [package for testing](https://pkg.go.dev/testing) that was used for the backend. The package can run tests, benchmarks, examples, and fuzzing. Files for testing should end in `_test.go` and will automatically execute functions of the form `TestXxx` and can allow for subtests that share the same setup and teardown code, which will be extremely useful since the system has such code that does take time to execute. For more information on the backend testing, see its [readme](../cmd/panda_executor_test/README.md) and the [test plan](./Test_Plan.md)


# gRPC

[gRPC](https://grpc.io/) is a modern open-source high performance Remote Procedure Call (RPC) framework that can run in any environment. gRPC has many similarities to Protocol Buffers and uses Protocol Buffers to structure data. For more information on Protocol Buffers (ProtoBuf) refer to its section below. The main advantage of gRPC is language-neutral so a client in Go can communicate with a client in Python. The server sends a request for the client to execute a method and the client will execute it and provide a response. 

This project uses gRPC and ProtoBuf to communicated between the frontend and backend. As a result, both ends use a language best suited for its implementation rather than making sacrifices for the sake of convenient communication. The frontend implementation can be seen in [grpc_agent.go](../internal/panda_controller/grpc_agent.go) and the backend implementation can be seen in [server.py](../panda_agent/server.py).

Also, the shorter communication of gRPC allows for more flexible error handling rather than a monolithic approach that would require restarting the entire program if an exception occurs.

# ProtoBuf

[Protocol Buffers](https://protobuf.dev/overview/) provide a language-neutral, platform-neutral, extensible mechanism for serializing structured data in a forward-compatible and backward-compatible way. The protocol would allow for easier data communication between different languages. This is especially important for communication between the frontend and backend as they are in different languages, so having a universal way to specify data is greatly beneficial. The data can be sent in a file for the other to read or sent over a network for the other to act upon. With gRPC it can define communication protocols.

The ProtoBuf can be seen in [panda_agent.proto](../panda_agent/proto/panda_agent.proto).

# QEMU

PANDA utilizes QEMU to emulate a whole system so that it can record all code executing and all data. [QEMU](https://www.qemu.org/) stands for Quick Emulator, and it is an open-source generic emulator and virtualizer. As an emulator, it runs a virtual model of a system (CPU, memory, etc.) to run any operating system. As a virtualizer, QEMU can run kernel virtual machines. Lastly, QEMU can be used in user mode emulation, which allows processes made for one type of CPU to be run on another type.

## qcow

QEMU, and by extension PANDA, use disk images for emulating machines. A disk image is called a qcow, which stands for "QEMU Copy on Write" and uses a disk storage optimization strategy that delays allocation of storage until it is needed. Files in qcow format can contain a variety of disk images which are generally associated with specific guest operating systems.
