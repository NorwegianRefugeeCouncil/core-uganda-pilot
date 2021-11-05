#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

set -e

# cleanup on exit
trap 'jobs -p | xargs -r kill' EXIT

(cd "${SCRIPT_DIR}/.." && go run . seed -u "http://localhost:9000")
