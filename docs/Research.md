# Research

The following document contains various topics that were researched during the course of the project that could be useful to future developers.

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
