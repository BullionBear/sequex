# gRPC Service

## 0. Install gRPC and Make Sure `protoc` is Working Correctly

### 0.1. Install Protocol Buffers Compiler (protoc): You need the protoc compiler to generate Go code from your .proto files.

#### On macOS, you can install it using Homebrew:
```bash
brew install protobuf
```

#### On Ubuntu, you can install it using `apt`:
```bash
sudo apt-get install protobuf-compiler
```

### 0.2. Install the Go plugin for `protoc`:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 0.3. Verify the installation

```bash
protoc --version
```

## 1. Define the Service in a `.proto` File

Create a file name `greet.proto` in a directory (e.g. `./proto/greet.proto`)

## 2. Generate `.pb.go`
Generate the go code from the `.proto` file:

```bash
protoc --go_out=pkg/ --go-grpc_out=pkg/ proto/greet.proto
```

## 3.1 Implement server (If needed)

Please see `cmd/greet/server.go`

## 3.2 Implement client (If needed)

Please see `cmd/greet/client/client.go`

## 3.3 Implement unittest

Please see `cmd/greet/server_test.go`