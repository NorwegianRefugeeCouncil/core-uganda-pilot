# Docker

Base: `nginx:alpine`

The Dockerfile for this app is at `/build/package/core-authnz-frontend.Dockerfile`

### Environment Variables

Environment variables expected by the `core-authnz-frontend` app when running it as a container

| Name | Description |
|------|-------------|
| `PUBLIC_URL` | Public URL of the `core-authnz-frontend` app
| `LISTEN_ADDRESS` | nginx [listen](http://nginx.org/en/docs/http/ngx_http_core_module.html#listen) configuration
| `SERVER_NAME` | nginx [server name](http://nginx.org/en/docs/http/server_names.html) configuration
| `OIDC_ISSUER` | URL of the OIDC Issuer
| `OAUTH_SCOPE` | Requested OAuth scope
| `OAUTH_REDIRECT_URI` | OAuth redirect URI
| `OAUTH_CLIENT_ID` | Oauth client ID
| `AUTHNZ_API_SERVER_URI` | `authnz-api-server` URL

