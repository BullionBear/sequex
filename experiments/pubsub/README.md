# NATS PubSub Example

This directory contains a simple NATS pubsub example with a publisher and subscriber.

## Prerequisites

1. Install NATS server:
   ```bash
   # Using Homebrew (macOS)
   brew install nats-io/nats-tools/nats
   
   # Or download from https://nats.io/download/
   ```

2. Start NATS server:
   ```bash
   nats-server
   ```

## Running the Example

### 1. Start the Subscriber

In one terminal, run the subscriber:

```bash
cd experiments/pubsub/sub
go run main.go
```

You should see output like:
```
Connected to NATS at nats://localhost:4222
Subscribed to sequex.trades
```

### 2. Start the Publisher

In another terminal, run the publisher:

```bash
cd experiments/pubsub/pub
go run main.go
```

You should see output like:
```
Connected to NATS at nats://localhost:4222
Published: Trade message #1 at 2024-01-15T10:30:00Z
Published: Trade message #2 at 2024-01-15T10:30:02Z
...
```

### 3. Observe the Communication

The subscriber will receive and display the messages published by the publisher:

```
Received message: Trade message #1 at 2024-01-15T10:30:00Z
Received message: Trade message #2 at 2024-01-15T10:30:02Z
...
```

## Features

- **Graceful Shutdown**: Both publisher and subscriber handle SIGINT/SIGTERM signals
- **Error Handling**: Proper error handling for connection and message operations
- **Message Acknowledgment**: Subscriber acknowledges received messages
- **Flush Operations**: Publisher ensures messages are sent before continuing
- **Structured Logging**: Clear logging of all operations

## Configuration

- **NATS URL**: Defaults to `nats://localhost:4222` (can be changed in the constants)
- **Subject**: Messages are published to `sequex.trades`
- **Publish Interval**: Publisher sends messages every 2 seconds

## Stopping the Example

Press `Ctrl+C` in either terminal to gracefully stop the publisher or subscriber. 