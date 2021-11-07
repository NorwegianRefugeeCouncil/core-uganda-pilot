#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
. "${SCRIPT_DIR}/utils.sh"

CORE_CONFIG_FILE="${SCRIPT_DIR}/../configs/config.custom.yaml"
HYDRA_CONFIG_FILE="${SCRIPT_DIR}/../deployments/hydra.custom.yaml"
POSTGRES_ENV_FILE="${SCRIPT_DIR}/../deployments/postgres.env"
REDIS_ENV_FILE="${SCRIPT_DIR}/../deployments/redis.env"
OIDC_CONFIG_FILE="${SCRIPT_DIR}/../deployments/oidc.config.json"
OIDC_USERS_FILE="${SCRIPT_DIR}/../deployments/oidc.users.json"

CREDS_DIR="${SCRIPT_DIR}/../creds"

function createUserPassword() {
  filePath="${CREDS_DIR}/user_$1_password"
  createFileIfNotExists "${filePath}" "openssl rand -hex 32"
}

HYDRA_SYSTEM_SECRET=$(createFileIfNotExists "${CREDS_DIR}/hydra_system_secret" "openssl rand -hex 32")
HYDRA_COOKIE_SECRET=$(createFileIfNotExists "${CREDS_DIR}/hydra_cookie_secret" "openssl rand -hex 32")
POSTGRES_ROOT_PASSWORD=$(createFileIfNotExists "${CREDS_DIR}/postgres_root_password" "openssl rand -hex 32")
POSTGRES_ROOT_USERNAME=$(createFileIfNotExists "${CREDS_DIR}/postgres_root_username" "echo -n postgres")
POSTGRES_CORE_DB=$(createFileIfNotExists "${CREDS_DIR}/postgres_core_db" "echo -n core")
POSTGRES_CORE_USERNAME=$(createFileIfNotExists "${CREDS_DIR}/postgres_core_username" "echo core")
POSTGRES_CORE_PASSWORD=$(createFileIfNotExists "${CREDS_DIR}/postgres_core_password" "openssl rand -hex 32")
POSTGRES_HYDRA_DB=$(createFileIfNotExists "${CREDS_DIR}/postgres_hydra_db" "echo -n hydra")
POSTGRES_HYDRA_USERNAME=$(createFileIfNotExists "${CREDS_DIR}/postgres_hydra_username" "echo -n hydra")
POSTGRES_HYDRA_PASSWORD=$(createFileIfNotExists "${CREDS_DIR}/postgres_hydra_password" "openssl rand -hex 32")
REDIS_PASSWORD=$(createFileIfNotExists "${CREDS_DIR}/redis_password" "openssl rand -hex 32")
OAUTH_CORE_ADMIN_CLIENT_ID=$(createFileIfNotExists "${CREDS_DIR}/core_admin_client_id" "echo -n core-admin")
OAUTH_CORE_ADMIN_CLIENT_SECRET=$(createFileIfNotExists "${CREDS_DIR}/core_admin_client_secret" "openssl rand -hex 32")
OAUTH_CORE_ADMIN_REDIRECT_URI=$(createFileIfNotExists "${CREDS_DIR}/core_admin_redirect_uri" "echo -n http://localhost:9001/oidc/callback")
OAUTH_CORE_ADMIN_ISSUER=$(createFileIfNotExists "${CREDS_DIR}/core_admin_issuer" "echo -n http://localhost:9005")
OAUTH_CORE_APP_CLIENT_ID=$(createFileIfNotExists "${CREDS_DIR}/core_app_client_id" "core-app")
OAUTH_CORE_APP_CLIENT_SECRET=$(createFileIfNotExists "${CREDS_DIR}/core_app_client_secret" "openssl rand -hex 32")
OAUTH_CORE_APP_REDIRECT_URI=$(createFileIfNotExists "${CREDS_DIR}/core_app_redirect_uri" "echo -n http://localhost:9000/oidc/callback")
OAUTH_CORE_APP_ISSUER=$(createFileIfNotExists "${CREDS_DIR}/core_app_issuer" "echo -n http://localhost:4444/")
OAUTH_NRC_CLIENT_ID=$(createFileIfNotExists "${CREDS_DIR}/nrc_idp_client_id" "nrc-idp")
OAUTH_NRC_CLIENT_SECRET=$(createFileIfNotExists "${CREDS_DIR}/nrc_idp_client_secret" "openssl rand -hex 32")
OAUTH_NRC_REDIRECT_URI=$(createFileIfNotExists "${CREDS_DIR}/nrc_idp_redirect_uri" "echo -n http://localhost:9002/oidc/callback")
OAUTH_NRC_ISSUER=$(createFileIfNotExists "${CREDS_DIR}/nrc_idp_issuer" "echo -n http://localhost:9005")
CORE_ADMIN_HASH_KEY=$(createFileIfNotExists "${CREDS_DIR}/core_admin_hash_key" "openssl rand -hex 64")
CORE_ADMIN_BLOCK_KEY=$(createFileIfNotExists "${CREDS_DIR}/core_admin_block_key" "openssl rand -hex 32")
CORE_LOGIN_HASH_KEY=$(createFileIfNotExists "${CREDS_DIR}/core_login_hash_key" "openssl rand -hex 64")
CORE_LOGIN_BLOCK_KEY=$(createFileIfNotExists "${CREDS_DIR}/core_login_block_key" "openssl rand -hex 32")
CORE_APP_HASH_KEY=$(createFileIfNotExists "${CREDS_DIR}/core_app_hash_key" "openssl rand -hex 64")
CORE_APP_BLOCK_KEY=$(createFileIfNotExists "${CREDS_DIR}/core_app_block_key" "openssl rand -hex 32")

