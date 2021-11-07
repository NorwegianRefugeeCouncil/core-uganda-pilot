#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
. "${SCRIPT_DIR}/utils.sh"

mkdir -p "${SCRIPT_DIR}/../creds"
CREDS_DIR="$(cd "${SCRIPT_DIR}/../creds" && pwd)"

HYDRA_SYSTEM_SECRET=$(createFileIfNotExists "${CREDS_DIR}/hydra/system_secret" "openssl rand -hex 32")
HYDRA_COOKIE_SECRET=$(createFileIfNotExists "${CREDS_DIR}/hydra/cookie_secret" "openssl rand -hex 32")
POSTGRES_ROOT_PASSWORD=$(createFileIfNotExists "${CREDS_DIR}/postgres/root_password" "openssl rand -hex 32")
POSTGRES_ROOT_USERNAME=$(createFileIfNotExists "${CREDS_DIR}/postgres/root_username" "echo -n postgres")
POSTGRES_CORE_DB=$(createFileIfNotExists "${CREDS_DIR}/core//db_name" "echo -n core")
POSTGRES_CORE_USERNAME=$(createFileIfNotExists "${CREDS_DIR}/core/db_username" "echo core")
POSTGRES_CORE_PASSWORD=$(createFileIfNotExists "${CREDS_DIR}/core/db_password" "openssl rand -hex 32")
POSTGRES_HYDRA_DB=$(createFileIfNotExists "${CREDS_DIR}/hydra/db_name" "echo -n hydra")
POSTGRES_HYDRA_USERNAME=$(createFileIfNotExists "${CREDS_DIR}/hydra/db_username" "echo -n hydra")
POSTGRES_HYDRA_PASSWORD=$(createFileIfNotExists "${CREDS_DIR}/hydra/db_password" "openssl rand -hex 32")
REDIS_PASSWORD=$(createFileIfNotExists "${CREDS_DIR}/redis/password" "openssl rand -hex 32")
OAUTH_CORE_ADMIN_CLIENT_ID=$(createFileIfNotExists "${CREDS_DIR}/core/admin_api/oauth_client_id" "echo -n core-admin")
OAUTH_CORE_ADMIN_CLIENT_SECRET=$(createFileIfNotExists "${CREDS_DIR}/core/admin_api/oauth_client_secret" "openssl rand -hex 32")
OAUTH_CORE_ADMIN_REDIRECT_URI=$(createFileIfNotExists "${CREDS_DIR}/core/admin_api/oauth_redirect_uri" "echo -n https://localhost:9001/oidc/callback")
OAUTH_CORE_ADMIN_ISSUER=$(createFileIfNotExists "${CREDS_DIR}/core/admin_api/oauth_issuer" "echo -n https://localhost:9005")
OAUTH_CORE_APP_CLIENT_ID=$(createFileIfNotExists "${CREDS_DIR}/core/app_api/oauth_client_id" "echo -n core-app")
OAUTH_CORE_APP_CLIENT_SECRET=$(createFileIfNotExists "${CREDS_DIR}/core/app_api/oauth_client_secret" "openssl rand -hex 32")
OAUTH_CORE_APP_REDIRECT_URI=$(createFileIfNotExists "${CREDS_DIR}/core/app_api/oauth_redirect_uri" "echo -n https://localhost:9000/oidc/callback")
OAUTH_CORE_APP_ISSUER=$(createFileIfNotExists "${CREDS_DIR}/core/app_api/oauth_issuer" "echo -n https://localhost:4444/")
OAUTH_NRC_CLIENT_ID=$(createFileIfNotExists "${CREDS_DIR}/nrc_idp/oauth_client_id" "echo -n nrc-idp")
OAUTH_NRC_CLIENT_SECRET=$(createFileIfNotExists "${CREDS_DIR}/nrc_idp/oauth_client_secret" "openssl rand -hex 32")
OAUTH_NRC_REDIRECT_URI=$(createFileIfNotExists "${CREDS_DIR}/nrc_idp/oauth_redirect_uri" "echo -n https://localhost:9002/oidc/callback")
OAUTH_NRC_ISSUER=$(createFileIfNotExists "${CREDS_DIR}/nrc_idp/oauth_issuer" "echo -n https://localhost:9005")
CORE_ADMIN_HASH_KEY=$(createFileIfNotExists "${CREDS_DIR}/core/admin_api/secret_hash_key" "openssl rand -hex 64")
CORE_ADMIN_BLOCK_KEY=$(createFileIfNotExists "${CREDS_DIR}/core/admin_api/secret_block_key" "openssl rand -hex 32")
CORE_LOGIN_HASH_KEY=$(createFileIfNotExists "${CREDS_DIR}/core/login/secret_hash_key" "openssl rand -hex 64")
CORE_LOGIN_BLOCK_KEY=$(createFileIfNotExists "${CREDS_DIR}/core/login/secret_block_key" "openssl rand -hex 32")
CORE_APP_HASH_KEY=$(createFileIfNotExists "${CREDS_DIR}/core/app_api/secret_hash_key" "openssl rand -hex 64")
CORE_APP_BLOCK_KEY=$(createFileIfNotExists "${CREDS_DIR}/core/app_api/secret_block_key" "openssl rand -hex 32")

OIDC_CONFIG_FILE="${CREDS_DIR}/oidc/config.json"
OIDC_USERS_FILE="${CREDS_DIR}/oidc/users.json"
REDIS_ENV_FILE="${CREDS_DIR}/redis/env"
HYDRA_CONFIG_FILE="${CREDS_DIR}/hydra/config.yaml"
POSTGRES_ENV_FILE="${CREDS_DIR}/postgres/env"
POSTGRES_INIT_FILE="${CREDS_DIR}/postgres/init.sh"
CORE_CONFIG_FILE="${SCRIPT_DIR}/../creds/core/config.yaml"

