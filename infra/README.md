# NATS JetStream Infrastructure

This directory contains the configuration files and documentation for managing NATS JetStream streams in the Sequex trading system.

## Overview

The system uses two main JetStream streams:

- **MARKET Stream**: Used for strategy execution and caching market data
- **RECORD Stream**: Used for recording and persistence of trading data

## Stream Configurations

### MARKET Stream (`stream_market.json`)
- **Purpose**: Real-time market data for strategy execution and caching
- **Subjects**: `market.trade.>`, `market.depth.>`
- **Storage**: Memory (for fast access)
- **Retention**: 5 seconds (short-term for real-time processing)
- **Acknowledgment**: Disabled (`no_ack: true`) for performance

### RECORD Stream (`stream_record.json`)
- **Purpose**: Long-term storage and recording of trading data
- **Subjects**: `record.trade.>`, `record.depth.>`
- **Storage**: File (for persistence)
- **Retention**: 1 hour (longer-term for analysis)
- **Acknowledgment**: Enabled for data integrity
- **Source**: Consumes from MARKET stream

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
nats-server --js
```

## Stream Management

### Create Streams

#### Option 1: Using JSON Configuration Files

```bash
# Create MARKET stream
nats stream add --config=stream_market.json

# Create RECORD stream
nats stream add --config=stream_record.json
```

#### Option 2: Using Command Line

```bash
# Create MARKET stream
nats stream add MARKET \
  --subjects="market.trade.>,market.depth.>" \
  --storage=memory \
  --retention=limits \
  --discard=old \
  --no-ack \
  --max-age=5s \
  --max-bytes=1GB

# Create RECORD stream
nats stream add RECORD \
  --subjects="record.trade.>,record.depth.>" \
  --storage=file \
  --retention=limits \
  --discard=old \
  --ack \
  --max-age=1h \
  --max-bytes=1GB \
  --sources=MARKET
```

### Delete Streams

```bash
# Delete MARKET stream (this will also delete RECORD stream due to dependency)
nats stream delete MARKET --force

# Or delete RECORD stream first, then MARKET
nats stream delete RECORD --force
nats stream delete MARKET --force
```

### Observe Streams

#### List All Streams
```bash
nats stream list
```

#### Get Stream Information
```bash
# Get detailed information about MARKET stream
nats stream info MARKET

# Get detailed information about RECORD stream
nats stream info RECORD
```

#### Monitor Stream Activity
```bash
# Monitor MARKET stream in real-time (new messages only)
nats subscribe --stream=MARKET --new

# Monitor RECORD stream in real-time (new messages only)
nats subscribe --stream=RECORD --new

# Monitor MARKET stream with all messages (including historical)
nats subscribe --stream=MARKET --all

# Monitor with specific subject pattern
nats subscribe market.trade.> --stream=MARKET --new

# Monitor with timestamps and rate graph
nats subscribe --stream=MARKET --new --timestamp --graph
```

#### View Stream Messages
```bash
# View recent messages from MARKET stream
nats stream view MARKET

# View recent messages from RECORD stream
nats stream view RECORD

# View messages with specific subject pattern
nats stream view MARKET --filter="market.trade.>"
```

#### Stream Statistics
```bash
# Get statistics for all streams
nats stream report

# Get statistics for specific stream
nats stream report MARKET
```

## Consumer Management

### Create Consumers

```bash
# Create consumer for MARKET stream
nats consumer add MARKET MARKET_CONSUMER \
  --deliver=all \
  --ack=explicit \
  --replay=instant

# Create consumer for RECORD stream
nats consumer add RECORD RECORD_CONSUMER \
  --deliver=all \
  --ack=explicit \
  --replay=instant
```

### List Consumers
```bash
# List consumers for MARKET stream
nats consumer list MARKET

# List consumers for RECORD stream
nats consumer list RECORD
```

### Consumer Information
```bash
# Get consumer information
nats consumer info MARKET MARKET_CONSUMER
nats consumer info RECORD RECORD_CONSUMER
```

## Data Flow

```
Market Data Sources → MARKET Stream → RECORD Stream
                           ↓              ↓
                    Strategy/Cache    Recording/Analysis
```

1. **Market Data** flows into the MARKET stream with subjects like `market.trade.BTCUSDT`
2. **Strategy components** consume from MARKET stream for real-time decision making
3. **RECORD stream** automatically receives data from MARKET stream for persistence
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
nats stream purge MARKET --filter="market.trade.>"

# Reset consumer
nats consumer reset MARKET MARKET_CONSUMER
```

## Configuration Details

### Stream Settings Explained

- **Storage**: `memory` for fast access, `file` for persistence
- **Retention**: `limits` means messages are removed based on age/bytes/msgs limits
- **Discard**: `old` removes oldest messages when limits are reached
- **Acknowledgment**: `no_ack` for performance, `ack` for reliability
- **Max Age**: Time before messages expire (5s for MARKET, 1h for RECORD)
- **Max Bytes**: Maximum storage size (1GB for both streams)
- **Subjects**: Wildcard patterns for message routing

### Performance Considerations

- **MARKET stream** uses memory storage for low latency
- **RECORD stream** uses file storage for durability
- **No acknowledgment** on MARKET stream for maximum throughput
- **Short retention** on MARKET stream to prevent memory bloat
- **Source relationship** ensures RECORD stream gets all MARKET data

## Integration with Application

The application code should:
1. **Error gracefully** if streams don't exist (don't create them in code)
2. **Use the predefined subject patterns** for publishing messages
3. **Handle stream unavailability** with appropriate retry logic
4. **Monitor stream health** and alert on issues

For more details on application integration, see the main project documentation.
