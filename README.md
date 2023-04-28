# PANDA Studio

The overall task of PANDA Studio is to design a service around PANDA which enables the creation of recordings from an image specification and a set of interactions. The following instructions help set up the environment for PANDA Studio to get users started. For information on deployment, development, testing, and a user guide see the [docs](./docs) folder.

## Ubuntu

PANDA Studio was designed to run on Ubuntu 20.04 and support for other operating systems is not officially supported. It is recommended to either use a [native ubuntu install](https://releases.ubuntu.com/focal/), [WSL and Windows](https://learn.microsoft.com/en-us/windows/wsl/install), or an [Ubuntu virtual machine](https://ubuntu.com/tutorials/how-to-run-ubuntu-desktop-on-a-virtual-machine-using-virtualbox#1-overview).

## Docker Instructions

PANDA Studio makes use of Docker to complete several parts of its functionality. Follow the instructions on [installing Docker](https://www.docker.com/products/docker-desktop/). For more information on Docker refer to the corresponding section in the [research memo](docs/RESEARCH.md#Docker).

## Download PANDA Studio

1. First clone this repository

```
git clone https://github.com/panda-re/panda_studio
```

2. Then change into the corresponding directory

```
cd panda_studio
```

## Deploying PANDA Studio

If you wish to deploy PANDA Studio, refer to the [deployment guide](./docs/DEPLOYMENT.md). If you wish to test and develop PANDA Studio, continue.

## PANDA Studio Backend 

The backend of PANDA Studio executes interactions specified in the frontend to create a recording for the user. It communicates with the frontend using gRPC and ProtoBuf. For more information about the specifics, see the [research document](/docs/RESEARCH.md#grpc).

### Running the Backend

1. Create directories and download files with a single command 
```
sudo make full_initial_setup
```

OR

1. Install privileged files into root directory
```
sudo make initial_setup_priviliged
```

2. Setup non-privileged aspects, including setting up the Docker Container
```
make initial_setup
```

3. Run it!

This will run an executor for the backend that will make a recording and replay that recording. The frontend uses this interface to create recordings from a web-user interface.

```
docker run -it --rm \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v /root/.panda:/root/.panda \
    -v /tmp/panda-studio:/tmp/panda-studio \
    pandare/panda_executor
```

The [executor](/cmd/panda_executor/panda_executor.go) is a good playground for trying new features or understanding the general flow of how the frontend interacts with the [backend](/panda_agent/).

Changes to the executor require rebuilding the docker container, which can be done with the command
```
make docker_executor
```
Similarly, changes to the agent require rebuilding
```
make docker_agent
```
Lastly, any changes to [ProtoBuf](/panda_agent/proto/panda_agent.proto) require running a script to update those changes
```
cd panda_agent/ && ./generate_protobuf.sh
```

4. Open the Dev Container (optional)

- This contains a mostly complete toolchain with all the tools needed for local development
- In VSCode, a toast will appear in the bottom right corner asking you to open the dev container

### Testing the Backend

The backend comes equipped with tests to ensure its functionality and they can be run using the command:
```
docker run -it --rm \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v /root/.panda:/root/.panda \
    -v /tmp/panda-studio:/tmp/panda-studio \
    pandare/panda_executor_test
```

Whenever a new feature is implemented, tests should be added to demonstrate how the feature is used and the expected inputs and outputs. Also, previous tests should still pass so that established features still function. More information can be found in the [test plan](/docs/TEST_PLAN.md) and [documentation](./cmd/panda_executor_test/README.md).

## Running the API

1. Follow the above instructions to build the required containers

2. Start the API server with the required services
```
docker compose -f ./docker-compose.dev.yml up --build panda_api
```

Note: you can use ctrl-C to quit the API. Changes made to the API will automatically be rebuilt

### Accessing the UI for backend services

- Minio's GUI is available on `localhost:9001`. The sign-in information can be found in `docker/docker-compose.dev.yml` as access key and secret key

- Mongo Express is a GUI for MongoDB and can be found at `localhost:8081`