# TLS
echo ">> Generating root CA"
CA_KEY_FILE="${CREDS_DIR}/ca/tls.key"
CA_CSR_FILE="${CREDS_DIR}/ca/cert.csr"
CA_CERT_FILE="${CREDS_DIR}/ca/tls.cert"
createFileIfNotExists "${CA_KEY_FILE}" "openssl genrsa 2048"
createFileIfNotExists "${CA_CSR_FILE}" "openssl req -new -key ${CA_KEY_FILE} -subj '/C=DE/ST=Berlin/L=Berlin/O=NRC/CN=core.dev'"
createFileIfNotExists "${CA_CERT_FILE}" "openssl x509 -in ${CA_CSR_FILE} -req -signkey ${CA_KEY_FILE} -days 365"

function createCert() {
  local DIR
  local NAME
  local DOMAIN
  local KEY_FILE
  local CSR_FILE
  local CERT_FILE
  local SERIAL_FILE
  DIR=$1
  NAME=$2
  DOMAIN=$3
  echo ">> Generating Server Certificate for ${NAME} (${DOMAIN})"
  mkdir -p "${CREDS_DIR}/${DIR}"
  KEY_FILE="${CREDS_DIR}/${DIR}/${NAME}_tls.key"
  CSR_FILE="${CREDS_DIR}/${DIR}/${NAME}.csr"
  CERT_FILE="${CREDS_DIR}/${DIR}/${NAME}_tls.cert"
  SERIAL_FILE="${CREDS_DIR}/${DIR}/${NAME}_serial"
  SSL_CONF_FILE="${CREDS_DIR}/${DIR}/${NAME}_ssl.conf"
  if [ ! -f "${SSL_CONF_FILE}" ]; then
    cat <<EOF >"${SSL_CONF_FILE}"
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
  createFileIfNotExists "${KEY_FILE}" "openssl genrsa 2048"
  createFileIfNotExists "${CSR_FILE}" "openssl req -new -key ${KEY_FILE} -subj '/C=DE/ST=Berlin/L=Berlin/O=NRC/CN=${DOMAIN}'"
  createFileIfNotExists "${CERT_FILE}" "openssl x509 -req -in ${CSR_FILE} -CA ${CA_CERT_FILE} -CAkey ${CA_KEY_FILE} -CAcreateserial -CAserial ${SERIAL_FILE} -days 1825 -sha256 -extfile ${SSL_CONF_FILE}"
}

createCert core/admin_api admin_api core-admin-api.dev
createCert core/admin_frontend admin_frontend core-admin-frontend.dev
createCert core/app_api app_api core-app-api.dev
createCert core/app_frontend app_frontend core-app-frontend.dev
createCert core/login login core-login.dev
createCert hydra hydra_public hydra.dev
createCert hydra hydra_admin hydra-admin.dev
createCert oidc oidc oidc.dev

echo ">> Creating Simple-OIDC Clients"

mkdir -p "$(dirname "${OIDC_CONFIG_FILE}")"
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
      "scope": "openid email profile",
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
    "nickname": "harley"
  },
  {
    "id": "barb",
    "email": "barb@nrc.no",
    "email_verified": true,
    "name": "Barb Stovin",
    "nickname": "barb"
  },
  {
    "id": "quinn",
    "email": "quinn@nrc.no",
    "email_verified": true,
    "name": "Quinn Leeming",
    "nickname": "quinn"
  },
  {
    "id": "sim",
    "email": "sim@nrc.no",
    "email_verified": true,
    "name": "Sim Cleaton",
    "nickname": "sim"
  },
  {
    "id": "phillie",
    "email": "phillie@nrc.no",
    "email_verified": true,
    "name": "Phillie Smeed",
    "nickname": "phillie"
  },
  {
    "id": "peta",
    "email": "peta@nrc.no",
    "email_verified": true,
    "name": "Peta Sammon",
    "nickname": "peta"
  },
  {
    "id": "marne",
    "email": "marne@nrc.no",
    "email_verified": true,
    "name": "Marne Probetts",
    "nickname": "marne"
  },
  {
    "id": "sibylla",
    "email": "sibylla@nrc.no",
    "email_verified": true,
    "name": "Sibylla Meadows",
    "nickname": "sibylla"
  },
  {
    "id": "evan",
    "email": "evan@nrc.no",
    "email_verified": true,
    "name": "Evan Highman",
    "nickname": "evan"
  },
  {
    "id": "franklin",
    "email": "franklin@nrc.no",
    "email_verified": true,
    "name": "Franklin Glamart",
    "nickname": "franklin"
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
EOF

echo ">> Creating Postgres Init File"

touch "${POSTGRES_INIT_FILE}"
cat <<EOF >"${POSTGRES_INIT_FILE}"
#!/bin/bash

set -e
set -u

function create_user_and_database() {
	local database=\$1
	local user=\$2
	local password=\$3
	echo ">> Creating user '\$user' and database '\$database'"
	psql -v ON_ERROR_STOP=1 --username "\$POSTGRES_USER" <<-EOSQL
      CREATE USER \$user WITH PASSWORD '\$password';
	    CREATE DATABASE \$database;
	    GRANT ALL PRIVILEGES ON DATABASE \$database TO \$user;
EOSQL
}

create_user_and_database "${POSTGRES_HYDRA_DB}" "${POSTGRES_HYDRA_USERNAME}" "${POSTGRES_HYDRA_PASSWORD}"
create_user_and_database "${POSTGRES_CORE_DB}" "${POSTGRES_CORE_USERNAME}" "${POSTGRES_CORE_PASSWORD}"

EOF

chmod +x "${POSTGRES_INIT_FILE}"

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
