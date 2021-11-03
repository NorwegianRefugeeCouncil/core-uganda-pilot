#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

CORE_CONFIG_FILE="${SCRIPT_DIR}/../configs/config.custom.yaml"
HYDRA_CONFIG_FILE="${SCRIPT_DIR}/../deployments/hydra.custom.yaml"
POSTGRES_ENV_FILE="${SCRIPT_DIR}/../deployments/postgres.env"
REDIS_ENV_FILE="${SCRIPT_DIR}/../deployments/redis.env"
OIDC_CONFIG_FILE="${SCRIPT_DIR}/../deployments/oidc.config.json"
OIDC_USERS_FILE="${SCRIPT_DIR}/../deployments/oidc.users.json"
ADMIN_APP_ENV_FILE="${SCRIPT_DIR}/../web/admin/.env"

if [ ! -f "${OIDC_USERS_FILE}" ] || [ ! -f "${OIDC_CONFIG_FILE}" ] || [ ! -f "${REDIS_ENV_FILE}" ] || [ ! -f "${CORE_CONFIG_FILE}" ] || [ ! -f "${POSTGRES_ENV_FILE}" ] || [ ! -f "${HYDRA_CONFIG_FILE}" ]; then

  POSTGRES_USER=postgres
  POSTGRES_PASSWORD=$(openssl rand -hex 16)
  HYDRA_DB=hydra
  HYDRA_USERNAME=hydra
  HYDRA_PASSWORD=$(openssl rand -hex 16)
  CORE_DB=core
  CORE_USERNAME=core
  CORE_PASSWORD=$(openssl rand -hex 16)
  REDIS_PASSWORD=$(openssl rand -hex 16)

  OIDC_ADMIN_CLIENT_SECRET=$(openssl rand -hex 16)
  OIDC_ADMIN_CLIENT_ID="core-admin"
  OIDC_ADMIN_REDIRECT_URI="http://localhost:3001"

  OIDC_LOGIN_CLIENT_SECRET=$(openssl rand -hex 16)
  OIDC_LOGIN_CLIENT_ID="core-login"
  OIDC_LOGIN_REDIRECT_URI="http://localhost:9002/oidc/callback"

  touch "${ADMIN_APP_ENV_FILE}"
  cat <<EOF >"${ADMIN_APP_ENV_FILE}"
REACT_APP_CLIENT_ID=${OIDC_ADMIN_CLIENT_ID}
REACT_APP_ISSUER=http://localhost:9005
REACT_APP_REDIRECT_URI=http://localhost:3001
REACT_APP_SILENT_REDIRECT_URI=http://localhost:3001
EOF

  touch "${OIDC_CONFIG_FILE}"
  cat <<EOF >"${OIDC_CONFIG_FILE}"
{
  "client_config": [
    {
      "client_id": "${OIDC_ADMIN_CLIENT_ID}",
      "client_secret": "${OIDC_ADMIN_CLIENT_SECRET}",
      "redirect_uris": [ "${OIDC_ADMIN_REDIRECT_URI}" ],
      "grant_types": [
        "authorization_code"
      ],
      "token_endpoint_auth_method": "none",
      "scope": "openid offline_access",
      "response_types": [
        "code"
      ]
    },{
      "client_id": "${OIDC_LOGIN_CLIENT_ID}",
      "client_secret": "${OIDC_LOGIN_CLIENT_SECRET}",
      "redirect_uris": [ "${OIDC_LOGIN_REDIRECT_URI}" ],
      "grant_types": [
        "authorization_code",
        "refresh_token"
      ],
      "token_endpoint_auth_method": "client_secret_post",
      "scope": "openid email profile offline_access",
      "response_types": [
        "code"
      ]
    }
  ]
}
EOF

  touch "${OIDC_USERS_FILE}"
  cat <<EOF >"${OIDC_USERS_FILE}"
[
  {
    "id": "admin",
    "email": "admin@nrc.no",
    "email_verified": true,
    "name": "Harley Kiffe",
    "nickname": "harley",
    "password": "$(openssl rand -hex 16)"
  },
  {
    "id": "barb",
    "email": "barb@nrc.no",
    "email_verified": true,
    "name": "Barb Stovin",
    "nickname": "barb",
    "password": "$(openssl rand -hex 16)"
  },
  {
    "id": "quinn",
    "email": "quinn@nrc.no",
    "email_verified": true,
    "name": "Quinn Leeming",
    "nickname": "quinn",
    "password": "$(openssl rand -hex 16)"
  },
  {
    "id": "sim",
    "email": "sim@nrc.no",
    "email_verified": true,
    "name": "Sim Cleaton",
    "nickname": "sim",
    "password": "$(openssl rand -hex 16)"
  },
  {
    "id": "phillie",
    "email": "phillie@nrc.no",
    "email_verified": true,
    "name": "Phillie Smeed",
    "nickname": "phillie",
    "password": "$(openssl rand -hex 16)"
  },
  {
    "id": "peta",
    "email": "peta@nrc.no",
    "email_verified": true,
    "name": "Peta Sammon",
    "nickname": "peta",
    "password": "$(openssl rand -hex 16)"
  },
  {
    "id": "marne",
    "email": "marne@nrc.no",
    "email_verified": true,
    "name": "Marne Probetts",
    "nickname": "marne",
    "password": "$(openssl rand -hex 16)"
  },
  {
    "id": "sibylla",
    "email": "sibylla@nrc.no",
    "email_verified": true,
    "name": "Sibylla Meadows",
    "nickname": "sibylla",
    "password": "$(openssl rand -hex 16)"
  },
  {
    "id": "evan",
    "email": "evan@nrc.no",
    "email_verified": true,
    "name": "Evan Highman",
    "nickname": "evan",
    "password": "$(openssl rand -hex 16)"
  },
  {
    "id": "franklin",
    "email": "franklin@nrc.no",
    "email_verified": true,
    "name": "Franklin Glamart",
    "nickname": "franklin",
    "password": "$(openssl rand -hex 16)"
  }
]
EOF

  touch "${REDIS_ENV_FILE}"
  cat <<EOF >"${REDIS_ENV_FILE}"
REDIS_PASSWORD=${REDIS_PASSWORD}
EOF

  touch "${HYDRA_CONFIG_FILE}"
  cat <<EOF >"${HYDRA_CONFIG_FILE}"
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
