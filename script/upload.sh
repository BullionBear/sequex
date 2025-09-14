#!/bin/bash

set -e

# Upload script - packages bin/* into tar.gz and uploads to MinIO
BIN_DIR="bin"
BUCKET="yt"
PREFIX="sequex"
ARCHIVE_NAME="sqx-linux-amd64.tar.gz"
VERSION=$(git rev-parse --short HEAD)

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# Check if mc is installed
check_mc() {
    if ! command -v mc &> /dev/null; then
        log_info "Installing MinIO Client..."
        wget -q https://dl.min.io/client/mc/release/linux-amd64/mc
        chmod +x mc
        sudo mv mc /usr/local/bin/
        log_success "MinIO Client installed"
    fi
}

# Main upload logic
log_info "Starting upload process..."

# Check if bin directory exists and has files
if [ ! -d "$BIN_DIR" ] || [ -z "$(ls -A $BIN_DIR 2>/dev/null)" ]; then
    echo "Error: $BIN_DIR directory not found or empty"
    echo "Please run 'make build' first"
    exit 1
fi

check_mc

# Create tar.gz archive
log_info "Creating archive from $BIN_DIR/*"
tar -czf "$ARCHIVE_NAME" -C "$BIN_DIR" .

# Upload archive with version tag
log_info "Uploading $ARCHIVE_NAME to $BUCKET/$PREFIX/$ARCHIVE_NAME (version: $VERSION)"
mc cp "$ARCHIVE_NAME" "$BUCKET/$PREFIX/${ARCHIVE_NAME%.tar.gz}-${VERSION}.tar.gz"
mc cp "$ARCHIVE_NAME" "$BUCKET/$PREFIX/$ARCHIVE_NAME"

# Clean up local archive
rm -f "$ARCHIVE_NAME"

log_success "Upload completed!"
log_info "Download URLs:"
log_info "  Latest: $BUCKET/$PREFIX/$ARCHIVE_NAME"
log_info "  Version: $BUCKET/$PREFIX/${ARCHIVE_NAME%.tar.gz}-${VERSION}.tar.gz"
