# Prerequisites for Operations

This document covers the prerequisites needed for operations, including uploading binaries from local to a self-hosted MinIO instance.

## 1. Install MinIO CLI (mc)

### Download and Install
```bash
# Download MinIO CLI for Linux AMD64
curl https://dl.min.io/client/mc/release/linux-amd64/mc \
  -o ./mc

# Make it executable
chmod +x ./mc
sudo mv ./mc /usr/local/bin/

# Verify installation
mc --help
```

## 2. Create Access Key and Secret Key

### For MinIO Server (Self-hosted)

If you're running your own MinIO server, you can create access keys through the MinIO Console (web interface) or using the MinIO CLI.

#### Method 1: Using MinIO Console (Web Interface)
1. Open MinIO Console in your browser (typically `http://your-server:9001`)
2. Login with your root credentials
3. Navigate to **Access Keys** section
4. Click **Create Access Key**
5. Copy the generated **Access Key** and **Secret Key**
6. Optionally set an expiration date and policies

#### Method 2: Using MinIO Admin CLI
```bash
# First configure admin access (use root credentials)
mc alias set myadmin http://your-minio-server:9000 root-user root-password
```

#### Method 3: Using Environment Variables (Docker)
If you're running MinIO in Docker, you can set custom root credentials:

```bash
# In your docker-compose.yml or docker run command
environment:
  MINIO_ROOT_USER: your-custom-access-key
  MINIO_ROOT_PASSWORD: your-custom-secret-key
```

### Default Credentials

Many MinIO installations use these default credentials:
- **Access Key**: `minioadmin`
- **Secret Key**: `minioadmin`

⚠️ **Security Note**: Always change default credentials in production environments!

### Best Practices for Access Keys

1. **Use Service Accounts**: Create dedicated service accounts instead of using root credentials
2. **Principle of Least Privilege**: Grant only necessary permissions
3. **Rotate Keys Regularly**: Set expiration dates and rotate keys periodically
4. **Store Securely**: Never hardcode keys in your source code
5. **Use Environment Variables**: Store keys in environment variables or secure vaults

```bash
# Example: Using environment variables
export MINIO_ACCESS_KEY="your-access-key"
export MINIO_SECRET_KEY="your-secret-key"

# Then use in mc alias command
mc alias set myminio http://your-server:9000 $MINIO_ACCESS_KEY $MINIO_SECRET_KEY
```

## 3. Setup MinIO CLI with Remote MinIO Server

### Configure MinIO Alias
```bash
# Add your MinIO server configuration
# Replace with your actual MinIO server details
mc alias set myminio http://your-minio-server:9000 your-access-key your-secret-key

# Example:
# mc alias set myminio http://localhost:9000 minioadmin minioadmin

# Verify the connection
mc admin info myminio
```

### List Available Aliases
```bash
# Show all configured aliases
mc alias list
```

## 4. Upload Files to MinIO Bucket

### Basic File Upload
```bash
# Upload a single file
mc cp /path/to/local/file myminio/bucket-name/

# Example: Upload binary to releases bucket
mc cp ./bin/sqx-linux-amd64 myminio/releases/
```

### Upload File with Tags
```bash
# Upload file with custom tags
mc cp --tags "version=v1.0.0,env=production,type=binary" \
  ./bin/sqx-linux-amd64 \
  myminio/releases/sqx-linux-amd64

# Upload with multiple tags
mc cp --tags "version=$VERSION,build=$(date +%Y%m%d),arch=amd64" \
  ./bin/sqx-linux-amd64 \
  myminio/releases/sqx-linux-amd64-$VERSION
```

### Advanced Upload Operations
```bash
# Upload with metadata
mc cp --attr "Content-Type=application/octet-stream,X-Custom-Header=value" \
  ./bin/sqx-linux-amd64 \
  myminio/releases/

# Upload directory recursively
mc cp --recursive ./bin/ myminio/releases/binaries/

# Upload with progress bar
mc cp --progress ./bin/sqx-linux-amd64 myminio/releases/
```

## 5. Bucket Management

### Create Bucket
```bash
# Create a new bucket
mc mb myminio/releases

# Create bucket with region
mc mb --region us-east-1 myminio/releases
```

### List Buckets and Objects
```bash
# List all buckets
mc ls myminio

# List objects in a bucket
mc ls myminio/releases

# List objects with details (size, date, etc.)
mc ls --recursive myminio/releases
```

### Set Bucket Policy (Optional)
```bash
# Set public read policy for releases bucket
mc policy set download myminio/releases

# Set custom policy from file
mc policy set-json /path/to/policy.json myminio/releases
```

## 6. Common Operations

### Download Files
```bash
# Download a file
mc cp myminio/releases/sqx-linux-amd64 ./downloads/

# Download with original tags preserved
mc cp --preserve myminio/releases/sqx-linux-amd64 ./downloads/
```

### View File Information
```bash
# Show file stats and metadata
mc stat myminio/releases/sqx-linux-amd64

# Show file tags
mc tag list myminio/releases/sqx-linux-amd64
```

### Remove Files
```bash
# Remove a file
mc rm myminio/releases/old-binary

# Remove recursively
mc rm --recursive myminio/releases/old-versions/
```

## 7. Environment Variables (Optional)

You can set environment variables to avoid passing credentials repeatedly:

```bash
# Set MinIO server configuration via environment
export MC_HOST_myminio=http://your-access-key:your-secret-key@your-minio-server:9000

# Now you can use mc commands without alias setup
mc ls myminio/
```

## 8. Troubleshooting

### Common Issues
1. **Permission denied**: Ensure the binary has execute permissions (`chmod +x`)
2. **Connection refused**: Verify MinIO server is running and accessible
3. **Access denied**: Check access key and secret key are correct
4. **Bucket not found**: Create the bucket first with `mc mb`

### Debug Mode
```bash
# Run mc commands with debug output
mc --debug cp file myminio/bucket/
```

### Check Configuration
```bash
# Verify alias configuration
mc alias list myminio

# Test connection
mc admin info myminio
```

## Example Workflow

Here's a complete example workflow for uploading a binary:

```bash
# 1. Install mc (if not already installed)
curl https://dl.min.io/client/mc/release/linux-amd64/mc \
  --create-dirs \
  -o /usr/local/bin/mc
chmod +x /usr/local/bin/mc

# 2. Configure MinIO alias
mc alias set myminio http://localhost:9000 minioadmin minioadmin

# 3. Create bucket (if it doesn't exist)
mc mb myminio/releases

# 4. Upload binary with version tag
VERSION=$(git describe --tags --always)
mc cp --tags "version=$VERSION,type=binary,arch=amd64" \
  ./bin/sqx-linux-amd64 \
  myminio/releases/sqx-linux-amd64-$VERSION

# 5. Verify upload
mc ls myminio/releases/
mc stat myminio/releases/sqx-linux-amd64-$VERSION
```
