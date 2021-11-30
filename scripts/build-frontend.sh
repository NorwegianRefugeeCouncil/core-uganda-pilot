#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

set -e

echo ">> Building web/app/frontend"
rm -rf "${ROOT_DIR}/web/app/frontend/node_modules"
(cd "${ROOT_DIR}/web/app/frontend" && yarn install)

