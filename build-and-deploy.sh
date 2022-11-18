#!/usr/bin/env bash

set -xv
set -euo pipefail

# hack: add jaeger ingress b/c helm chart doesn't specify host
#   TODO there might be a way to do this through the helm chart
#   or could just use a port-forward
kubectl create -n telemetry \
  -f jaeger-ingress.yaml \
  -o yaml --dry-run=client | kubectl apply -f -
  

pushd cmd
  ./build.sh
popd

pushd postgres
  ./build-image.sh
popd

# TODO deploy postgres chart ????

helm upgrade --install -n scaling \
  my-scaling \
  charts/server \
  --set postgres.image=localhost:5000/scaling-database:latest \
  --set webserver.binary=/main \
  --set webserver.image=localhost:5000/scaling:latest \
  --set "jaegerUrl=http://my-hunter-jaeger-collector.telemetry.svc.cluster.local:14268/api/traces" \
  --set "logLevel=debug"

helm upgrade --install -n scaling \
  my-loadgen \
  charts/loadgen \
  --set loadgen.binary=/main \
  --set loadgen.image=localhost:5000/scaling:latest \
  --set loadgen.job.enabled=true \
  --set loadgen.deployment.enabled=false \
  --set "jaegerUrl=http://my-hunter-jaeger-collector.telemetry.svc.cluster.local:14268/api/traces" \
  --set "logLevel=debug" \
  --set "webserver.host=my-scaling-server-webserver"

# grafana:
#   http://scaling.local/d/ZZ1598S4k/webserver?orgId=1
# jaeger:
#   http://jaeger.local/search
