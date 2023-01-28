
panda_executor: panda_agent_protoc_go
	go build -o ./bin/panda_executor ./cmd/panda_executor

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