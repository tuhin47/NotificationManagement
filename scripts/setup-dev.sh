#!/bin/bash

# Development Setup Script for Notification Management System
# This script sets up the development environment with LocalStack

set -e

echo "🚀 Setting up Notification Management System development environment..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

echo "📦 Starting LocalStack..."
docker compose up -d config-server postgres keycloak mailcatcher redis

echo "⏳ Waiting for config-server to be ready..."
sleep 10

# Check if LocalStack is running
if curl -s http://localhost:4566/_localstack/health > /dev/null; then
    echo "✅ LocalStack is running successfully!"
else
    echo "❌ LocalStack failed to start. Check logs with: docker-compose logs localstack"
    exit 1
fi

echo "🔧 Setting up environment variables..."
export AWS_ENDPOINT=http://localhost:4566
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test

echo "📋 Environment variables set:"
echo "  AWS_ENDPOINT=$AWS_ENDPOINT"
echo "  AWS_REGION=$AWS_REGION"

# Check if config is already pushed to SSM; if not, push using env/upload-config.sh

echo "🔍 Checking if configuration is already pushed to SSM..."

CONFIG_SSM_PARAM="/myapp/config"
CONFIG_FROM_SSM=true

# Try to get the parameter from LocalStack SSM
PARAM_EXISTS=$(aws --endpoint-url=$AWS_ENDPOINT --region $AWS_REGION ssm get-parameter --name "$CONFIG_SSM_PARAM" > /dev/null 2>&1 && echo "yes" || echo "no")

if [ "$PARAM_EXISTS" = "yes" ]; then
    echo "✅ Configuration already exists in SSM ($CONFIG_SSM_PARAM)."
else
    echo "⚠️  Configuration not found in SSM. Uploading from env/app-config.json..."
    if [ -f env/upload-config.sh ]; then
        bash env/upload-config.sh
        if [ $? -eq 0 ]; then
            echo "✅ Configuration uploaded to SSM successfully."
        else
            echo "❌ Failed to upload configuration to SSM."
            exit 1
        fi
    else
        echo "❌ env/upload-config.sh not found!"
        exit 1
    fi
fi


echo "🧪 Testing AWS Config service connection..."
if go run main.go config test-connection; then
    echo "✅ AWS Config service connection successful!"
else
    echo "❌ AWS Config service connection failed. Check LocalStack logs."
    exit 1
fi
