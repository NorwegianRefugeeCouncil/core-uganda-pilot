#!/bin/sh
set -e

find /usr/share/nginx/html/ \
  -type f \( -iname \*.js -o -iname \*.html -o -iname \*.json -o -iname \*.map \) \
  -exec sed -i "s~%{PUBLIC_URL}%~${PUBLIC_URL}~g" {} + \
  -exec sed -i "s~%{OIDC_ISSUER}%~${OIDC_ISSUER}~g" {} + \
  -exec sed -i "s~%{OAUTH_SCOPE}%~${OAUTH_SCOPE}~g" {} + \
  -exec sed -i "s~%{OAUTH_REDIRECT_URI}%~${OAUTH_REDIRECT_URI}~g" {} + \
  -exec sed -i "s~%{OAUTH_CLIENT_ID}%~${OAUTH_CLIENT_ID}~g" {} + \
  -exec sed -i "s~%{SERVER_URL}%~${SERVER_URL}~g" {} +
