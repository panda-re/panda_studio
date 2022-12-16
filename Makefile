
panda_executor: panda_agent_protoc
	go build -o ./bin/panda_executor ./cmd/panda_executor

panda_agent_protoc: panda_agent/proto/*.proto
	mkdir -p panda_agent/pb
	protoc \
		--go_out=./panda_agent/pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=./panda_agent/pb \
		--go-grpc_opt=paths=source_relative \
		-I./panda_agent/proto \
		./panda_agent/proto/*.proto