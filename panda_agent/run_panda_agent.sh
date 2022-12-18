#!/bin/bash

mkdir -p shared

docker build --pull --rm -f "Dockerfile" -t pandare/panda_agent "."
docker run --rm -it \
  -v ~/.panda:/root/.panda \
  -v $(pwd)/shared:/panda/shared \
  -p 50051:50051 \
  pandare/panda_agent