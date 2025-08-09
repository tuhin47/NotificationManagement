# AI Notification Management

<!-- TOC -->
* [AI Notification Management](#ai-notification-management)
  * [Feature Overview](#feature-overview)
    * [Functional Features](#functional-features)
    * [Non-Functional Features](#non-functional-features)
  * [Quick Start](#quick-start)
    * [Prerequisites](#prerequisites)
    * [1. Setup Environment](#1-setup-environment)
    * [2. Build and Run the Application](#2-build-and-run-the-application)
  * [Services](#services)
  * [Authentication](#authentication)
  * [Future Plans](#future-plans)
<!-- TOC -->

## Feature Overview

This project is a notification management system that leverages AI for processing notifications. It supports various formats such as JSON, HTML, XML, CSV, and text. The system includes features like scheduled notifications, Single Sign-On (SSO) with Keycloak, and integration with AWS services using LocalStack for local development.

### Functional Features

- AI based notification management (supports json, html, xml, csv, and text formats)
- Scheduled notifications
- Single Sign-On (SSO) with Keycloak

### Non-Functional Features

- AWS Config service integration
    - LocalStack for local AWS service emulation
    - Environment Variables override from Environment Variables
- Gemini AI integration for processing notifications
- Telegram bot integration for notifications
- Configured :
    - Logging and error handling (using Zap logger)
    - Dependency injection (using uber-fx)
- Generic worker for processing scheduled tasks
- GORM for database operations
- Postman collection for API testing

## Quick Start

### Prerequisites

- Go 1.24 or higher
- Docker and Docker Compose (for LocalStack)
- AWS CLI (optional, for production)

### 1. Setup Environment

```bash
# Start LocalStack with AWS Config service using docker
sh scripts/setup-dev.sh

# Verify LocalStack is running
curl http://localhost:4566/_localstack/health
```

### 2. Build and Run the Application

Please create the GEMINI_KEY from [here](https://aistudio.google.com/app/apikey) and TELEGRAM_TOKEN from [here](https://core.telegram.org/bots/tutorial#obtain-your-bot-token)

```bash
# Build the application
go build -o notification-management
# Run the application
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export GEMINI_KEY=
go run main.go serve

# Run Worker
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test
export GEMINI_KEY=
export TELEGRAM_ENABLED=true
export TELEGRAM_TOKEN=
go run main.go worker
```

## Services

| Service     | Description            | Port  |
|-------------|------------------------|-------|
| Application | Main app server        | 8080  |
| PostgreSQL  | Database in Docker     | 54322 |
| LocalStack  | AWS services emulation | 4566  |
| Keycloak    | SSO                    | 8081  |
| Mailcatcher | Email testing          | 1080  |

## Authentication

To use the application, you need to follow these steps:

- Import postman collection from `postman/notification-management.postman_collection.json`
- Expand Notification Management folder > Authentication section
- Click on `Get New Access Token`
  Username: `nms`
  Password: `admin`
- Click on `Use Token`

Now you can use the APIs in the collection.

## Future Plans

- Authentication in requests
- AI based todo setup
- Report generation based on periodic results
- Additional notification formats (e.g., PDF, Excel)
- Add prometheus and grafana for monitoring
- Add OPENAI integration for AI processing
