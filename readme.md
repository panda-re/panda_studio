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

5. Run it!
```
docker run -it --rm \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v /root/.panda:/root/.panda \
    -v /tmp/panda-studio:/tmp/panda-studio \
    pandare/panda_executor
```

6. Open the Dev Container (optional)
    - This contains a mostly complete toolchain with all the tools needed for local development
    - In VSCode, a toast will appear in the bottom right corner asking you to open the dev container
    - The dev container is not yet fully supported (working on it!)

## Other Instructions
Prerequisites:
- For building/running the web ui, you will need node 14.15.1.
- You may also need node 14.13.1 to install dependencies. If 14.15.1 doesn't work, switch to 14.13.1 to install dependencies then
switch back to 14.15.1 to build.
- use nvm (https://github.com/coreybutler/nvm-windows/releases - nvm-setup.zip) to manage multiple versions of node
- With node 14.13.1, run "npm install --force" to install necessary dependencies
- Python 3.10 for the server project.
- Python 3.7 for the backend project.
- You should create a separate venv in each of these projects to install respective dependencies. A tutorial on how to do
so can be found here: https://docs.python.org/3/library/venv.html
- Once inside the active venv, you will need to pip install flask and flask_cors for the server project
- You will also need to pip install pandare inside the backend project

Running the web ui:
- Install the necessary dependencies (described above)
- Ensure that you are in the panda-studio-ui directory
- Run "npm run dev" to start the web client

Running the web server:
- Make sure necessary dependencies are installed
- Run "flask --app server.py run"
- The web server needs to be running in order for the UI to properly call to the API

Running the backend process:
- Don't. It just exists. Something will call this code to run it, but the dependencies do need to exist.



