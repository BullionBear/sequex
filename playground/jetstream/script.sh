#!/bin/bash

# JetStream Management Script
# Usage: ./script.sh [create|clean]

NATS_URL="nats://localhost:4222"
STREAM_NAME="TEST_STREAM"
CONSUMER_NAME="TEST_CONSUMER"
SUBJECT="test.>"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}


print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_nats_connection() {
    print_info "Checking NATS connection..."
    if ! nc -z localhost 4222; then
        print_error "NATS server is not running on localhost:4222"
        print_info "Please start NATS with: docker-compose up -d"
        exit 1
    fi
    print_info "NATS connection successful"
}

create_stream() {
    print_info "Creating JetStream stream: $STREAM_NAME"
    
    # Create stream using command line flags with all required parameters
    nats stream add $STREAM_NAME \
        --subjects="$SUBJECT" \
        --storage=file \
        --retention=limits \
        --discard=old \
        --ack \
        --max-age=1h \
        --max-msgs=1000 \
        --max-bytes=1MB \
        --defaults \
        --server="$NATS_URL"
    
    if [ $? -eq 0 ]; then
        print_info "Stream '$STREAM_NAME' created successfully"
    else
        print_error "Failed to create stream '$STREAM_NAME'"
        exit 1
    fi
}

create_consumer() {
    print_info "Creating consumer: $CONSUMER_NAME"
    
    # Create consumer using command line flags
    nats consumer add $STREAM_NAME $CONSUMER_NAME \
        --deliver=all \
        --ack=explicit \
        --replay=instant \
        --defaults \
        --server="$NATS_URL"
    
    if [ $? -eq 0 ]; then
        print_info "Consumer '$CONSUMER_NAME' created successfully"
    else
        print_error "Failed to create consumer '$CONSUMER_NAME'"
        exit 1
    fi
}

show_stream_info() {
    print_info "Stream information:"
    nats stream info $STREAM_NAME --server="$NATS_URL"
    
    print_info "Consumer information:"
    nats consumer info $STREAM_NAME $CONSUMER_NAME --server="$NATS_URL"
}

clean_stream() {
    print_warning "Removing JetStream stream: $STREAM_NAME"
    
    nats stream delete $STREAM_NAME --force --server="$NATS_URL"
    
    if [ $? -eq 0 ]; then
        print_info "Stream '$STREAM_NAME' removed successfully"
    else
        print_error "Failed to remove stream '$STREAM_NAME'"
        exit 1
    fi
}

case "$1" in
    "create")
        print_info "Setting up JetStream test environment..."
        check_nats_connection
        create_stream
        create_consumer
        show_stream_info
        print_info "JetStream test environment ready!"
        print_info "You can now run the Go example: go run main.go"
        ;;
    "clean")
        print_info "Cleaning up JetStream test environment..."
        check_nats_connection
        clean_stream
        print_info "JetStream test environment cleaned up!"
        ;;
    *)
        echo "Usage: $0 {create|clean}"
        echo ""
        echo "Commands:"
        echo "  create  - Create test JetStream stream and consumer"
        echo "  clean   - Remove test JetStream stream and consumer"
        exit 1
        ;;
esac
