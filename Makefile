.PHONY:

lint:
	golangci-lint run

cover:
	go test -race -coverprofile=cover.out -coverpkg=./... ./...
	go tool cover -html=cover.out -o cover.html

proto:
	protoc --go_out=internal/server --go_opt=paths=source_relative \
        --go-grpc_out=internal/server --go-grpc_opt=paths=source_relative \
        -I api/protobuf api/protobuf/autotest.proto