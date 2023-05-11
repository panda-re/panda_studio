.PHONY: all test full_initial_setup initial_setup_priviliged initial_setup docker_agent panda_agent_protoc docker_executor panda_executor

all: panda_executor panda_agent_protoc docker_agent docker_executor

test: docker_agent docker_executor_test

full_initial_setup: initial_setup_priviliged initial_setup

initial_setup_priviliged:
	mkdir -p /root/.panda
	wget -O /root/.panda/bionic-server-cloudimg-amd64-noaslr-nokaslr.qcow2 \
    	"https://www.dropbox.com/s/4avqfxqemd29i5j/bionic-server-cloudimg-amd64-noaslr-nokaslr.qcow2?dl=1"

initial_setup: docker_agent docker_executor
	mkdir -p /tmp/panda-agent

docker_agent:
	docker build -f ./docker/Dockerfile.panda-agent -t pandare/panda_agent ./panda_agent

docker_di_agent:
	docker build -f ./docker/Dockerfile.derive-image -t pandare/panda_di_agent ./derive_image

docker_executor:
	docker build -f docker/Dockerfile.panda-executor -t pandare/panda_executor .

docker_executor_test:
	docker build -f docker/Dockerfile.panda-executor-test -t pandare/panda_test_executor .

docker_api:
	docker build -f docker/Dockerfile.panda-api -t pandare/panda_api .

panda_executor: panda_agent_protoc_go
	go build -o ./bin/panda_executor ./cmd/panda_executor

panda_replay_executor: panda_agent_protoc_go
	go build -o ./bin/panda_executor ./cmd/panda_replay_executor

panda_executor_test: panda_agent_protoc_go
	go test -c ./cmd/panda_test_executor

panda_api: panda_agent_protoc_go
	go generate ./internal/api
	go build -o ./bin/panda_api ./cmd/panda_api

panda_agent_protoc: panda_agent_protoc_go panda_agent_protoc_py

panda_agent_protoc_go: panda_agent/proto/*.proto
	mkdir -p panda_agent/pb
	protoc \
		--go_out=./panda_agent/pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=./panda_agent/pb \
		--go-grpc_opt=paths=source_relative \
		-I./panda_agent/proto \
		./panda_agent/proto/*.proto
	
panda_agent_protoc_py: panda_agent/proto/*.proto
	mkdir -p panda_agent/pb
	python3 -m grpc_tools.protoc \
		--python_out=panda_agent/pb \
		--pyi_out=panda_agent/pb \
		--grpc_python_out=panda_agent/pb \
		-I./panda_agent/proto \
		./panda_agent/proto/*.proto