# Multi-stage Dockerfile using custom base image from GHCR
# Stage 1: Build stage using custom base image with Go and protobuf pre-installed
FROM ghcr.io/bullionbear/sequex-base:latest AS builder

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Protobuf tools are already installed in the base image

# Copy source code
COPY . .

# Build the application
RUN make build

# Stage 2: Runtime stage using Debian slim (glibc compatible)
FROM debian:bookworm-slim AS runtime

# Install ca-certificates
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN groupadd -g 1001 appgroup && \
    useradd -u 1001 -g appgroup -s /bin/false appuser

# Set working directory
WORKDIR /app

# Copy binaries from builder stage
COPY --from=builder /app/bin/feed-linux-amd64 /app/feed

# Copy configuration files
COPY --from=builder /app/config/ /app/config/

# Make binaries executable and change ownership
RUN chmod +x /app/feed && \
    chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Default command
CMD ["./feed"]