#!/bin/bash

PROJECT_BASE_PATH=$(pwd)

SECRETS_SED_COMMANDS=""
CONFIGMAP_SED_COMMANDS=""

while IFS='=' read -r key value; do
  if [[ -z "$key" || -z "$value" ]]; then
    continue
  fi

  if grep -q "  $key:" k8/secrets.yaml; then
    ENCODED_VALUE=$(echo -n "$value" | base64)
    SECRETS_SED_COMMANDS+=" -e \"s|  $key: \".*\"|  $key: \\\"$ENCODED_VALUE\\\"|g\""
    SECRETS_SED_COMMANDS+=" -e \"s|  $key: .*|  $key: \\\"$ENCODED_VALUE\\\"|g\""
  fi

  if grep -q "  $key:" k8/config-maps.yaml; then
    CONFIGMAP_SED_COMMANDS+=" -e \"s|  $key: \".*\"|  $key: \\\"$value\\\"|g\""
    CONFIGMAP_SED_COMMANDS+=" -e \"s|  $key: .*|  $key: \\\"$value\\\"|g\""
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

kubectl create configmap prometheus-config --from-file=prometheus.yml=prometheus/prometheus.yml
kubectl create configmap keycloak-import --from-file=gocloak-realm.json=keycloak/import/gocloak-realm.json

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
