PACKAGE="github.com/BullionBear/sequex"
VERSION := $(shell git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null)
COMMIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_TIMESTAMP := $(shell date '+%Y-%m-%dT%H:%M:%S')
LDFLAGS := -X '${PACKAGE}/env.Version=${VERSION}' \
           -X '${PACKAGE}/env.CommitHash=${COMMIT_HASH}' \
           -X '${PACKAGE}/env.BuildTime=${BUILD_TIMESTAMP}'

PROTO_DIR = protobuf
GO_OUT_DIR = internal/model
PROTOC = protoc
PROTOC_GEN_GO = protoc-gen-go
# Find all proto files
PROTO_FILES := $(shell find $(PROTO_DIR) -name "*.proto")

build:
	rm -rf docs
	swag init --parseDependency --parseInternal -g cmd/master/main.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o ./bin/master-linux-x86 cmd/master/main.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o ./bin/node-linux-x86 cmd/node/main.go	

test:
	go test -v ./...

proto:
	@echo "Generating Go code from protobuf files..."
	@mkdir -p $(GO_OUT_DIR)
	@for proto_file in $(PROTO_FILES); do \
		echo "Processing $$proto_file..."; \
		$(PROTOC) \
			--proto_path=. \
			--go_out=$(GO_OUT_DIR) \
			--go_opt=paths=source_relative \
			--go-grpc_out=$(GO_OUT_DIR) \
			--go-grpc_opt=paths=source_relative \
			$$proto_file; \
	done
	@echo "Protobuf generation completed!"

clean:
	rm -rf bin/*
	rm -rf logs/*
	rm -rf $(GO_OUT_DIR)/protobuf


.PHONY: build, clean, proto