#!/bin/bash

docker build --pull --rm -f "executor/Dockerfile" -t pandare/pandaagent "executor"
docker run --rm -it -v ~/.panda:/root/.panda -p 50051:50051 pandare/pandaagent