#!/bin/bash

set -e

# Configuration
REGISTRY="ghcr.io"
OWNER="bullionbear"
IMAGE_NAME="sequex"
FULL_IMAGE="${REGISTRY}/${OWNER}/${IMAGE_NAME}"

# Get commit hash for tagging
COMMIT_HASH=$(git rev-parse --short HEAD)

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

# Ensure we're logged in to GHCR
log_info "Checking GHCR authentication..."
if ! sudo docker pull ${REGISTRY}/${OWNER}/sequex-base:latest > /dev/null 2>&1; then
    echo "Please login to GHCR first:"
    echo "  echo \$GITHUB_TOKEN | sudo docker login ghcr.io -u USERNAME --password-stdin"
    exit 1
fi

# Build the image
log_info "Building Docker image..."
sudo docker build -t "${FULL_IMAGE}:${COMMIT_HASH}" -t "${FULL_IMAGE}:latest" .

log_success "Image built with tags:"
log_info "  - ${FULL_IMAGE}:${COMMIT_HASH}"
log_info "  - ${FULL_IMAGE}:latest"

# Push to GHCR
log_info "Pushing ${FULL_IMAGE}:${COMMIT_HASH}..."
sudo docker push "${FULL_IMAGE}:${COMMIT_HASH}"

log_info "Pushing ${FULL_IMAGE}:latest..."
sudo docker push "${FULL_IMAGE}:latest"

log_success "Successfully pushed images to GHCR!"
