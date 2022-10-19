#!/usr/bin/env bash

set -xv
set -euo pipefail

SCALING_NS=${SCALING_NS:-"scaling"}

kubectl create ns "${SCALING_NS}" || true

kubectl create configmap webserver-config \
  --namespace "${SCALING_NS}" \
  --from-file=config.json=../cmd/webserver.json \
  -o yaml --dry-run=client | kubectl apply -f -

helm upgrade --install my-scaling \
  --namespace "${SCALING_NS}" \
  ../charts/scaling \
  --set configMap=webserver-config
