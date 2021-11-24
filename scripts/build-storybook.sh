#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

set -e

echo ">> Installing web/designsystem npm dependencies"
rm -rf "${ROOT_DIR}/web/designsystem/node_modules"
(cd "${ROOT_DIR}/web/designsystem" && yarn install && tsc)

echo ">> Building web/storyBook"
rm -rf "${ROOT_DIR}/web/storyBook/node_modules"
(cd "${ROOT_DIR}/web/storyBook" && yarn install)
