#!/bin/bash

docker build --pull --rm -f "panda_agent/Dockerfile" -t pandare/panda_agent "panda_agent"
docker run --rm -it -v ~/.panda:/root/.panda -p 50051:50051 pandare/panda_agent