#!/bin/bash

rm -rf pb/ *_pb2*
mkdir -p pb

PYTHON_OUT=pb
GO_OUT=pb

# Generate python stubs
# pip3 install grpcio-tools
python3 -m grpc_tools.protoc \
    --python_out=$PYTHON_OUT \
    --pyi_out=$PYTHON_OUT \
    --grpc_python_out=$PYTHON_OUT \
    --go_out=$GO_OUT \
    --go_opt=paths=source_relative \
    --go-grpc_out=$GO_OUT \
    --go-grpc_opt=paths=source_relative \
    -I./proto \
    ./proto/*.proto \