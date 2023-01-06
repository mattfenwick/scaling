#!/usr/bin/env bash

set -xv
set -euo pipefail


CREATE_CLUSTER=${CREATE_CLUSTER:-false}
BUILD_IMAGES=${BUILD_IMAGES:-true}
DEPLOY_CHARTS=${DEPLOY_CHARTS:-true}

export TELEMETRY_NS=${TELEMETRY_NS:-"telemetry"}
export KIND_NODE_IMAGE=${KIND_NODE_IMAGE:-"kindest/node:v1.24.6"}
export REGISTRY_PORT=${REGISTRY_PORT:-5000}
export REGISTRY_IMAGE=${REGISTRY_IMAGE:-"docker.io/library/registry:2"}
export REGISTRY_NAME=${REGISTRY_NAME:-"kind-registry"}
export KIND_CLUSTER_NAME=${KIND_CLUSTER_NAME:-"scaling"}
export NGINX_NS=${NGINX_NS:-"nginx-ingress"}
export PULL_IMAGES=${PULL_IMAGES:-true}

SCALING_NS=${SCALING_NS:-"scaling"}
LOADGEN_RELEASE=${LOADGEN_RELEASE:-"my-loadgen"}
SERVER_RELEASE=${SERVER_RELEASE:-"my-scaling"}
POSTGRES_RELEASE_NAME=${POSTGRES_RELEASE_NAME:-"my-pg"}


if [[ $CREATE_CLUSTER == true ]]; then
    pushd kind
        ./deploy.sh
    popd

    ./load-images.sh

    pushd telemetry
        ./deploy.sh
    popd
fi


if [[ $BUILD_IMAGES == true ]]; then
    pushd ../cmd
        ./build.sh
    popd
fi


if [[ $DEPLOY_CHARTS == true ]]; then
    kubectl create ns "${SCALING_NS}" || true

    # postgres
    helm upgrade "${POSTGRES_RELEASE_NAME}" postgresql \
        --install \
        --repo https://charts.bitnami.com/bitnami \
        --version 11.6.2 \
        --namespace "$SCALING_NS" \
        --timeout 5m0s \
        --wait \
        -f postgres-values.yaml

    # server
    helm upgrade --install "$SERVER_RELEASE" \
        --namespace "${SCALING_NS}" \
        ../charts/server \
        --values server-values.yaml \
        --wait

    # loadgen
    helm upgrade --install "${LOADGEN_RELEASE}" \
        ../charts/server \
        --namespace "${SCALING_NS}" \
        --values loadgen-values.yaml
fi


# grafana:
#   http://scaling.local/d/ZZ1598S4k/webserver?orgId=1
# jaeger:
#   http://jaeger.local/search