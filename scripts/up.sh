#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

set -e

"${SCRIPT_DIR}/init_secrets.sh"

COMPOSE_PROJECT_NAME=core docker-compose -f deployments/webapp.docker-compose.yaml up --build -d
