FROM mcr.microsoft.com/vscode/devcontainers/go:0-1.19-bullseye

# Install required packages
RUN apt-get update && apt-get -y install protobuf-compiler

# Install required go tools
USER vscode
ARG go_proto_version="1.28"
ARG go_grpc_version="1.2"
RUN go install "google.golang.org/protobuf/cmd/protoc-gen-go@v${go_proto_version}" \
    && go install "google.golang.org/grpc/cmd/protoc-gen-go-grpc@v${go_grpc_version}"
