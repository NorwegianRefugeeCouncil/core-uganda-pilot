# Run the server

- Install dependencies `make install-all`
- Add these entries to your `/etc/hosts` file

```
127.0.0.1 core-admin-api.dev core-admin-frontend.dev core-app-api.dev core-app-frontend.dev core-login.dev hydra.dev hydra-admin.dev oidc.dev
::1       core-admin-api.dev core-admin-frontend.dev core-app-api.dev core-app-frontend.dev core-login.dev hydra.dev hydra-admin.dev oidc.dev
```

- Create your secrets `make init-secrets`
- **Add the generated CA Certificate to your trust store**
    - Fedora: `sudo cp creds/ca/tls.cert /etc/pki/ca-trust/source/anchors/nrc_core_dev.pem && sudo update-ca-trust`
    - Ubuntu: `sudo cp creds/ca/tls.cert /usr/local/share/ca-certificates/nrc_core_dev.pem && sudo update-ca-certificates`
    - Mac: `sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain creds/ca/tls.cert`
- Start docker resources `make up`
- Migrate the database `make migrate`
- Create initial configuration `make bootstrap`
- Start the server `make serve`
- Start the frontend `make serve-pwa`
- Start the admin frontend `make serve-admin`
- Open the browser https://core-app-frontend.dev:3000
	- Authenticate with \<whatever\>`@nrc.no` (any email ending with `nrc.no`)
	- Click `Login with Norwegian Refugee Council`
	- Put any password (there is no password verification for this development oidc-provider)
- Open the browser https://core-admin-frontend.dev:3001
- Credentials will be generated in `creds/`

## What is my user/password?

You can put any password you want in the fake oidc provider.

# Component Overview

| Component | Address | Purpose | 
|-----------|---------|---------|
Core React App | https://core-app-frontend.dev:3000 | Frontend for form management & data collection
Core App Backend| https://core-app-api.dev:9000 | API for Core React App
Core Admin React App | https://core-admin-frontend.dev:3001 | Frontend for Core Admin Backend
Core Admin Backend | https://core-admin-api.dev:9001 | OAuth & Identity Provider Management
Core Login Server | https://core-login.dev:9002 | Federates Identity Providers & Password Credentials
Hydra Public| https://hydra.dev:4444 | Provides OIDC Protocol to Core Login Server
Hydra Private | https://hydra-admin.dev:4445 | Private management api for Hydra
Simple OIDC | https://oidc.dev:9005 | Local OIDC provider for development purposes
