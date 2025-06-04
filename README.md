# IAM Credentials Proxy

A lightweight proxy server that exposes AWS IAM credentials from a Kubernetes pod using IRSA (IAM Role for ServiceAccount) through a simple HTTP API compatible with AWS credential_process.

## Features

- Exposes temporary AWS credentials via HTTP API
- Compatible with AWS credential_process in .aws/config
- Supports optional security features:
  - Bearer token authentication
  - IP whitelisting
- Health check endpoint
- Graceful shutdown
- Minimal container image

## Usage

### Running in Kubernetes

1. Deploy the proxy with IRSA configuration:

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: iam-proxy
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::ACCOUNT_ID:role/YOUR_ROLE_NAME
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: iam-proxy
spec:
  replicas: 1
  selector:
    matchLabels:
      app: iam-proxy
  template:
    metadata:
      labels:
        app: iam-proxy
    spec:
      serviceAccountName: iam-proxy
      containers:
      - name: iam-proxy
        image: ghcr.io/kotaicode/iam-proxy:latest
        ports:
        - containerPort: 8080
        env:
        - name: PORT
          value: "8080"
        - name: SECURITY_TOKEN
          value: "your-secure-token"  # Optional
        - name: ALLOWED_IPS
          value: "10.0.0.0/8,172.16.0.0/12"  # Optional
```

### Configuration

Environment variables:

- `PORT`: Server port (default: 8080)
- `LOG_LEVEL`: Logging level (default: info)
- `SECURITY_TOKEN`: Bearer token for authentication (optional)
- `ALLOWED_IPS`: Comma-separated list of allowed IP addresses (optional)

### AWS Configuration

Add to your `~/.aws/config`:

```ini
[profile irsa-proxy]
credential_process = curl -s -H "Authorization: Bearer your-secure-token" http://localhost:8080/credentials
region = eu-central-1
```

## Development

### Prerequisites

- Go 1.21 or later
- Docker
- Task (optional, for using Taskfile)

### Building

```bash
# Using Task
task build

# Or manually
go build -o bin/iam-proxy ./cmd/main.go
```

### Running Locally

```bash
# Using Task
task run

# Or manually
./bin/iam-proxy
```

### Testing

```bash
task test
```

### Docker

```bash
# Build image
task docker

# Push to GitHub Container Registry
task docker-push
```

### Creating a Release

1. Install GoReleaser:
```bash
brew install goreleaser
```

2. Test the release process:
```bash
task release-dry-run
```

3. Create a new release:
```bash
VERSION=0.1.0 task release
```

This will:
- Create and push a version tag
- Trigger GitHub Actions workflow
- Build binaries for multiple platforms
- Create Docker images
- Create a GitHub release with all artifacts

## Security Considerations

- Always use HTTPS in production
- Configure IP whitelisting when possible
- Use a strong security token
- Consider using client certificates for additional security
- Monitor access logs for suspicious activity 