# Docker

Base: `nginx:alpine`

The Dockerfile for this app is at `/build/package/authnz-frontend.Dockerfile`

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

### Environment variable substitution

The problem with using environment variables with react apps is that these will be built-in at compile time.

What if we want to deploy this application multiple times, and want to use different environment variables? e.g. for
production, staging, etc.

The solution is to use fake, known values for environment variables, that we can then replace again when starting the
nginx server. Eg `REACT_APP_OIDC_ISSUER=FAKE_OIDC_ISSUER`

When launching the `nginx` server, we iterate through the statically built `js/json/html` files, and replace all
instances of `\${FAKE_OIDC_ISSUER}` with the actual value `https://my-oidc-issuer`.
