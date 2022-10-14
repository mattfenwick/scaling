#!/usr/bin/env bash

set -xv
set -euo pipefail

declare -a IMAGES=(
  "grafana/grafana:9.1.7"
  "quay.io/kiwigrid/k8s-sidecar:1.19.2"

  "registry.k8s.io/kube-state-metrics/kube-state-metrics:v2.5.0"
  "quay.io/prometheus/alertmanager:v0.24.0"
  "jimmidyson/configmap-reload:v0.5.0"
  "quay.io/prometheus/node-exporter:v1.3.1"
  "prom/pushgateway:v1.4.3"
  "quay.io/prometheus/prometheus:v2.39.1"

  "quay.io/prometheuscommunity/postgres-exporter:v0.10.0"

  "jaegertracing/all-in-one:1.37.0")


for image in "${IMAGES[@]}"
do
  docker pull "$image"

  kind load docker-image "$image"
done