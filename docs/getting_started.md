# Run the server

- Install dependencies `make install-all`
- Add these entries to your `/etc/hosts` file

```
127.0.0.1 core-admin-api.dev
127.0.0.1 core-admin-frontend.dev
127.0.0.1 core-app-api.dev
127.0.0.1 core-app-frontend.dev
127.0.0.1 core-login.dev
127.0.0.1 hydra.dev
127.0.0.1 hydra-admin.dev
127.0.0.1 oidc.dev

::1  core-admin-api.dev
::1  core-admin-frontend.dev
::1  core-app-api.dev
::1  core-app-frontend.dev
::1  core-login.dev
::1  hydra.dev
::1  hydra-admin.dev
::1  oidc.dev
```

- Start docker resources `make up`
- Create initial configuration `make bootstrap`
- Start the server `make serve`
- Start the frontend `make serve-pwa`
- Open the browser https://core-app-frontend.dev:3000
- Credentials will be generated in `creds/`
- **Add the generated CA Certificate to your trust store**
  - Fedora: `sudo cp creds/ca/tls.cert /etc/pki/ca-trust/source/anchors/nrc_core_dev.pem && sudo update-ca-trust`
  - Ubuntu: `sudo cp creds/ca/tls.cert /usr/local/share/ca-certificates/nrc_core_dev.pem && sudo update-ca-certificates`
  - Mac: `sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain creds/ca/tls.cert`

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