# TLS
createFileIfNotExists "${CREDS_DIR}/ca_key" "openssl genrsa 2048"
createFileIfNotExists "${CREDS_DIR}/ca_csr" "openssl req -new -key ${CREDS_DIR}/ca_key -subj '/C=DE/ST=Berlin/L=Berlin/O=NRC/CN=core.dev'"
createFileIfNotExists "${CREDS_DIR}/ca_cert" "openssl x509 -in ${CREDS_DIR}/ca_csr -req -signkey ${CREDS_DIR}/ca_key -days 365"

function createCert() {
  local NAME
  local DOMAIN
  local SSL_CONF
  NAME=$1
  DOMAIN=$2
  SSL_CONF="${CREDS_DIR}/${NAME}_ssl.conf"
  if [ ! -f "${SSL_CONF}" ]; then
    cat <<EOF > "${SSL_CONF}"
authorityKeyIdentifier=keyid,issuer
basicConstraints=CA:FALSE
keyUsage = digitalSignature, nonRepudiation, keyEncipherment, dataEncipherment
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
DNS.2 = ${DOMAIN}
IP.1 = 127.0.0.1
EOF
  fi
  createFileIfNotExists "${CREDS_DIR}/${NAME}_key" "openssl genrsa 2048"
  createFileIfNotExists "${CREDS_DIR}/${NAME}_csr" "openssl req -new -key ${CREDS_DIR}/${NAME}_key -subj '/C=DE/ST=Berlin/L=Berlin/O=NRC/CN=${DOMAIN}'"
  createFileIfNotExists "${CREDS_DIR}/${NAME}_cert" "openssl x509 -req -in ${CREDS_DIR}/${NAME}_csr -CA ${CREDS_DIR}/ca_cert -CAkey ${CREDS_DIR}/ca_key -CAcreateserial -CAserial ${CREDS_DIR}/${NAME}_serial -days 1825 -sha256 -extfile ${SSL_CONF}"
}

createCert core_admin_api core-admin-api.dev
createCert core_admin_frontend core-admin-frontend.dev
createCert core_app_api core-app-api.dev
createCert core_app_frontend core-app-frontend.dev
createCert core_login core-login.dev
createCert hydra hydra.dev
createCert hydra_admin hydra-admin.dev

echo ">> Creating Simple-OIDC Clients"

touch "${OIDC_CONFIG_FILE}"
cat <<EOF >"${OIDC_CONFIG_FILE}"
{
  "scopes":[
    "openid",
    "profile",
    "email",
    "offline_access"
  ],
  "clients": [
    {
      "client_id": "${OAUTH_CORE_ADMIN_CLIENT_ID}",
      "client_secret": "${OAUTH_CORE_ADMIN_CLIENT_SECRET}",
      "redirect_uris": [ "${OAUTH_CORE_ADMIN_REDIRECT_URI}" ],
      "grant_types": [
        "authorization_code",
        "refresh_token"
      ],
      "token_endpoint_auth_method": "client_secret_post",
      "scope": "openid",
      "response_types": [
        "code"
      ]
    },{
      "client_id": "${OAUTH_NRC_CLIENT_ID}",
      "client_secret": "${OAUTH_NRC_CLIENT_SECRET}",
      "redirect_uris": [ "${OAUTH_NRC_REDIRECT_URI}" ],
      "grant_types": [
        "authorization_code",
        "refresh_token"
      ],
      "token_endpoint_auth_method": "client_secret_post",
      "scope": "openid email profile",
      "response_types": [
        "code"
      ]
    }
  ]
}
EOF

echo ">> Creating Simple-OIDC Users"

