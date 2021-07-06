.PHONY: lint
lint:
	golangci-lint run

.PHONY: proto
proto:
	protoc --go_out=internal/server --go_opt=paths=source_relative \
        --go-grpc_out=internal/server --go-grpc_opt=paths=source_relative \
        -I api/protobuf api/protobuf/autotest.proto