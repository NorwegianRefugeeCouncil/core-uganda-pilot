#!/bin/bash
set -e

COMPOSE_PROJECT_NAME=core docker-compose -f deployments/webapp.docker-compose.yaml logs "$@"
