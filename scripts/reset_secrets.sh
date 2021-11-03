#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

set -e

CORE_CONFIG_FILE="${SCRIPT_DIR}/../configs/config.custom.yaml"
HYDRA_CONFIG_FILE="${SCRIPT_DIR}/../deployments/hydra.custom.yaml"
POSTGRES_ENV_FILE="${SCRIPT_DIR}/../deployments/postgres.env"
OIDC_CONFIG_FILE="${SCRIPT_DIR}/../deployments/oidc.config.json"
OIDC_USERS_FILE="${SCRIPT_DIR}/../deployments/oidc.users.json"
REDIS_ENV_FILE="${SCRIPT_DIR}/../deployments/redis.env"
PWA_ENV_FILR="${SCRIPT_DIR}/../web/pwa/.env"
ADMIN_ENV_FILR="${SCRIPT_DIR}/../web/admin/.env"

rm "${CORE_CONFIG_FILE}" || echo "${CORE_CONFIG_FILE}" does not exist
rm "${HYDRA_CONFIG_FILE}" || echo "${HYDRA_CONFIG_FILE}" does not exist
rm "${POSTGRES_ENV_FILE}" || echo "${POSTGRES_ENV_FILE}" does not exist
rm "${OIDC_CONFIG_FILE}" || echo "${OIDC_CONFIG_FILE}" does not exist
rm "${OIDC_USERS_FILE}" || echo "${OIDC_USERS_FILE}" does not exist
rm "${REDIS_ENV_FILE}" || echo "${REDIS_ENV_FILE}" does not exist
rm "${PWA_ENV_FILR}" || echo "${PWA_ENV_FILR}" does not exist
rm "${ADMIN_ENV_FILR}" || echo "${ADMIN_ENV_FILR}" does not exist
