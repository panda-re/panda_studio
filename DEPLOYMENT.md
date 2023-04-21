# Deployment Instructions

## 1. Install Docker

Instructions to install Docker may be different depending on your system.
Please refer to [Docker's instructions](https://docs.docker.com/engine/install/) for your specific system.

## 2. Download PANDA Studio

First clone this repository to the system you wish to deploy to

`git clone https://github.com/panda-re/panda_studio`

Then change your directory to the repository

`cd panda_studio`

## 3. Set up PANDA Studio

The PANDA Agent container must be built separately with the command

`make docker_agent`

Additionally, a temporary folder must be created to allow PANDA Studio to communicate with the ephemeral agent containers
`mkdir /tmp/panda-studio`

## 4. Configuration

Copy the file `.env.test` to `.env` and change the environment variables to your liking.

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

For more advanced deployment, see the [Docker Compose documentation](https://docs.docker.com/compose/). 