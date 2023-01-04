#!/usr/bin/env bash

set -xv
set -euo pipefail


kubectl create ns "$TELEMETRY_NS" || true


helm upgrade --install my-prom prometheus \
  --repo https://prometheus-community.github.io/helm-charts \
  --debug \
  --version 15.16.0 \
  --namespace "$TELEMETRY_NS"


# create datasource configuration

kubectl create secret generic grafana-datasource \
  --namespace "$TELEMETRY_NS" \
  --from-file=./grafana-datasource.yaml \
  --dry-run=client \
  -o yaml \
  | kubectl apply -f -

kubectl patch secret grafana-datasource \
  --namespace "$TELEMETRY_NS" \
  -p '{"metadata":{"labels":{"grafana_datasource": "1"}}}'

# create dashboard configuration

kubectl create secret generic grafana-dashboards \
  --namespace "$TELEMETRY_NS" \
  --from-file=./grafana-dashboards \
  --dry-run=client \
  -o yaml \
  | kubectl apply -f -

kubectl patch secret grafana-dashboards \
  --namespace "$TELEMETRY_NS" \
  -p '{"metadata":{"labels":{"grafana_dashboard": "1"}}}'

# set up grafana

helm upgrade --install my-grafana grafana \
  --repo https://grafana.github.io/helm-charts \
  --version 6.21.2 \
  --debug \
  --namespace "$TELEMETRY_NS" \
  -f grafana-values.yaml

# set up exporters

helm upgrade --install my-pg-prom-exporter prometheus-postgres-exporter \
  --repo https://prometheus-community.github.io/helm-charts \
  --version 2.5.0 \
  --debug \
  --namespace "$TELEMETRY_NS" \
  --set "config.datasource.host=my-pg-postgresql.${SCALING_NS}.svc.cluster.local" \
  --set config.datasource.password=postgres \
  --set config.datasource.database=postgres \
  --set config.autoDiscoverDatabases=true \
  --set config.logLevel=debug \
  --set config.logFormat=json \
  --set "config.includeDatabases={scaling}"

# set up jaeger

helm upgrade --install my-hunter jaeger \
  --repo https://jaegertracing.github.io/helm-charts \
  --version 0.62.1 \
  --wait \
  -n "$TELEMETRY_NS" \
  -f jaeger-values.yaml

# hack: add jaeger ingress b/c helm chart doesn't specify host
#   TODO there might be a way to do this through the helm chart
#   or could just use a port-forward
kubectl create -n "$TELEMETRY_NS" \
  -f jaeger-ingress.yaml \
  -o yaml --dry-run=client | kubectl apply -f -
