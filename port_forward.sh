#!/bin/bash
kubectl port-forward service/app 8080:8080& 
kubectl port-forward service/keycloak 8081:8081 &
kubectl port-forward service/mailcatcher 1080:1080 