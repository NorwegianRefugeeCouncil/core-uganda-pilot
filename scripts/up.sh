#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

set -e

"${SCRIPT_DIR}/init_secrets.sh"

docker-compose -f deployments/webapp.docker-compose.yaml up --build -d

echo "Waiting for database"
(cd "${SCRIPT_DIR}" && ./waitforit.sh -h localhost -p 5433)

sleep 15

echo "Database ready. migrating"
(cd "${SCRIPT_DIR}" && ./migrate.sh -h localhost -p 5433)

echo "Waiting for hydra migrations"

sleep 10

echo "Bootstrapping initial data for dev environment"
"${SCRIPT_DIR}/bootstrap.sh"
