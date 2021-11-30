#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

set -e

echo ">> Building web/pwa"
rm -rf "${ROOT_DIR}/web/pwa/node_modules"
(cd "${ROOT_DIR}/web/pwa" && yarn install)

