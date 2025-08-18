#!/usr/bin/env bash

HOST_DIR=$(pwd)
VM_DIR=$(pwd)
minikube start \
  --mount \
  --mount-string="${HOST_DIR}:${VM_DIR}" \
  --driver=docker \
  --cpus=3 \
  --memory=10g \
  --force

