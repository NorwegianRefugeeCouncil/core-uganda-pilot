echo ">> Creating Core Server OAuth Client..."

RESP=$(
  curl --request POST -L \
    --url 'http://localhost:4445/clients' \
    --data-binary @- <<EOF
{
  "allowed_cors_origins":[],
  "redirect_uris":["http://localhost:19006"],
  "client_name":"Core Native App",
  "client_uri":"http://localhost:3000",
  "client_id":"core-react-native-app",
  "grant_types":["authorization_code", "refresh_token"],
  "scope":"openid offline email profile",
  "response_types":["code"],
  "token_endpoint_auth_method": "none"
}
EOF
)

echo "${RESP}" | grep "a resource with that value exists already" &&
  echo "Found Hydra OAuth client with same id. Updating..." &&
  curl --request PUT -L \
    --url "http://localhost:4445/clients/react-native" \
    --data-binary @- <<EOF
{
  "allowed_cors_origins":[],
  "redirect_uris":["http://localhost:19006"],
  "client_name":"Core Native App",
  "client_uri":"http://localhost:3000",
  "client_id":"core-react-native-app",
  "grant_types":["authorization_code", "refresh_token"],
  "scope":"openid offline email profile",
  "response_types":["code"],
  "token_endpoint_auth_method": "none"
}
EOF
