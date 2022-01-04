#!/bin/bash
set -e
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

cd "${SCRIPT_DIR}/.."

TUNNEL_PID=
PROXY_PID=
SERVE_PID=
DOCKER_UP=0

function cleanup {
  echo "Cleaning up..."

  echo "Killing tunnel..."
  pkill -f $(cat creds/pet)

  echo "Killing proxy..."
  pkill -f "deployments/envoy-local.yaml"

  if [ -n "${SERVE_PID}" ]; then
    if ps -p $SERVE_PID > /dev/null; then
      echo "Killing backend..."
      kill "${SERVE_PID}"
    fi
  fi

  if [ "${DOCKER_UP}" -eq 1 ]; then
    echo "Stopping docker..."
    make down
  fi
}

trap cleanup EXIT

echo "Starting tunnels..."
make tunnels > /dev/null 2>&1 &
TUNNEL_PID=$!
echo "Tunnel running on PID ${TUNNEL_PID}"

echo "Starting proxy..."
make proxy-local > /dev/null 2>&1 &
PROXY_PID=$!
echo "Proxy running on PID ${PROXY_PID}"

echo "Bringing up docker..."
make up
DOCKER_UP=1

echo "Waiting 40s for docker images to come up..."
sleep 40

echo "Bootstrapping backend..."
make bootstrap

echo "Serving backend..."
make serve &
SERVE_PID=$!
echo "Backend running on PID ${SERVE_PID}"

wait