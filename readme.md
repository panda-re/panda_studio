# PANDA Studio

The overall task of PANDA Studio is to design a service around PANDA which enables the creation of recordings from an image specification and a set of interaction.

## Docker Instructions

First, make sure Docker is properly installed on your system.

Note that some of these instructions are merely workarounds for the time being.

### Installing PANDA Studio 
1. Install privlidged files into root directory
	makes a .panda directory in root user, and downloads the main server QCOW to it.
```
sudo make initial_setup_priviliged
```

2. Setup non-privlidged aspects, including setting up the Docker Container
```
make initial_setup
```

## OR

1. Install all files with single command 
```
sudo make full_initial_setup
```

### Preparation
1. Build the PANDA agent image
```
docker build -f ./docker/Dockerfile.panda-agent -t pandare/panda_agent ./panda_agent
```

2. Build the PANDA executor image
```
docker build -f docker/Dockerfile.panda-executor -t pandare/panda_executor .
```

3. Run it!
```
docker run -it --rm \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v /root/.panda:/root/.panda \
    -v /tmp/panda-studio:/tmp/panda-studio \
    pandare/panda_executor
```

4. Open the Dev Container (optional)
    - This contains a mostly complete toolchain with all the tools needed for local development
    - In VSCode, a toast will appear in the bottom right corner asking you to open the dev container
    - The dev container is not yet fully supported (working on it!)

## Running Tests
PANDA Studio is equipped with a test to verify the functionality of the agent that controls PANDA. More information can be found in [`cmd/panda_test_executor/`](https://github.com/panda-re/panda_studio/tree/main/cmd/panda_test_executor)

Run the tests with
```
docker run -it --rm \
     -v /var/run/docker.sock:/var/run/docker.sock \
     -v /root/.panda:/root/.panda \
     -v /tmp/panda-studio:/tmp/panda-studio \
     pandare/panda_test_executor
```

## Running the API

1. Follow the above instructions to build the required containers

2. Start the API server with the required services
```
docker compose -f ./docker/docker-compose.dev.yml up --build panda_api
```

Note: you can use ctrl-C to quit the API. Changes made to the API will automatically be rebuilt

### Accessing the UI for backend services

- Minio's GUI is available on `localhost:9001`. The sign-in information can be found in `docker/docker-compose.dev.yml` as access key and secret key

- Mongo Express is a GUI for MongoDB and can be found at `localhost:8081`
