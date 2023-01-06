#!/usr/bin/env bash

set -xv
set -euo pipefail


SCALING_NS=${SCALING_NS:-"scaling"}


pushd ../cmd
    ./build.sh
popd

kubectl delete pod \
  -n "${SCALING_NS}" \
  -l app.kubernetes.io/component=webserver
