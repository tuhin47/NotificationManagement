# Notification Management System

A Go-based notification management system with AWS Config service integration for development and production environments.

## Features

- AWS Config service integration
- LocalStack support for development
- Configuration management with environment variables
- CLI commands for AWS Config operations

## Prerequisites

- Go 1.24 or higher
- Docker and Docker Compose (for LocalStack)
- AWS CLI (optional, for production)

## Quick Start

### 1. Start LocalStack for Development

```bash
# Start LocalStack with AWS Config service
docker-compose up -d

# Verify LocalStack is running
curl http://localhost:4566/health
```

### 2. Build and Run the Application

```bash
# Build the application
go build -o notification-management

# Run the application
./notification-management serve
```

### 3. Test AWS Config Service

```bash
# Test connection to AWS Config service
./notification-management config test-connection

# List AWS Config rules
./notification-management config list-rules

# Check configuration recorder status
./notification-management config check-status
```

## Configuration

The application uses a hierarchical configuration system:

1. **Default values** (hardcoded in the application)
2. **Configuration file** (`config.yaml`)
3. **Environment variables** (override file settings)

### Configuration File

Create a `config.yaml` file in the project root:

```yaml
app:
  name: "NotificationManagement"
  version: "1.0.0"
  port: 8080
  env: "development"

aws:
  region: "us-east-1"
  access_key_id: "test"
  secret_access_key: "test"
  endpoint: "http://localhost:4566"
  use_localstack: true
  config_service:
    enabled: true

# ... other configurations
```

### Environment Variables

You can override configuration using environment variables:

```bash
export AWS_REGION=us-west-2
export AWS_USE_LOCALSTACK=true
export AWS_ENDPOINT=http://localhost:4566
export APP_ENV=development
```

## LocalStack Development Setup

### Starting LocalStack

```bash
# Start all services
docker-compose up -d

# Start only LocalStack (without web UI)
docker-compose up -d localstack
```

### LocalStack Web UI

Access the LocalStack web interface at: http://localhost:8080

### AWS CLI with LocalStack

Configure AWS CLI to use LocalStack:

```bash
aws configure set aws_access_key_id test
aws configure set aws_secret_access_key test
aws configure set region us-east-1
aws configure set output json

# Test with LocalStack
aws --endpoint-url=http://localhost:4566 configservice list-config-rules
```

## AWS Config Service Usage

### Available Commands

```bash
# Test connection to AWS Config service
./notification-management config test-connection

# List all AWS Config rules
./notification-management config list-rules

# Check configuration recorder status
./notification-management config check-status
```

### Programmatic Usage

```go
package main

import (
    "context"
    "NotificationManagement/aws"
)

func main() {
    // Create AWS Config client
    client, err := aws.NewConfigClient()
    if err != nil {
        panic(err)
    }

    // List config rules
    ctx := context.Background()
    rules, err := client.ListConfigRules(ctx)
    if err != nil {
        panic(err)
    }

    // Process rules...
}
```

## Development Workflow

### 1. Start LocalStack

```bash
docker-compose up -d
```

### 2. Set Environment Variables

```bash
export AWS_USE_LOCALSTACK=true
export AWS_ENDPOINT=http://localhost:4566
export AWS_REGION=us-east-1
```

### 3. Run Application

```bash
go run main.go serve
```

### 4. Test AWS Config Operations

```bash
go run main.go config test-connection
go run main.go config list-rules
```

## Production Deployment

For production, set the following environment variables:

```bash
export AWS_USE_LOCALSTACK=false
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=your-access-key
export AWS_SECRET_ACCESS_KEY=your-secret-key
```

## Project Structure

```
NotificationManagement/
├── aws/                    # AWS service clients
│   └── config_client.go   # AWS Config service client
├── cmd/                    # CLI commands
│   ├── root.go            # Root command
│   ├── serve.go           # Serve command
│   ├── worker.go          # Worker command
│   └── config.go          # AWS Config commands
├── config/                 # Configuration management
│   └── config.go          # Configuration structure and loading
├── docker-compose.yml      # LocalStack development setup
├── config.yaml            # Configuration file
├── go.mod                 # Go module dependencies
├── main.go                # Application entry point
└── README.md              # This file
```

## Troubleshooting

### LocalStack Issues

1. **LocalStack not starting**: Check Docker is running and ports are available
2. **Connection refused**: Ensure LocalStack is fully started (check logs with `docker-compose logs localstack`)
3. **AWS Config service not available**: Verify `config` is in the SERVICES environment variable

### AWS Config Service Issues

1. **Authentication errors**: Check AWS credentials and region configuration
2. **Permission errors**: Ensure proper IAM permissions for AWS Config operations
3. **Endpoint errors**: Verify endpoint URL and LocalStack configuration

### Configuration Issues

1. **Config not loading**: Check `config.yaml` syntax and file permissions
2. **Environment variables not working**: Ensure proper variable names and values
3. **Default values**: Check the `setDefaults()` function in `config/config.go`

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License. 

---

## 1. **Install AWS CLI**

**On Linux:**

```bash
# Download the AWS CLI v2 installer
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"

# Unzip the installer
unzip awscliv2.zip

# Run the install script (may require sudo)
sudo ./aws/install

# Verify installation
aws --version
```

If you don’t have `unzip`:
```bash
sudo apt-get update && sudo apt-get install unzip -y
```

---

## 2. **Install awslocal**

`awslocal` is a Python package. You need `pip` (Python’s package manager):

```bash
# If you don't have pip, install it:
sudo apt-get install python3-pip -y

# Install awslocal
pip3 install awscli-local

# Verify installation
awslocal --version
```

---

## 3. **Usage Examples**

- Use `awslocal` just like `aws`, but it automatically points to LocalStack:
  ```bash
  awslocal configservice list-config-rules
  ```

- Or use `aws` with the `--endpoint-url` flag:
  ```bash
  aws --endpoint-url=http://localhost:4566 configservice list-config-rules
  ```

---

Let me know if you want a script to automate these steps! 