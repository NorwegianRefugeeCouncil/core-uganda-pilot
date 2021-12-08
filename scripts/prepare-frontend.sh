#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

set -e

(cd "${ROOT_DIR}/frontend" && yarn install)
(cd "${ROOT_DIR}/frontend" && yarn build:packages)