touch "${OIDC_USERS_FILE}"
cat <<EOF >"${OIDC_USERS_FILE}"
[
  {
    "id": "admin",
    "email": "admin@nrc.no",
    "email_verified": true,
    "name": "Harley Kiffe",
    "nickname": "harley",
    "password": "$(createUserPassword admin)"
  },
  {
    "id": "barb",
    "email": "barb@nrc.no",
    "email_verified": true,
    "name": "Barb Stovin",
    "nickname": "barb",
    "password": "$(createUserPassword barb)"
  },
  {
    "id": "quinn",
    "email": "quinn@nrc.no",
    "email_verified": true,
    "name": "Quinn Leeming",
    "nickname": "quinn",
    "password": "$(createUserPassword quinn)"
  },
  {
    "id": "sim",
    "email": "sim@nrc.no",
    "email_verified": true,
    "name": "Sim Cleaton",
    "nickname": "sim",
    "password": "$(createUserPassword sim)"
  },
  {
    "id": "phillie",
    "email": "phillie@nrc.no",
    "email_verified": true,
    "name": "Phillie Smeed",
    "nickname": "phillie",
    "password": "$(createUserPassword phillie)"
  },
  {
    "id": "peta",
    "email": "peta@nrc.no",
    "email_verified": true,
    "name": "Peta Sammon",
    "nickname": "peta",
    "password": "$(createUserPassword peta)"
  },
  {
    "id": "marne",
    "email": "marne@nrc.no",
    "email_verified": true,
    "name": "Marne Probetts",
    "nickname": "marne",
    "password": "$(createUserPassword marne)"
  },
  {
    "id": "sibylla",
    "email": "sibylla@nrc.no",
    "email_verified": true,
    "name": "Sibylla Meadows",
    "nickname": "sibylla",
    "password": "$(createUserPassword sibylla)"
  },
  {
    "id": "evan",
    "email": "evan@nrc.no",
    "email_verified": true,
    "name": "Evan Highman",
    "nickname": "evan",
    "password": "$(createUserPassword evan)"
  },
  {
    "id": "franklin",
    "email": "franklin@nrc.no",
    "email_verified": true,
    "name": "Franklin Glamart",
    "nickname": "franklin",
    "password": "$(createUserPassword franklin)"
  }
]
EOF

echo ">> Creating Redis Env File"

touch "${REDIS_ENV_FILE}"
cat <<EOF >"${REDIS_ENV_FILE}"
REDIS_PASSWORD=${REDIS_PASSWORD}
EOF

echo ">> Creating Hydra Config File"

touch "${HYDRA_CONFIG_FILE}"
cat <<EOF >"${HYDRA_CONFIG_FILE}"
secrets:
  system:
    - ${HYDRA_SYSTEM_SECRET}
  cookie:
    - ${HYDRA_COOKIE_SECRET}
dsn: postgres://${POSTGRES_HYDRA_USERNAME}:${POSTGRES_HYDRA_PASSWORD}@db:5432/${POSTGRES_HYDRA_DB}?sslmode=disable&max_conns=20&max_idle_conns=4
EOF

echo ">> Creating Postgres Env File"

touch "${POSTGRES_ENV_FILE}"
cat <<EOF >"${POSTGRES_ENV_FILE}"
POSTGRES_USER=${POSTGRES_ROOT_USERNAME}
POSTGRES_PASSWORD=${POSTGRES_ROOT_PASSWORD}
HYDRA_DB=${POSTGRES_HYDRA_DB}
HYDRA_USERNAME=${POSTGRES_HYDRA_USERNAME}
HYDRA_PASSWORD=${POSTGRES_HYDRA_PASSWORD}
CORE_DB=${POSTGRES_CORE_DB}
CORE_USERNAME=${POSTGRES_CORE_USERNAME}
CORE_PASSWORD=${POSTGRES_CORE_PASSWORD}
EOF

echo ">> Creating Core Config File"

touch "${CORE_CONFIG_FILE}"
cat <<EOF >"${CORE_CONFIG_FILE}"
serve:
  admin:
    cache:
      redis:
        password: ${REDIS_PASSWORD}
    secrets:
      hash:
      - ${CORE_ADMIN_HASH_KEY}
      block:
      - ${CORE_ADMIN_BLOCK_KEY}
    oidc:
      client_id: ${OAUTH_CORE_ADMIN_CLIENT_ID}
      client_secret: ${OAUTH_CORE_ADMIN_CLIENT_SECRET}
      issuer: ${OAUTH_CORE_ADMIN_ISSUER}
  public:
    cache:
      redis:
        password: ${REDIS_PASSWORD}
    secrets:
      hash:
      - ${CORE_APP_HASH_KEY}
      block:
      - ${CORE_APP_BLOCK_KEY}
    oidc:
      client_id: ${OAUTH_CORE_APP_CLIENT_ID}
      client_secret: ${OAUTH_CORE_APP_CLIENT_SECRET}
      issuer: ${OAUTH_CORE_APP_ISSUER}
  login:
    cache:
      redis:
        password: ${REDIS_PASSWORD}
    secrets:
      hash:
      - ${CORE_LOGIN_HASH_KEY}
      block:
      - ${CORE_LOGIN_BLOCK_KEY}
dsn: "host=localhost port=5433 user=${POSTGRES_CORE_USERNAME} password=${POSTGRES_CORE_PASSWORD} dbname=${POSTGRES_CORE_DB} sslmode=disable"
EOF
