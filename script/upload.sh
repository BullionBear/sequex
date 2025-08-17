#!/bin/bash

set -e

# Simple upload script for current binary to S3
BINARY="bin/sqx-linux-amd64"
BUCKET="minio"
PREFIX="sequex"
VERSION=$(git rev-parse --short HEAD)

echo "Uploading binary to S3..."

# Check if binary exists
if [ ! -f "$BINARY" ]; then
    echo "Error: Binary not found at $BINARY"
    echo "Please run 'make build' first"
    exit 1
fi

# Check if mc is installed
if ! command -v mc &> /dev/null; then
    echo "Installing MinIO Client..."
    wget -q https://dl.min.io/client/mc/release/linux-amd64/mc
    chmod +x mc
    sudo mv mc /usr/local/bin/
fi

# Upload binary
echo "Uploading $BINARY to $BUCKET/$PREFIX/sqx-linux-amd64:$VERSION"
mc cp "$BINARY" "$BUCKET/$PREFIX/sqx-linux-amd64:$VERSION"
mc cp "$BINARY" "$BUCKET/$PREFIX/sqx-linux-amd64:latest"

echo "Upload completed!"
echo "Download URL: $BUCKET/$PREFIX/sqx-linux-amd64:latest"