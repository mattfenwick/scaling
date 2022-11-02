#!/usr/bin/env bash

set -xv
set -euo pipefail

CREATE_CLUSTER=${CREATE_CLUSTER:-false}

if [[ $CREATE_CLUSTER == true ]]; then
  ./cluster.sh
fi

TELEMETRY_NS=telemetry
NGINX_NS="nginx"


kubectl create ns $NGINX_NS || true
kubectl create ns scaling || true
kubectl create ns $TELEMETRY_NS || true

kubectl create secret generic grafana-datasource \
  --namespace "$TELEMETRY_NS" \
  --from-file=grafana-datasource.yaml=./dev/scripts/grafana-datasource.yaml \
  --dry-run=client \
  -o yaml \
  | kubectl apply -f -

kubectl patch secret grafana-datasource \
  --namespace "$TELEMETRY_NS" \
  -p '{"metadata":{"labels":{"grafana_datasource": "1"}}}'

# create dashboard configuration

kubectl create secret generic grafana-dashboards \
  --namespace "$TELEMETRY_NS" \
  --from-file=./dev/scripts/grafana-dashboards \
  --dry-run=client \
  -o yaml \
  | kubectl apply -f -

kubectl patch secret grafana-dashboards \
  --namespace "$TELEMETRY_NS" \
  -p '{"metadata":{"labels":{"grafana_dashboard": "1"}}}'

skaffold run

#skaffold dev -m server
#skaffold debug -m server

