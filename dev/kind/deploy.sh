#!/usr/bin/env bash

set -xv
set -euo pipefail


# cluster
kind create cluster \
  --name "$KIND_CLUSTER_NAME" \
  --config=./kind-config.yaml \
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


# local image registry
running="$(docker inspect -f '{{.State.Running}}' "${REGISTRY_NAME}" 2>/dev/null || true)"
if [ "${running}" != 'true' ]; then
  docker run \
    --network "kind" \
    -d --restart=always -p "127.0.0.1:${REGISTRY_PORT}:5000" --name "${REGISTRY_NAME}" \
    "$REGISTRY_IMAGE"
fi


kubectl get nodes
kubectl wait --for=condition="Ready" nodes --all --timeout="15m"


# metrics server
kubectl apply -f ./metrics-server.yaml


# ingress controller
kubectl create ns "$NGINX_NS" || true

helm install --namespace "$NGINX_NS" my-nginx ingress-nginx \
  --repo https://kubernetes.github.io/ingress-nginx \
  --debug \
  --wait \
  --version 4.0.17 \
  -f ./nginx-values.yaml
