#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

echo ">> Creating Core Server OAuth Client..."

RESP=$(
  curl --request POST -L \
    --url 'https://hydra-admin.dev:4445/clients' \
    --data-binary @- <<EOF
{
  "allowed_cors_origins":[],
  "redirect_uris":["$(cat "${SCRIPT_DIR}/../creds/core/app_api/oauth_redirect_uri")"],
  "client_name":"Core App",
  "client_uri":"https://core-app-frontend.dev:3000",
  "client_id":"$(cat "${SCRIPT_DIR}/../creds/core/app_api/oauth_client_id")",
  "client_secret":"$(cat "${SCRIPT_DIR}/../creds/core/app_api/oauth_client_secret")",
  "grant_types":["authorization_code"],
  "scope":"openid offline email profile",
  "response_types":["code"]
}
EOF
)

echo "${RESP}" | grep "a resource with that value exists already" &&
  echo "Found Hydra OAuth client with same id. Updating..." &&
  curl --request PUT -L \
    --url "https://hydra-admin.dev:4445/clients/$(cat "${SCRIPT_DIR}/../creds/core/app_api/oauth_client_id")" \
    --data-binary @- <<EOF
{
 "allowed_cors_origins":[],
 "redirect_uris":["$(cat "${SCRIPT_DIR}/../creds/core/app_api/oauth_redirect_uri")"],
 "client_name":"Core App",
 "client_uri":"https://localhost:3000",
 "client_secret":"$(cat "${SCRIPT_DIR}/../creds/cpre/app_api/oauth_client_secret")",
 "grant_types":["authorization_code"],
 "scope":"openid offline email profile",
 "response_types":["code"]
}
EOF

echo ">> Creating Core Server OAuth Client...Done!"

echo ">> Registering Organization..."

ORG_UUID="$(uuidgen)"
docker exec -it "$(docker ps -aqf "name=db")" /usr/bin/psql \
  -d "$(cat "${SCRIPT_DIR}/../creds/core/db_name")" \
  -U "$(cat "${SCRIPT_DIR}/../creds/core/db_username")" \
  -c "
    INSERT INTO organizations (id, name, email_domain)
    values ('${ORG_UUID}','Norwegian Refugee Council','nrc.no');"

echo ">> Registering Organization... Done!"

echo ">> Registering Organization Identity Provider..."

docker exec -it "$(docker ps -aqf "name=db")" /usr/bin/psql \
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
      'https://oidc.dev:9005',
      '$(cat "${SCRIPT_DIR}/../creds/nrc_idp/oauth_client_id")',
      '$(cat "${SCRIPT_DIR}/../creds/nrc_idp/oauth_client_secret")',
      'nrc.no',
      'Simple Oidc Provider'
    );"

echo ">> Registering Organization Identity Provider...Done!"
