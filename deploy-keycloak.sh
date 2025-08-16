#!/bin/bash

PROJECT_BASE_PATH=$(pwd)

for file in k8/*.yaml; do
  if grep -q "__PROJECT_BASE_PATH__" "$file"; then
    echo "Processing $file..."
    sed "s|__PROJECT_BASE_PATH__|$PROJECT_BASE_PATH|g" "$file" | kubectl apply -f -
  else
    echo "Applying $file without path replacement..."
    kubectl apply -f "$file"
  fi
done

echo "All Kubernetes manifests in 'k8' directory processed with dynamic path: $PROJECT_BASE_PATH"
