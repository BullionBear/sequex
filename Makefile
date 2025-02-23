BINARY := alex

PACKAGE="github.com/BullionBear/crypto-trade"
VERSION := $(shell git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null)
COMMIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_TIMESTAMP := $(shell date '+%Y-%m-%dT%H:%M:%S')
LDFLAGS := -X '${PACKAGE}/env.Version=${VERSION}' \
           -X '${PACKAGE}/env.CommitHash=${COMMIT_HASH}' \
           -X '${PACKAGE}/env.BuildTime=${BUILD_TIMESTAMP}'

gen:
	protoc --go_out=. --go-grpc_out=. api/proto/feed.proto

build:
	go build -ldflags="$(LDFLAGS)" -o ./bin/$(BINARY) cmd/main.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o ./bin/$(BINARY)-linux-x86 cmd/$(BINARY)/*.go


clean:
	rm -rf bin/*
	rm -rf logs/*


.PHONY: gen, build, clean