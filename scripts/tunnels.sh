#!/bin/bash
set -e
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
(cd "${SCRIPT_DIR}/.." && go run "." dev tunnels)
