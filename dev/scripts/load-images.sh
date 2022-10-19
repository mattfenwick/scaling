#!/usr/bin/env bash

set -xv
set -euo pipefail

PULL_IMAGE=${PULL_IMAGE:-false}

declare -a IMAGES=(
  "grafana/grafana:8.3.4"
  "quay.io/kiwigrid/k8s-sidecar:1.19.2"

  "registry.k8s.io/kube-state-metrics/kube-state-metrics:v2.5.0"
  "quay.io/prometheus/alertmanager:v0.24.0"
  "jimmidyson/configmap-reload:v0.5.0"
  "quay.io/prometheus/node-exporter:v1.3.1"
  "prom/pushgateway:v1.4.3"
  "quay.io/prometheus/prometheus:v2.39.1"

  "docker.io/bitnami/postgresql:14.4.0-debian-11-r4"

  "grafana/loki:2.6.1"

  "docker.io/grafana/promtail:2.6.1"

  "quay.io/prometheuscommunity/postgres-exporter:v0.10.0"

  "k8s.gcr.io/ingress-nginx/controller:v1.1.1@sha256:0bc88eb15f9e7f84e8e56c14fa5735aaa488b840983f87bd79b1054190e660de"

  "jaegertracing/all-in-one:1.37.0")


for image in "${IMAGES[@]}"
do
  if [[ $PULL_IMAGE == true ]]; then
    docker pull "$image"
  fi

  kind load docker-image "$image"
done
