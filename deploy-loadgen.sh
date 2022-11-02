#!/usr/bin/env bash

set -xv
set -euo pipefail


SCALING_NAMESPACE=${SCALING_NAMESPACE:-"scaling"}
LOADGEN_RELEASE=${LOADGEN_RELEASE:-"my-loadgen"}
LOADGEN_IMAGE="localhost:5000/scaling:latest"


# TODO: build, push image
# TODO: provide helm override to point to image
pushd cmd
    IMAGE=$LOADGEN_IMAGE ./build.sh
popd


helm upgrade --install "${LOADGEN_RELEASE}" \
  ./charts/loadgen \
  --namespace "${SCALING_NAMESPACE}" \
  --values ./dev/loadgen/values.yaml \
  --set loadgen.binary="/main" \
  --set loadgen.image="$LOADGEN_IMAGE"
