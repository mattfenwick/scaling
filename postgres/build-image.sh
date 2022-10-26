#!/usr/bin/env bash

set -xv
set -euo pipefail

IMAGE=${IMAGE:-"localhost:5000/scaling-database:latest"}

docker build -t "${IMAGE}" .

docker push "${IMAGE}"
