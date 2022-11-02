#!/usr/bin/env bash

set -xv
set -euo pipefail

IMAGE=${IMAGE:-localhost:5000/scaling:latest}


CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go


docker build -t $IMAGE .

docker push $IMAGE
