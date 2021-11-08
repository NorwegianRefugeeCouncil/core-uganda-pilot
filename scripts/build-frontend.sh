#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

set -e

echo ">> Installing web/app/client npm dependencies"
(cd "${ROOT_DIR}/web/app/client" && yarn install)

echo ">> Transpiling web/app/client"
tsc --build "${ROOT_DIR}/web/app/client/tsconfig.json"

echo ">> Building web/app/frontend"
rm -rf "${ROOT_DIR}/web/app/frontend/node_modules"
(cd "${ROOT_DIR}/web/app/frontend" && yarn install)

echo ">> Building core server"
go build -o "${ROOT_DIR}/tmp/main" "${ROOT_DIR}/cmd"

