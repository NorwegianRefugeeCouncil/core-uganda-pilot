#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

CORE_CONFIG_FILE="${SCRIPT_DIR}/../configs/config.custom.yaml"
HYDRA_CONFIG_FILE="${SCRIPT_DIR}/../deployments/hydra.custom.yaml"
POSTGRES_ENV_FILE="${SCRIPT_DIR}/../deployments/postgres.env"
REDIS_ENV_FILE="${SCRIPT_DIR}/../deployments/redis.env"

if [ ! -f "${REDIS_ENV_FILE}"  ] || [ ! -f "${CORE_CONFIG_FILE}" ] || [ ! -f "${POSTGRES_ENV_FILE}" ] || [ ! -f "${HYDRA_CONFIG_FILE}" ]; then

  POSTGRES_USER=postgres
  POSTGRES_PASSWORD=$(openssl rand -hex 16)
  HYDRA_DB=hydra
  HYDRA_USERNAME=hydra
  HYDRA_PASSWORD=$(openssl rand -hex 16)
  CORE_DB=core
  CORE_USERNAME=core
  CORE_PASSWORD=$(openssl rand -hex 16)
  REDIS_PASSWORD=$(openssl rand -hex 16)

  touch "${REDIS_ENV_FILE}"
  cat <<EOF > "${REDIS_ENV_FILE}"
REDIS_PASSWORD=${REDIS_PASSWORD}
EOF

  touch "${HYDRA_CONFIG_FILE}"
  cat <<EOF > "${HYDRA_CONFIG_FILE}"
secrets:
  system:
    - $(openssl rand -hex 16 | base64)
dsn: postgres://${HYDRA_USERNAME}:${HYDRA_PASSWORD}@db:5432/${HYDRA_DB}?sslmode=disable&max_conns=20&max_idle_conns=4
EOF

  touch "${POSTGRES_ENV_FILE}"
  cat <<EOF >"${POSTGRES_ENV_FILE}"
POSTGRES_USER=${POSTGRES_USER}
POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
HYDRA_DB=${HYDRA_DB}
HYDRA_USERNAME=${HYDRA_USERNAME}
HYDRA_PASSWORD=${HYDRA_PASSWORD}
CORE_DB=${CORE_DB}
CORE_USERNAME=${CORE_USERNAME}
CORE_PASSWORD=${CORE_PASSWORD}
EOF

  touch "${CORE_CONFIG_FILE}"
  cat <<EOF >"${CORE_CONFIG_FILE}"
serve:
  admin:
    cache:
      redis:
        password: ${REDIS_PASSWORD}
    secrets:
      hash:
      - $(openssl rand -hex 64)
      block:
      - $(openssl rand -hex 32)
  public:
    cache:
      redis:
        password: ${REDIS_PASSWORD}
    secrets:
      hash:
      - $(openssl rand -hex 64)
      block:
      - $(openssl rand -hex 32)
  login:
    cache:
      redis:
        password: ${REDIS_PASSWORD}
    secrets:
      hash:
      - $(openssl rand -hex 64)
      block:
      - $(openssl rand -hex 32)
dsn: "host=localhost port=5433 user=${CORE_USERNAME} password=${CORE_PASSWORD} dbname=core sslmode=disable"
EOF
fi
