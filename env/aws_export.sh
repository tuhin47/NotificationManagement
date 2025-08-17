echo "🔧 Setting up environment variables..."
export AWS_ENDPOINT=http://localhost:4566
export AWS_REGION=us-east-1
export AWS_ACCESS_KEY_ID=test
export AWS_SECRET_ACCESS_KEY=test

echo "📋 Environment variables set:"
echo "  AWS_ENDPOINT=$AWS_ENDPOINT"
echo "  AWS_REGION=$AWS_REGION"

export CONFIG_SSM_PARAM="/myapp/config"
export CONFIG_FROM_SSM=true

aws --endpoint-url=http://localhost:4566 ssm put-parameter \
  --name "$CONFIG_SSM_PARAM" \
  --region "$AWS_REGION" \
  --type "String" \
  --value "$(cat /tmp/app-config.json)" \
  --overwrite

