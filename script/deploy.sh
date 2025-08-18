#!/bin/bash

set -e

# Deploy script - combines upload and install functionality
BINARY="bin/sqx-linux-amd64"
BUCKET="minio"
PREFIX="sequex"
BINARY_NAME="sqx-linux-amd64"
INSTALL_PATH="/usr/local/bin/sqx"
TEMP_PATH="/tmp/sqx-linux-amd64"
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

# Upload function
upload() {
    log_info "Starting upload process..."
    
    # Check if binary exists
    if [ ! -f "$BINARY" ]; then
        echo "Error: Binary not found at $BINARY"
        echo "Please run 'make build' first"
        exit 1
    fi
    
    check_mc
    
    # Upload binary
    log_info "Uploading $BINARY to $BUCKET/$PREFIX/$BINARY_NAME:$VERSION"
    mc cp "$BINARY" "$BUCKET/$PREFIX/$BINARY_NAME:$VERSION"
    mc cp "$BINARY" "$BUCKET/$PREFIX/$BINARY_NAME:latest"
    
    log_success "Upload completed!"
    log_info "Download URL: $BUCKET/$PREFIX/$BINARY_NAME:latest"
}

# Install function
install() {
    log_info "Starting installation process..."
    
    check_mc
    
    # Download the latest binary
    log_info "Downloading latest binary from $BUCKET/$PREFIX/$BINARY_NAME:latest"
    mc cp "$BUCKET/$PREFIX/$BINARY_NAME:latest" "$TEMP_PATH"
    
    # Make it executable
    chmod +x "$TEMP_PATH"
    
    # Move to /usr/local/bin
    log_info "Installing to $INSTALL_PATH"
    sudo mv "$TEMP_PATH" "$INSTALL_PATH"
    
    log_success "Installation completed!"
    log_info "You can now run: sqx --help"
}

# Show usage
show_usage() {
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  upload    Build and upload binary to S3"
    echo "  install   Download and install latest binary from S3"
    echo "  deploy    Upload and then install (full deployment)"
    echo "  help      Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 upload   # Upload current binary to S3"
    echo "  $0 install  # Install latest binary from S3"
    echo "  $0 deploy   # Upload and install (full deployment)"
}

# Deploy function (upload + install)
deploy() {
    log_info "Starting full deployment process..."
    upload
    echo ""
    install
    log_success "Full deployment completed!"
}

# Main script logic
main() {
    local command="${1:-help}"
    
    case "$command" in
        "upload")
            upload
            ;;
        "install")
            install
            ;;
        "deploy")
            deploy
            ;;
        "help"|"-h"|"--help")
            show_usage
            ;;
        *)
            echo "Error: Unknown command '$command'"
            show_usage
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
