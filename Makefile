BINARY=sequex
PACKAGE="github.com/BullionBear/sequex"
VERSION := $(shell git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null)
COMMIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_TIMESTAMP := $(shell date '+%Y-%m-%dT%H:%M:%S')
LDFLAGS := -X '${PACKAGE}/env.Version=${VERSION}' \
           -X '${PACKAGE}/env.CommitHash=${COMMIT_HASH}' \
           -X '${PACKAGE}/env.BuildTime=${BUILD_TIMESTAMP}'

build:
	env GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o ./bin/$(BINARY)-linux-x86 cmd/main.go

test:
	go test -v ./...

proto:
	@echo "Generating protobuf files..."
	@find proto -name "*.proto" -type f | while read file; do \
		rel_path=$$(echo "$$file" | sed 's|^proto/||'); \
		dir=$$(dirname "internal/model/protobuf/$$rel_path"); \
		mkdir -p "$$dir"; \
		protoc --proto_path=proto --go_out=internal/model/protobuf --go_opt=paths=source_relative "$$file"; \
		echo "Generated: $$file -> $$dir/*.pb.go"; \
	done

clean:
	rm -rf bin/*
	rm -rf logs/*


.PHONY: build, clean, proto