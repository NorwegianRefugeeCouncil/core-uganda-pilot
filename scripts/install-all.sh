#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

echo Downloading go modules
(cd "${SCRIPT_DIR}/.." && go mod download)

echo Installing pwa dependencies
(cd "${SCRIPT_DIR}/../web" && yarn install)

