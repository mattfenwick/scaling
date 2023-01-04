#!/usr/bin/env bash

set -xv
set -euo pipefail


CREATE_CLUSTER=${CREATE_CLUSTER:-false}
export TELEMETRY_NS=${TELEMETRY_NS:-"telemetry"}
export APP_NS=${APP_NS:-"scaling"}
export KIND_NODE_IMAGE=${KIND_NODE_IMAGE:-"kindest/node:v1.24.6"}
export REGISTRY_PORT=${REGISTRY_PORT:-5000}
export REGISTRY_IMAGE=${REGISTRY_IMAGE:-"docker.io/library/registry:2"}
export REGISTRY_NAME=${REGISTRY_NAME:-"kind-registry"}
export KIND_CLUSTER_NAME=${KIND_CLUSTER_NAME:-"scaling"}
export NGINX_NS=${NGINX_NS:-"nginx-ingress"}
export PULL_IMAGE=${PULL_IMAGE:-true}

export SCALING_NS=${SCALING_NS:-"scaling"}
export LOADGEN_RELEASE=${LOADGEN_RELEASE:-"my-loadgen"}
export SERVER_RELEASE=${SERVER_RELEASE:-"my-scaling"}


if [[ $CREATE_CLUSTER == true ]]; then
    pushd kind
        ./deploy.sh
    popd

    ./load-images.sh

    pushd telemetry
        ./deploy.sh
    popd
fi


if [[ $BUILD == true ]]; then
    pushd cmd
    ./build.sh
    popd

    pushd postgres
    ./build-image.sh
    popd
fi


if [[ $DEPLOY == true ]]; then
    kubectl create ns "${SCALING_NS}" || true

    # TODO deploy postgres chart ????

    # server
    helm upgrade --install "$SERVER_RELEASE" \
        --namespace "${SCALING_NS}" \
        ../charts/server \
        --values server-values.yaml

    # loadgen
    helm upgrade --install "${LOADGEN_RELEASE}" \
        ../charts/loadgen \
        --namespace "${SCALING_NS}" \
        --values loadgen-values.yaml
fi


# grafana:
#   http://scaling.local/d/ZZ1598S4k/webserver?orgId=1
# jaeger:
#   http://jaeger.local/search