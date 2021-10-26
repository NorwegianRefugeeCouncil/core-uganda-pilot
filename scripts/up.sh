#!/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

set -e

docker-compose -f deployments/webapp.docker-compose.yaml up --build -d

echo waiting for database
(cd "${SCRIPT_DIR}" && ./waitforit.sh -h localhost -p 5433)

echo database ready. migrating
(cd "${SCRIPT_DIR}" && ./migrate.sh -h localhost -p 5433)
