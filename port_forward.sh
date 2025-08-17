#!/bin/bash

REMOTE_IP="113.11.63.102"
SSH_PORT="42022"
SSH_USER="towhidul"

LOCAL_PORT1="8080"
LOCAL_PORT2="8081"

kubectl port-forward service/app 8080:8080 &
kubectl port-forward service/keycloak 8081:8081 &

ssh -p $SSH_PORT \
    -R $LOCAL_PORT1:localhost:$LOCAL_PORT1 \
    -R $LOCAL_PORT2:localhost:$LOCAL_PORT2 \
    -o "ServerAliveInterval 30" \
    -o "ServerAliveCountMax 3" \
    $SSH_USER@$REMOTE_IP

kill $(jobs -p)