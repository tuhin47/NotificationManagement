#!/bin/bash

# Upload env/app-config.json to LocalStack SSM Parameter Store
aws --endpoint-url=http://localhost:4566 ssm put-parameter \
  --name "/myapp/config" \
  --type "String" \
  --value "$(cat env/app-config.json)" \
  --overwrite