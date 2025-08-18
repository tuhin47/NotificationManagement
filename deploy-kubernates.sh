#!/bin/bash

PROJECT_BASE_PATH=$(pwd)

# Check for --del argument
if [[ "$1" == "--del" ]]; then
  echo "Deleting all Kubernetes resources..."
  kubectl delete deployment,service,configmap,secret,ingress,statefulset --all -n default
  echo "Kubernetes resources deleted."
  exit 0
fi

SECRETS_SED_COMMANDS=""
CONFIGMAP_SED_COMMANDS=""

while IFS='=' read -r key value; do
  if [[ -z "$key" ]]; then
    continue
  fi

  if [[ "$key" == "GEMINI_KEY" || "$key" == "TELEGRAM_TOKEN" ]]; then
    if [[ -n "${!key}" ]]; then
      value="${!key}"
      echo "Using environment variable for $key"
    else
      echo "Warning: $key not found in environment variables. Using value from .env file."
    fi
  fi
  
  if [[ -z "$value" ]]; then
    continue
  fi

  if grep -q "  $key:" k8/secrets.yaml; then
    ENCODED_VALUE=$(echo -n "$value" | base64)
    SECRETS_SED_COMMANDS+=" -e \"s|  $key: \".*\"|  $key: \\\"$ENCODED_VALUE\\\"|g\""
    SECRETS_SED_COMMANDS+=" -e \"s|  $key: .*|  $key: \\\"$ENCODED_VALUE\\\"|g\""
  fi

  if grep -q "  $key:" k8/config-maps.yaml; then
    CONFIGMAP_SED_COMMANDS+=" -e \"s|^[[:space:]]*$key:.*$|  $key: \\\"$value\\\"|g\""
  fi
done < .env

if [[ -n "$SECRETS_SED_COMMANDS" ]]; then
  echo "Updating k8/secrets.yaml with values from .env..."
  eval "cat k8/secrets.yaml | sed $SECRETS_SED_COMMANDS" | kubectl apply -f -
else
  echo "No matching secrets found in .env to update k8/secrets.yaml."
  kubectl apply -f k8/secrets.yaml
fi

if [[ -n "$CONFIGMAP_SED_COMMANDS" ]]; then
  echo "Updating k8/config-maps.yaml with values from .env..."
  eval "cat k8/config-maps.yaml | sed $CONFIGMAP_SED_COMMANDS" | kubectl apply -f -
else
  echo "No matching config map values found in .env to update k8/config-maps.yaml."
  kubectl apply -f k8/config-maps.yaml
fi

for file in k8/*.yaml; do
  if [ "$file" == "k8/secrets.yaml" ] || [ "$file" == "k8/config-maps.yaml" ]; then
    echo "Skipping $file as it's already processed."
    continue
  fi

  if grep -q "__PROJECT_BASE_PATH__" "$file"; then
    echo "Processing $file..."
    sed "s|__PROJECT_BASE_PATH__|$PROJECT_BASE_PATH|g" "$file" | kubectl apply -f -
  else
    echo "Applying $file without path replacement..."
    kubectl apply -f "$file"
  fi
done
echo "All Kubernetes manifests in 'k8' directory processed with dynamic path: $PROJECT_BASE_PATH"

echo "Waiting for config-server statefulset to be ready..."
kubectl rollout status --watch --timeout=300s statefulset/config-server

CONFIG_SERVER_POD=$(kubectl get pods -l app=config-server -o jsonpath='{.items[0].metadata.name}')

if [ -z "$CONFIG_SERVER_POD" ]; then
  echo "Error: config-server pod not found."
  exit 1
fi

kubectl cp env/app-config-prod.json "$CONFIG_SERVER_POD":/tmp/app-config.json
kubectl cp env/aws_export.sh "$CONFIG_SERVER_POD":/tmp/aws_export.sh
kubectl exec "$CONFIG_SERVER_POD" -- bash /tmp/aws_export.sh;
