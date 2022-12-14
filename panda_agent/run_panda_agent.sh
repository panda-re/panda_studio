#!/bin/bash

docker build --pull --rm -f "Dockerfile" -t pandare/panda_agent "."
docker run --rm -it -v ~/.panda:/root/.panda -p 50051:50051 pandare/panda_agent