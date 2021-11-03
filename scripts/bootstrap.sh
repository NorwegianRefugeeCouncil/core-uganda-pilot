#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

PWA_ENV_FILE="${SCRIPT_DIR}/../web/pwa/.env"

echo Creating Core App Hydra Client
CLIENT=$(
  curl --request POST -sL \
    --url 'http://localhost:4445/clients' \
    --data-binary @- <<EOF
{
  "allowed_cors_origins":["http://localhost:3000"],
  "redirect_uris":["http://localhost:3000"],
  "client_name":"Core React App",
  "client_uri":"http://localhost:3000",
  "grant_types":["authorization_code"],
  "token_endpoint_auth_method":"none",
  "scope":"openid offline offline_access email profile",
  "response_types":["id_token","refresh_token","code"]
}
EOF
)

echo Creating Core App env file
touch "${PWA_ENV_FILE}"
cat <<EOF >"${PWA_ENV_FILE}"
REACT_APP_CLIENT_ID=$(echo "${CLIENT}" | jq -r ".client_id")
REACT_APP_ISSUER=http://localhost:4444
REACT_APP_REDIRECT_URI=http://localhost:3000
REACT_APP_SILENT_REDIRECT_URI=http://localhost:3000
EOF

echo Registering Organization
ORG_UUID="$(uuidgen)"
docker exec -it "$(docker ps -aqf "name=core_db")" /usr/bin/psql \
  -d core \
  -U core \
  -c "
    INSERT INTO organizations (id, name, email_domain)
    values ('${ORG_UUID}','Norwegian Refugee Council','nrc.no');"

CLIENT_INFO=$(jq '.client_config | .[] | select( .client_id == "core-login" )' <"${SCRIPT_DIR}/../deployments/oidc.config.json")
CLIENT_ID=$(echo "${CLIENT_INFO}" | jq -r .client_id)
CLIENT_SECRET=$(echo "${CLIENT_INFO}" | jq -r .client_secret)

echo Registering Organization Identity Provider
docker exec -it "$(docker ps -aqf "name=core_db")" /usr/bin/psql \
  -d core \
  -U core \
  -c "
    INSERT INTO identity_providers
    (
      id,
      organization_id,
      domain,
      client_id,
      client_secret,
      email_domain,
      name
    )
    values
    (
      '$(uuidgen)',
      '${ORG_UUID}',
      'http://localhost:9005',
      '${CLIENT_ID}',
      '${CLIENT_SECRET}',
      'nrc.no',
      'Simple Oidc Provider'
    );"
