#!/bin/sh
set -e

find /usr/share/nginx/html/ \
  -type f \( -iname \*.js -o -iname \*.html -o -iname \*.json \) \
  -exec sed -i "s/\${PUBLIC_URL}/${PUBLIC_URL}/g" {} + \
  -exec sed -i "s/\${REACT_APP_OIDC_ISSUER}/${OIDC_ISSUER}/g" {} + \
  -exec sed -i "s/\${REACT_APP_OAUTH_SCOPE}/${OAUTH_SCOPE}/g" {} + \
  -exec sed -i "s/\${REACT_APP_OAUTH_REDIRECT_URI}/${OAUTH_REDIRECT_URI}/g" {} + \
  -exec sed -i "s/\${REACT_APP_OAUTH_CLIENT_ID}/${OAUTH_CLIENT_ID}/g" {} +
