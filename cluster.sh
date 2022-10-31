#!/usr/bin/env bash

set -xv
set -euo pipefail

KIND_NODE_IMAGE="kindest/node:v1.24.6"
REGISTRY_PORT=5000
REGISTRY_IMAGE="docker.io/library/registry:2"
REGISTRY_NAME="kind-registry"
KIND_CLUSTER_NAME="scaling"


kind create cluster \
  --name "$KIND_CLUSTER_NAME" \
  --config=./dev/scripts/kind-config.yaml \
  --image "$KIND_NODE_IMAGE"

# Document the local registry
# https://github.com/kubernetes/enhancements/tree/master/keps/sig-cluster-lifecycle/generic/1755-communicating-a-local-registry
cat <<-EOF  | kubectl apply -f -
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: local-registry-hosting
      namespace: kube-public
    data:
      localRegistryHosting.v1: |
        host: "localhost:${REGISTRY_PORT}"
        help: "https://kind.sigs.k8s.io/docs/user/local-registry/"
EOF

running="$(docker inspect -f '{{.State.Running}}' "${REGISTRY_NAME}" 2>/dev/null || true)"
if [ "${running}" != 'true' ]; then
  docker run \
    --network "kind" \
    -d --restart=always -p "127.0.0.1:${REGISTRY_PORT}:5000" --name "${REGISTRY_NAME}" \
    "$REGISTRY_IMAGE"
fi
