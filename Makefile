BINARY=sequex
PACKAGE="github.com/BullionBear/sequex"
VERSION := $(shell git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null)
COMMIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_TIMESTAMP := $(shell date '+%Y-%m-%dT%H:%M:%S')
LDFLAGS := -X '${PACKAGE}/env.Version=${VERSION}' \
           -X '${PACKAGE}/env.CommitHash=${COMMIT_HASH}' \
           -X '${PACKAGE}/env.BuildTime=${BUILD_TIMESTAMP}'

codegen:
	protoc --go_out=pkg/ --go-grpc_out=pkg/ proto/greet.proto
	protoc --go_out=pkg/ --go-grpc_out=pkg/ proto/sequex.proto
	protoc --go_out=pkg/ --go-grpc_out=pkg/ proto/solvexity.proto

build:
	env GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o ./bin/$(BINARY)-linux-x86 cmd/sequex/server.go

test:
	go test -v ./...

clean:
	rm -rf bin/*
	rm -rf logs/*


.PHONY: codegen, build, clean