#!/usr/bin/env bash
set -euo pipefail
NAME="${1:?image name is required}"
PLATFORM="${2:-linux/amd64}"
IMAGE="benzhi/${NAME}:latest"
echo "RUN go build ./..."
docker build --platform "${PLATFORM}" -f benzhi.Dockerfile -t "${IMAGE}" .
echo "DONE"
