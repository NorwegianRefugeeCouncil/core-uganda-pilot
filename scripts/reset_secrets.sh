#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

set -e

CORE_CONFIG_FILE="${SCRIPT_DIR}/../configs/config.custom.yaml"
HYDRA_CONFIG_FILE="${SCRIPT_DIR}/../deployments/hydra.custom.yaml"
POSTGRES_ENV_FILE="${SCRIPT_DIR}/../deployments/postgres.env"

rm "${CORE_CONFIG_FILE}" || echo "${CORE_CONFIG_FILE}" does not exist
rm "${HYDRA_CONFIG_FILE}" || echo "${HYDRA_CONFIG_FILE}" does not exist
rm "${POSTGRES_ENV_FILE}" || echo "${POSTGRES_ENV_FILE}" does not exist
