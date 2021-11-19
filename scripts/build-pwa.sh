#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

set -e

echo ">> Installing web/app/client npm dependencies"
(cd "${ROOT_DIR}/web/app/client" && yarn install && tsc)

echo ">> Installing web/auth npm dependencies"
rm -rf "${ROOT_DIR}/web/auth/node_modules"
(cd "${ROOT_DIR}/web/auth" && yarn install && tsc)

echo ">> Building web/pwa"
rm -rf "${ROOT_DIR}/web/pwa/node_modules"
(cd "${ROOT_DIR}/web/pwa" && yarn install)

echo ">> Building core server"
go build -o "${ROOT_DIR}/tmp/main" "${ROOT_DIR}/cmd"

