#!/bin/bash

# Generate golang stubs
# protoc --go-out=. --go_opt=paths=source_relative proto/panda_interface.proto
python3 -m grpc_tools.protoc --python_out=proto --pyi_out=proto --grpc_python_out=proto -I./proto ./proto/panda_interface.proto

# Generate python stubs
# pip3 install grpcio-tools
python3 -m grpc_tools.protoc \
    --python_out=proto \
    --pyi_out=proto \
    --grpc_python_out=proto \
    --go_out=proto \
    --go_opt=paths=source_relative \
    -I./proto \
    ./proto/panda_interface.proto \