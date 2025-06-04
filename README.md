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

## Installation

### Using Helm

1. Add the Helm repository:
```bash
helm repo add kotaicode https://kotaico.de/iam-proxy
helm repo update
```

2. Create a values file (e.g., `my-values.yaml`):
```yaml
serviceAccount:
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::ACCOUNT_ID:role/YOUR_ROLE_NAME

config:
  securityToken: "your-secure-token"  # Optional
  allowedIPs:  # Optional
    - "10.0.0.0/8"
    - "172.16.0.0/12"
```

3. Install the chart:
```bash
helm install iam-proxy kotaicode/iam-proxy -f my-values.yaml
```

For more configuration options, see the [values.yaml](helm/iam-proxy/values.yaml) file.

### Using Docker

You can run the proxy directly using Docker. Choose the appropriate image for your architecture:

For AMD64:
```bash
docker run -p 8080:8080 \
  -e SECURITY_TOKEN=your-secure-token \
  -e ALLOWED_IPS=10.0.0.0/8,172.16.0.0/12 \
  ghcr.io/kotaicode/iam-proxy:latest-amd64
```

For ARM64:
```bash
docker run -p 8080:8080 \
  -e SECURITY_TOKEN=your-secure-token \
  -e ALLOWED_IPS=10.0.0.0/8,172.16.0.0/12 \
  ghcr.io/kotaicode/iam-proxy:latest-arm64
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
- Helm (for local development)

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

### Local Helm Development

1. Install dependencies:
```bash
helm dependency update helm/iam-proxy
```

2. Test the chart:
```bash
helm lint helm/iam-proxy
helm template iam-proxy helm/iam-proxy -f examples/values.yaml
```

3. Install locally:
```bash
helm install iam-proxy helm/iam-proxy -f examples/values.yaml
```

## Security Considerations

- Always use HTTPS in production
- Configure IP whitelisting when possible
- Use a strong security token
- Consider using client certificates for additional security
- Monitor access logs for suspicious activity

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 