# NATS JetStream Infrastructure

This directory contains the configuration files and documentation for managing NATS JetStream streams in the Sequex trading system.

## Overview

The system uses two main JetStream streams:

- **TRADE Stream**: Used for strategy execution and caching TRADE data
- **RECORD Stream**: Used for recording and persistence of trading data

## Stream Configurations

### TRADE Stream (`stream_TRADE.json`)
- **Purpose**: Real-time TRADE data for strategy execution and caching
- **Subjects**: `TRADE.trade.>`, `TRADE.depth.>`
- **Storage**: Memory (for fast access)
- **Retention**: 5 seconds (short-term for real-time processing)
- **Acknowledgment**: Disabled (`no_ack: true`) for performance

### RECORD Stream (`stream_record.json`)
- **Purpose**: Long-term storage and recording of trading data
- **Subjects**: `record.trade.>`, `record.depth.>`
- **Storage**: File (for persistence)
- **Retention**: 1 hour (longer-term for analysis)
- **Acknowledgment**: Enabled for data integrity
- **Source**: Consumes from TRADE stream

## Prerequisites

1. **NATS Server**: Ensure NATS server is running with JetStream enabled
2. **NATS CLI**: Install the NATS CLI tool for stream management

### Starting NATS Server

If using Docker (recommended):
```bash
# From the playground/jetstream directory
cd playground/jetstream
docker compose up -d
```

Or start NATS server directly:
```bash
curl -L https://github.com/nats-io/nats-server/releases/download/v2.11.9/nats-server-v2.11.9-linux-amd64.tar.gz -o nats-server-v2.11.9-linux-amd64.tar.gz
tar -xzf nats-server-v2.11.9-linux-amd64.tar.gz /bin/nats-server
mv nats-server-v2.11.9-linux-amd64/nats-server 
nats-server -js
```

## Stream Management

### Create Streams

#### Using JSON Configuration Files

```bash
# Create TRADE stream
nats stream add --config=stream_trade.json
```


### Delete Streams

```bash
# Delete TRADE stream (this will also delete RECORD stream due to dependency)
nats stream delete TRADE --force
```

### Observe Streams

#### List All Streams
```bash
nats stream list
```

#### Get Stream Information
```bash
# Get detailed information about TRADE stream
nats stream info TRADE

```

#### Monitor Stream Activity
```bash
# Monitor TRADE stream in real-time (new messages only)
nats subscribe --stream=TRADE --new

# Monitor TRADE stream with all messages (including historical)
nats subscribe --stream=TRADE --all

# Monitor with specific subject pattern
nats subscribe TRADE.trade.> --stream=TRADE --new

# Monitor with timestamps and rate graph
nats subscribe --stream=TRADE --new --timestamp --graph
```

#### View Stream Messages
```bash
# View recent messages from TRADE stream
nats stream view TRADE

# View messages with specific subject pattern
nats stream view TRADE --filter="TRADE.trade.>"
```

#### Stream Statistics
```bash
# Get statistics for all streams
nats stream report

# Get statistics for specific stream
nats stream report TRADE
```

## Consumer Management

### Create Consumers

```bash
# Create consumer for TRADE stream
nats consumer add TRADE --config consumer_trade_pubsub.json

# Create consumer for RECORD stream
nats consumer add TRADE --config consumer_trade_work_queue.json
```

### List Consumers
```bash
# List consumers for TRADE stream
nats consumer list TRADE
```

### Consumer Information
```bash
# Get consumer information
nats consumer info TRADE TRADE_CONSUMER
nats consumer info RECORD RECORD_CONSUMER
```

### Delete Consumer
```bash
# Delete consumer from TRADE stream
nats consumer rm TRADE TRADE_PUBSUB

# Delete consumer from TRADE stream
nats consumer rm TRADE TRADE_WORK_QUEUE
```


## Data Flow

```
TRADE Data Sources → TRADE Stream → RECORD Stream
                           ↓              ↓
                    Strategy/Cache    Recording/Analysis
```

1. **TRADE Data** flows into the TRADE stream with subjects like `TRADE.trade.BTCUSDT`
2. **Strategy components** consume from TRADE stream for real-time decision making
3. **RECORD stream** automatically receives data from TRADE stream for persistence
4. **Analysis tools** consume from RECORD stream for historical analysis

## Troubleshooting

### Common Issues

1. **Stream Creation Fails**
   ```bash
   # Check if NATS server is running
   nc -z localhost 4222
   
   # Check JetStream status
   nats server info jetstream
   ```

2. **Permission Denied**
   ```bash
   # Ensure NATS server has JetStream enabled
   nats-server --jetstream --store_dir=./data
   ```

3. **Stream Not Found**
   ```bash
   # List all streams to verify existence
   nats stream list
   ```

### Useful Commands

```bash
# Check NATS server status
nats server info

# Check JetStream status
nats server info jetstream

# Purge stream messages
nats stream purge TRADE --filter="TRADE.trade.>"

# Reset consumer
nats consumer reset TRADE TRADE_CONSUMER
```

## Configuration Details

### Stream Settings Explained

- **Storage**: `memory` for fast access, `file` for persistence
- **Retention**: `limits` means messages are removed based on age/bytes/msgs limits
- **Discard**: `old` removes oldest messages when limits are reached
- **Acknowledgment**: `no_ack` for performance, `ack` for reliability
- **Max Age**: Time before messages expire (5s for TRADE, 1h for RECORD)
- **Max Bytes**: Maximum storage size (1GB for both streams)
- **Subjects**: Wildcard patterns for message routing

### Performance Considerations

- **TRADE stream** uses memory storage for low latency
- **RECORD stream** uses file storage for durability
- **No acknowledgment** on TRADE stream for maximum throughput
- **Short retention** on TRADE stream to prevent memory bloat
- **Source relationship** ensures RECORD stream gets all TRADE data

## Integration with Application

The application code should:
1. **Error gracefully** if streams don't exist (don't create them in code)
2. **Use the predefined subject patterns** for publishing messages
3. **Handle stream unavailability** with appropriate retry logic
4. **Monitor stream health** and alert on issues

For more details on application integration, see the main project documentation.
