#!/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

docker-compose -f deployments/webapp.docker-compose.yaml down --remove-orphans

(cd $SCRIPT_DIR && ./reset_local_config.sh)
