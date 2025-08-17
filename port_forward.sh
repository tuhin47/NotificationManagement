#!/bin/bash
kubectl port-forward service/app 4747:80 -n ingress-nginx