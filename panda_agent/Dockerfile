FROM pandare/panda

RUN mkdir /panda
WORKDIR /panda

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
        protobuf-compiler

COPY requirements.txt .

RUN pip3 install -r requirements.txt

RUN mkdir proto pb
COPY proto/* ./proto
RUN ls /panda/proto

RUN python -m grpc_tools.protoc --python_out=pb --grpc_python_out=pb -I./proto ./proto/*.proto

COPY *.py ./

EXPOSE 50051

ENTRYPOINT ["python3", "/panda/server.py"]