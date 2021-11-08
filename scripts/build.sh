#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

set -e

echo ">> Installing web/app/client npm dependencies"
(cd "${ROOT_DIR}/web/app/client" && yarn install)

echo ">> Installing web/app/designSystem npm dependencies"
(cd "${ROOT_DIR}/web/app/designSystem" && yarn install)

echo ">> Transpiling web/app/client"
tsc --build "${ROOT_DIR}/web/app/client/tsconfig.json"

echo ">> Transpiling web/app/designSystem"
tsc --build "${ROOT_DIR}/web/app/designSystem/tsconfig.json"

echo ">> Building web/app/frontend"
rm -rf "${ROOT_DIR}/web/app/frontend/node_modules"
(cd "${ROOT_DIR}/web/app/frontend" && yarn install)

echo ">> Building web/pwa"
rm -rf "${ROOT_DIR}/web/pwa/node_modules"
(cd "${ROOT_DIR}/web/pwa" && yarn install)

echo ">> Building web/admin"
rm -rf "${ROOT_DIR}/web/admin/node_modules"
(cd "${ROOT_DIR}/web/admin" && yarn install)

echo ">> Building core server"
go build -o "${ROOT_DIR}/tmp/main" "${ROOT_DIR}/cmd"

