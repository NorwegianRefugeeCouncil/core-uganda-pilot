#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
(cd "${SCRIPT_DIR}/.." && go run "." dev bootstrap)




#echo ">> Creating React Native OAuth Client..."
#
#RESP=$(
#  curl --request POST -L \
#    --url 'http://admin.hydra.core.nrc-tech.dev:4445/clients' \
#    --data-binary @- <<EOF
#{
#  "allowed_cors_origins":[],
#  "redirect_uris":["http://localhost:19006"],
#  "client_name":"Core Native App",
#  "client_uri":"http://localhost:3000",
#  "client_id":"core-react-native-app",
#  "grant_types":["authorization_code", "refresh_token"],
#  "scope":"openid offline email profile",
#  "response_types":["code"],
#  "token_endpoint_auth_method": "none"
#}
#EOF
#)

#echo ">> Creating React Native OAuth Client...Done!"
#
#echo ">> Registering Organization..."
#
#ORG_UUID="$(uuidgen)"
#docker exec -it "$(docker ps -aqf "name=db")" /usr/bin/psql \
#  -d "$(cat "${SCRIPT_DIR}/../creds/core/db_name")" \
#  -U "$(cat "${SCRIPT_DIR}/../creds/core/db_username")" \
#  -c "
#    INSERT INTO organizations (id, name, email_domain)
#    values ('${ORG_UUID}','Norwegian Refugee Council','nrc.no');"
#
#echo ">> Registering Organization... Done!"
#
#echo ">> Registering Organization Identity Provider..."
#
#docker exec -it "$(docker ps -aqf "name=db")" /usr/bin/psql \
#  -d core \
#  -U core \
#  -c "
#    INSERT INTO identity_providers
#    (
#      id,
#      organization_id,
#      domain,
#      client_id,
#      client_secret,
#      email_domain,
#      name
#    )
#    values
#    (
#      '$(uuidgen)',
#      '${ORG_UUID}',
#      'http://localhost:9005',
#      '$(cat "${SCRIPT_DIR}/../creds/nrc_idp/oauth_client_id")',
#      '$(cat "${SCRIPT_DIR}/../creds/nrc_idp/oauth_client_secret")',
#      'nrc.no',
#      'Simple Oidc Provider'
#    );"
#
#echo ">> Registering Organization Identity Provider...Done!"
