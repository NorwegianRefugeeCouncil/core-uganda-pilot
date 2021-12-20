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
  if [ -n "${TUNNEL_PID}" ]; then
    echo "Killing tunnel..."
    kill "${TUNNEL_PID}"
  fi
  if [ -n "${PROXY_PID}" ]; then
    echo "Killing proxy..."
    kill "${PROXY_PID}"
  fi
  if [ -n "${SERVE_PID}" ]; then
    echo "Killing backend..."
    kill "${SERVE_PID}"
  fi
  if [ "${DOCKER_UP}" -eq 1 ]; then
    echo "Stopping docker..."
    make down
  fi
}

trap cleanup SIGINT EXIT ERR HUP INT QUIT TERM

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

echo "Waiting 30s for docker images to come up..."
sleep 30

echo "Bootstrapping backend..."
make bootstrap

echo "Serving backend..."
make serve &
SERVE_PID=$!
echo "Backend running on PID ${SERVE_PID}"