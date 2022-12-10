FROM pandare/panda

RUN mkdir /panda
WORKDIR /panda

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
        protobuf-compiler

COPY python/requirements.txt .

RUN pip3 install -r requirements.txt

COPY proto/panda_interface.proto ./

RUN python -m grpc_tools.protoc --python_out=. --grpc_python_out=. -I. ./*.proto

COPY python/panda_slave.py ./

EXPOSE 50051

ENTRYPOINT ["python3", "/panda/panda_slave.py"]