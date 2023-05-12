# Deployment Instructions

This guide is recommended if you wish to deploy PANDA Studio on a server or other system, which is useful for granting multiple people access to the installation.

If you wish to test and develop PANDA Studio, refer to the getting started [instructions](../README.md).

## 1. Install Docker

Instructions to install Docker may be different depending on your system.
Please refer to [Docker's instructions](https://docs.docker.com/engine/install/) for your specific system.

## 2. Download PANDA Studio

First clone this repository to the system you wish to deploy to

`git clone https://github.com/panda-re/panda_studio`

Then change your directory to the repository

`cd panda_studio`

## 3. Set up PANDA Studio

The PANDA Agent container must be built separately with the command. Note that you may need to use `sudo` if you get permission errors.

`make docker_agent`

Additionally, a temporary folder must be created to allow PANDA Studio to communicate with the ephemeral agent containers

`mkdir /tmp/panda-studio`

## 4. Configuration

Copy the file `.env.example` to `.env` and change the environment variables to your liking.

`cp .env.example .env`

### List of Environment Variables

| Variable           | Description                                                    |
| ------------------ | -------------------------------------------------------------- |
| `S3_ACCESS_KEY`    | The root username for the MinIO                                |
| `S3_SECRET_KEY`    | The root password for the MinIO                                |
| `MONGODB_USERNAME` | The username to set for MongoDB                                |
| `MONGODB_PASSWORD` | The password to set for MongoDB                                |
| `MONGODB_DATABASE` | The name of the MongoDB database to use with PANDA Studio      |

## 5. Start Docker Compose

PANDA Studio and all its required components can be started with the command

`docker compose up -d`

Following this, Docker will pull and build all the required images and start the containers. The web UI should be available at port 8080 on the host machine.


For more advanced deployment, see the [Docker Compose documentation](https://docs.docker.com/compose/).