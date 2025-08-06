# Notification Management System


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

### 1. Setup Environment

```bash
# Start LocalStack with AWS Config service using docker
sh scripts/setup-dev.sh

# Verify LocalStack is running
curl http://localhost:4566/_localstack/health
```

### 2. Build and Run the Application

```bash
# Build the application
go build -o notification-management
# Build the application
go run main.go serve
```