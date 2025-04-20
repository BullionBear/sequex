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

## 0.4. Install `buf` (Optional but Recommended)

`buf` is a tool for managing `.proto` files and generating code.

#### On macOS, you can install it using Homebrew:
```bash
brew install bufbuild/buf/buf
```

#### On Linux, you can install it using `curl`:
```bash
curl -sSL https://github.com/bufbuild/buf/releases/download/v1.20.0/buf-Linux-x86_64 -o /usr/local/bin/buf
chmod +x /usr/local/bin/buf
```

### 0.5. Initialize `buf` in Your Project

Run the following command in the root of your project to create a `buf.yaml` file:
```bash
buf config init
```

### 0.6. Generate Code Using `buf`

To generate Go code from your `.proto` files using `buf`, run:
```bash
buf generate
```

Make sure your `buf.gen.yaml` file is configured correctly for Go code generation.

---

## 1. Define the Service in a `.proto` File

Create a file name `greet.proto` in a directory (e.g. `./proto/greet.proto`)

## 2. Generate `.pb.go`
Generate the go code from the `.proto` file:

```bash
protoc --go_out=pkg/ --go-grpc_out=pkg/ proto/greet.proto
```

Alternatively, if you are using `buf`, you can run:
```bash
buf generate
```

## 3.1 Implement server (If needed)

Please see `cmd/greet/server.go`

## 3.2 Implement client (If needed)

Please see `cmd/greet/client/client.go`

## 3.3 Implement unittest

Please see `cmd/greet/server_test.go`