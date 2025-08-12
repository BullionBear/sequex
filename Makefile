BINARY=node
PACKAGE="github.com/BullionBear/sequex"
VERSION := $(shell git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null)
COMMIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_TIMESTAMP := $(shell date '+%Y-%m-%dT%H:%M:%S')
LDFLAGS := -X '${PACKAGE}/env.Version=${VERSION}' \
           -X '${PACKAGE}/env.CommitHash=${COMMIT_HASH}' \
           -X '${PACKAGE}/env.BuildTime=${BUILD_TIMESTAMP}'

build:
	rm -rf docs
	swag init --parseDependency --parseInternal -g cmd/master/main.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o ./bin/master-linux-x86 cmd/master/main.go
	env GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o ./bin/node-linux-x86 cmd/node/main.go	

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
	@echo "Organizing generated files..."
	@if [ -f internal/model/protobuf/example/rng.pb.go ]; then \
		echo "RNG protobuf file generated in internal/model/protobuf/example/"; \
	fi

clean:
	rm -rf bin/*
	rm -rf logs/*


.PHONY: build, clean, proto