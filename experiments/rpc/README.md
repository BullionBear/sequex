# NATS RPC Client/Server Example

This example demonstrates a simple RPC (Remote Procedure Call) system using NATS messaging. The system consists of a server that handles RPC requests and a client that makes RPC calls.

## Features

- **Request-Response Pattern**: Uses NATS request-reply messaging
- **JSON Protocol**: All messages are serialized as JSON
- **Error Handling**: Proper error handling for invalid requests and server errors
- **Multiple Operations**: Supports addition and string operations
- **Unique Request IDs**: Each request gets a unique identifier

## Architecture

```
Client                    NATS Server                    RPC Server
   |                         |                              |
   |-- Request "add" ------->|                              |
   |                         |-- Request "add" ----------->|
   |                         |                              |-- Process
   |                         |                              |-- Response
   |                         |<-- Response ----------------|
   |<-- Response ------------|                              |
```

## Available RPC Methods

### 1. Addition Operation (`add`)
- **Request**: `{"a": 10, "b": 20}`
- **Response**: `{"sum": 30}`

### 2. String Operation (`string`)
- **Request**: `{"text": "Hello"}`
- **Response**: `{"length": 5, "upper": "Hello (length: 5)"}`

## Running the Example

### Prerequisites

1. Install NATS server:
   ```bash
   go install github.com/nats-io/nats-server/v2@latest
   ```

2. Install NATS Go client:
   ```bash
   go get github.com/nats-io/nats.go
   ```

### Step 1: Start NATS Server

```bash
nats-server
```

### Step 2: Start the RPC Server

In one terminal:
```bash
cd experiments/rpc/srv
go run main.go
```

You should see:
```
Connected to NATS server
RPC Server started. Listening for requests on 'rpc.requests'
Available methods: add, string
```

### Step 3: Run the RPC Client

In another terminal:
```bash
cd experiments/rpc/client
go run main.go
```

You should see output like:
```
Connected to NATS server
Making addition RPC call...
Addition result: 10 + 20 = 30
Making string RPC call...
String result: length=16, text='Hello, NATS RPC! (length: 16)'
Making multiple RPC calls...
Call 1 result: 10 + 5 = 15
Call 2 result: 20 + 10 = 30
Call 3 result: 30 + 15 = 45
Testing error handling with unknown method...
Expected error received: Unknown method: unknown_method
RPC client demo completed
```

## Code Structure

### Server (`srv/main.go`)
- `RPCHandler`: Handles incoming RPC requests
- `handleAdd()`: Processes addition requests
- `handleString()`: Processes string requests
- `handleRequest()`: Main request dispatcher

### Client (`client/main.go`)
- `RPCClient`: Makes RPC calls to the server
- `Call()`: Generic RPC call method
- `CallAdd()`: Specific method for addition
- `CallString()`: Specific method for string operations

## Message Format

### Request Format
```json
{
  "id": "req_1234567890",
  "method": "add",
  "params": {
    "a": 10,
    "b": 20
  }
}
```

### Response Format
```json
{
  "id": "req_1234567890",
  "result": {
    "sum": 30
  }
}
```

### Error Response Format
```json
{
  "id": "req_1234567890",
  "error": "Invalid request parameters"
}
```

## Extending the System

To add new RPC methods:

1. **Server Side**:
   - Add new request/response types
   - Add a new handler method
   - Update the switch statement in `handleRequest()`

2. **Client Side**:
   - Add new request/response types
   - Add a new convenience method (e.g., `CallMultiply()`)
   - Update the client to use the new method

## Error Handling

The system handles various error scenarios:
- Invalid JSON requests
- Unknown methods
- Invalid parameters
- Network timeouts
- Server errors

## Performance Considerations

- Uses NATS request-reply with 5-second timeout
- Unique request IDs prevent response confusion
- JSON marshaling/unmarshaling for flexibility
- Stateless design allows for horizontal scaling

## Security Notes

This is a basic example. For production use, consider:
- Authentication and authorization
- Request validation
- Rate limiting
- TLS encryption
- Message signing 