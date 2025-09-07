# JetStream Example

This directory contains a complete example demonstrating NATS JetStream functionality with Docker, including stream creation, message publishing, and consumption.

## Prerequisites

- Docker and Docker Compose
- Go 1.23+ (for running the example)
- NATS CLI tool (for script operations)
- This example uses the existing sequex project dependencies

### Installing NATS CLI

```bash
go install github.com/nats-io/natscli/nats@latest
```

## Quick Start

### 1. Start NATS with JetStream

```bash
# Start NATS server with JetStream enabled
docker compose up -d

# Verify NATS is running
docker compose ps
```

The NATS server will be available at:
- **Client Port**: `nats://localhost:4222`
- **HTTP Monitoring**: `http://localhost:8222`
- **Cluster Port**: `nats://localhost:6222`

### 2. Set up JetStream Stream and Consumer

```bash
# Create test stream and consumer
./script.sh create
```

This will create:
- **Stream**: `TEST_STREAM` with subject pattern `test.>`
- **Consumer**: `TEST_CONSUMER` with explicit acknowledgment

### 3. Run the Go Example

```bash
# Run the JetStream example (from the jetstream directory)
cd playground/jetstream
go run main.go
```

The example will:
- Connect to NATS JetStream
- Verify the stream and consumer exist (created by script.sh)
- Publish 5 test messages
- Consume and acknowledge the messages
- Display stream and consumer statistics

**Note**: The Go example assumes the stream and consumer already exist. If they don't exist, it will show an error message and exit. Always run `./script.sh create` first.

### 4. Clean Up

```bash
# Remove test stream and consumer
./script.sh clean

# Stop NATS server
docker-compose down
```

## Files Overview

### `docker-compose.yml`
- Configures NATS server with JetStream enabled
- Sets up persistent storage for streams
- Exposes monitoring and client ports
- Includes health checks

### `script.sh`
- **`create`**: Creates test stream and consumer
- **`clean`**: Removes test stream and consumer
- Includes connection validation and error handling

### `main.go`
- Complete Go example demonstrating JetStream usage
- Shows stream creation, message publishing, and consumption
- Includes proper error handling and acknowledgment
- Displays stream and consumer statistics

## JetStream Configuration

### Stream Configuration
- **Name**: `TEST_STREAM`
- **Subjects**: `test.>`
- **Storage**: File-based storage
- **Retention**: Limits policy
- **Max Age**: 1 hour
- **Max Messages**: 1000
- **Max Bytes**: 1MB
- **Replicas**: 1

### Consumer Configuration
- **Name**: `TEST_CONSUMER`
- **Durable**: Yes
- **Deliver Policy**: All messages
- **Ack Policy**: Explicit acknowledgment
- **Replay Policy**: Instant

## Monitoring

### NATS Server Monitoring
Access the NATS monitoring dashboard at: http://localhost:8222

### Stream Information
```bash
# View stream information
nats stream info TEST_STREAM

# List all streams
nats stream list
```

### Consumer Information
```bash
# View consumer information
nats consumer info TEST_STREAM TEST_CONSUMER

# List consumers for a stream
nats consumer list TEST_STREAM
```

## Advanced Usage

### Manual Stream Creation
```bash
nats stream add TEST_STREAM \
  --subjects="test.>" \
  --storage=file \
  --retention=limits \
  --max-age=1h \
  --max-msgs=1000 \
  --max-bytes=1MB \
  --replicas=1
```

### Manual Consumer Creation
```bash
nats consumer add TEST_STREAM TEST_CONSUMER \
  --deliver=all \
  --ack=explicit \
  --replay=instant
```

### Publishing Messages
```bash
# Publish a single message
nats pub test.message.1 "Hello JetStream!"

# Publish multiple messages
for i in {1..5}; do
  nats pub test.message.$i "Message $i"
done
```

### Consuming Messages
```bash
# Pull messages from consumer
nats consumer next TEST_STREAM TEST_CONSUMER

# Subscribe to messages
nats sub test.>
```

## Troubleshooting

### Connection Issues
- Ensure NATS server is running: `docker-compose ps`
- Check server logs: `docker-compose logs nats`
- Verify port accessibility: `nc -z localhost 4222`

### Stream/Consumer Issues
- Check if stream exists: `nats stream list`
- Verify consumer configuration: `nats consumer list TEST_STREAM`
- Review server logs for errors

### Go Dependencies
This example uses the existing sequex project dependencies. The NATS Go client is already included in the main project's go.mod file.

```bash
# Verify NATS Go client version (from project root)
go list -m github.com/nats-io/nats.go
```

## Example Output

```
🚀 Starting JetStream Example
==============================
✅ Connected to NATS server
✅ JetStream context created
✅ Stream 'TEST_STREAM' already exists
✅ Consumer 'TEST_CONSUMER' already exists

📤 Publishing test messages...
  📨 Published to test.message.1: Hello JetStream! Message #1 - 2024-01-15T10:30:00Z
  📨 Published to test.message.2: Hello JetStream! Message #2 - 2024-01-15T10:30:00Z
  📨 Published to test.message.3: Hello JetStream! Message #3 - 2024-01-15T10:30:00Z
  📨 Published to test.message.4: Hello JetStream! Message #4 - 2024-01-15T10:30:00Z
  📨 Published to test.message.5: Hello JetStream! Message #5 - 2024-01-15T10:30:00Z

📥 Consuming messages...
📬 Received 5 messages:
  1. Subject: test.message.1
     Data: Hello JetStream! Message #1 - 2024-01-15T10:30:00Z
     Timestamp: 2024-01-15T10:30:00.123Z
     ✅ Acknowledged

  2. Subject: test.message.2
     Data: Hello JetStream! Message #2 - 2024-01-15T10:30:00Z
     Timestamp: 2024-01-15T10:30:00.223Z
     ✅ Acknowledged

📊 Stream Information:
  Name: TEST_STREAM
  Subjects: [test.>]
  Storage: File
  Retention: Limits
  Messages: 5
  Bytes: 250

📊 Consumer Information:
  Name: TEST_CONSUMER
  Durable: TEST_CONSUMER
  Deliver Policy: All
  Ack Policy: Explicit
  Num Pending: 0
  Num Delivered: 5

🎉 JetStream example completed successfully!
💡 Use './script.sh clean' to clean up the test stream and consumer
```

## Learn More

- [NATS JetStream Documentation](https://docs.nats.io/jetstream)
- [NATS Go Client](https://github.com/nats-io/nats.go)
- [NATS CLI Tool](https://github.com/nats-io/natscli)
