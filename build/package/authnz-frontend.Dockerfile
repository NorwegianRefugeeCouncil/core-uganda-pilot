FROM node:lts-slim as build

RUN mkdir -p /home/node/app && chown -R node:node /home/node/app
USER node
WORKDIR /home/node/app

ADD --chown=node:node frontend/package.json .
ADD --chown=node:node frontend/yarn.lock .
ADD --chown=node:node frontend/lerna.json .
ADD --chown=node:node frontend/packages/core-api-client/package.json ./packages/core-api-client/package.json
ADD --chown=node:node frontend/packages/core-auth/package.json ./packages/core-auth/package.json
ADD --chown=node:node frontend/apps/admin/package.json ./apps/admin/package.json

RUN yarn --immutable

COPY --chown=node:node frontend/. .

# Using this for further envsubst
ENV REACT_APP_OIDC_ISSUER='%{OIDC_ISSUER}%'
ENV REACT_APP_OAUTH_SCOPE='%{OAUTH_SCOPE}%'
ENV REACT_APP_OAUTH_REDIRECT_URI='%{OAUTH_REDIRECT_URI}%'
ENV REACT_APP_OAUTH_CLIENT_ID='%{OAUTH_CLIENT_ID}%'
ENV PUBLIC_URL='%{PUBLIC_URL}%'
ENV REACT_APP_AUTHNZ_API_SERVER_URI='%{AUTHNZ_API_SERVER_URI}%'

RUN yarn build:core-auth && \
    yarn build:core-api-client && \
    yarn build:admin

FROM nginx:alpine

RUN  touch /var/run/nginx.pid && \
     chown -R nginx:nginx /var/cache/nginx /var/run/nginx.pid && \
     rm /etc/nginx/conf.d/default.conf && \
     chown -R nginx:nginx /etc/nginx/conf.d && \
     chown -R nginx:nginx /usr/share/nginx/html/.

USER nginx

ADD --chown=nginx:nginx configs/authnz-frontend/nginx.conf.template  /etc/nginx/templates/nginx.conf.template
ADD --chown=nginx:nginx configs/authnz-frontend/replace-env-vars.sh /docker-entrypoint.d/40-subst-on-assets.sh
COPY --from=build --chown=nginx:nginx /home/node/app/apps/admin/build /usr/share/nginx/html


