# PANDA Studio

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
1. Create a cache to store the default PANDA images:
```
sudo mkdir -p /root/.panda
```
1.5. Download the default x86_64 image
```
sudo wget -O /root/.panda/bionic-server-cloudimg-amd64-noaslr-nokaslr.qcow2 \
    "https://www.dropbox.com/s/4avqfxqemd29i5j/bionic-server-cloudimg-amd64-noaslr-nokaslr.qcow2?dl=1"
```

2. Create a directory for the executor to communicate with the agent:
```
mkdir -p /tmp/panda-agent
```

3. Build the PANDA agent image
```
docker build -f ./docker/Dockerfile.panda-agent -t pandare/panda_agent ./panda_agent
```

4. Build the PANDA executor image
```
docker build -f docker/Dockerfile.panda-executor -t pandare/panda_executor .
```

6. Run it!
```
docker run -it --rm \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v /root/.panda:/root/.panda \
    -v /tmp/panda-studio:/tmp/panda-studio \
    pandare/panda_executor
```

7. Open the Dev Container (optional)
    - This contains a mostly complete toolchain with all the tools needed for local development
    - In VSCode, a toast will appear in the bottom right corner asking you to open the dev container
    - The dev container is not yet fully supported (working on it!)

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
