#!/bin/bash

set -e

# Install script for downloading and installing the latest binary
BUCKET="minio"
PREFIX="sequex"
BINARY_NAME="sqx-linux-amd64"
INSTALL_PATH="/usr/local/bin/sqx"
TEMP_PATH="/tmp/sqx-linux-amd64"

echo "Installing sequex binary..."

# Check if mc is installed
if ! command -v mc &> /dev/null; then
    echo "Installing MinIO Client..."
    wget -q https://dl.min.io/client/mc/release/linux-amd64/mc
    chmod +x mc
    sudo mv mc /usr/local/bin/
fi

# Download the latest binary
echo "Downloading latest binary from $BUCKET/$PREFIX/$BINARY_NAME:latest"
mc cp "$BUCKET/$PREFIX/$BINARY_NAME:latest" "$TEMP_PATH"

# Make it executable
chmod +x "$TEMP_PATH"

# Move to /usr/local/bin
echo "Installing to $INSTALL_PATH"
sudo mv "$TEMP_PATH" "$INSTALL_PATH"

echo "Installation completed!"
echo "You can now run: sqx --help"
