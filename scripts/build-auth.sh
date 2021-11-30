#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

set -e

echo ">> Building web/auth"
rm -rf "${ROOT_DIR}/web/auth/node_modules"
(cd "${ROOT_DIR}/web/auth" && yarn install && tsc)
