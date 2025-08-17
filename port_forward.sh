#!/bin/bash
kubectl port-forward service/ingress-nginx 4747:80 -n ingress-nginx