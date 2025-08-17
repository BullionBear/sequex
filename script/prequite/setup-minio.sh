#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Default MinIO configuration
DEFAULT_ENDPOINT="http://localhost:9000"
DEFAULT_ACCESS_KEY="minioadmin"
DEFAULT_SECRET_KEY="minioadmin"

echo "MinIO Setup Script"
echo "=================="
echo ""

# Check if mc is installed
if ! command -v mc &> /dev/null; then
    log_info "Installing MinIO Client (mc)..."
    wget -q https://dl.min.io/client/mc/release/linux-amd64/mc
    chmod +x mc
    sudo mv mc /usr/local/bin/
    log_success "MinIO Client installed successfully"
else
    log_info "MinIO Client already installed"
fi

# Get MinIO configuration
read -p "MinIO Endpoint [$DEFAULT_ENDPOINT]: " endpoint
endpoint=${endpoint:-$DEFAULT_ENDPOINT}

read -p "Access Key [$DEFAULT_ACCESS_KEY]: " access_key
access_key=${access_key:-$DEFAULT_ACCESS_KEY}

read -p "Secret Key [$DEFAULT_SECRET_KEY]: " secret_key
secret_key=${secret_key:-$DEFAULT_SECRET_KEY}

# Configure MinIO
log_info "Configuring MinIO client..."
mc config host add minio "$endpoint" "$access_key" "$secret_key"

# Test connection
log_info "Testing connection..."
if mc ls minio &> /dev/null; then
    log_success "MinIO connection successful!"
    
    # Create bucket if it doesn't exist
    if ! mc ls minio/sequex &> /dev/null; then
        log_info "Creating sequex bucket..."
        mc mb minio/sequex
        log_success "Bucket created successfully!"
    else
        log_info "Bucket already exists"
    fi
else
    log_error "Failed to connect to MinIO. Please check your configuration."
    exit 1
fi

log_success "MinIO setup completed successfully!"
echo ""
echo "You can now use the upload script:"
echo "  ./script/upload.sh all"
echo "  ./script/upload.sh sqx"
echo "  ./script/upload.sh master"
