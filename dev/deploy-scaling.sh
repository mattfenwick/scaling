#!/usr/bin/env bash

set -xv
set -euo pipefail

kubectl create configmap webserver-config \
  --from-file=config.json=../cmd/webserver.json \
  -o yaml --dry-run=client | kubectl apply -f -

helm upgrade --install my-scaling ../charts/scaling \
  --set configMap=webserver-config
