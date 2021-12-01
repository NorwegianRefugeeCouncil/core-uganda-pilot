# Run the server

- Install Go
  - `brew install go`
- Install `docker` and `docker-compose`
- Install dependencies `make install-all`
- Create your secrets `make init-secrets`
- Start tunnels `make tunnels` (install localtunnels if you don't have it `yarn global add localtunnel`)
    - When you start the tunnels for the first time, you will see 2 URLs that correspond to your tunnel URI. Open the
      links, and approve the dialog that will be presented to you. Otherwise all calls through the tunnels will fail.
      This has to be done every ~1 week as they are reset.
- **Add the generated CA Certificate to your trust store**
    - Fedora: `sudo cp certs/ca.crt /etc/pki/ca-trust/source/anchors/nrc_core_dev.crt && sudo update-ca-trust`
    - Ubuntu: `sudo cp certs/ca.crt /usr/local/share/ca-certificates/nrc_core_dev.crt && sudo update-ca-certificates`
    - Mac: `sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain certs/ca.crt`
    - If you are using **Chrome on Linux**:
        - settings
        - privacy and security
        - manage certificates
        - trust anchors
        - select certs/ca.crt
        - trust for identifying websites
        - Do the same sort of thing in firefox if you use that
- Start docker resources `make up`
- Migrate the database `make migrate`
- Create initial configuration `make bootstrap`
- Start the server `make serve`
- Start the frontend `make serve-pwa`
- Start the admin frontend `make serve-admin`
- Open the browser http://localhost:3000
    - Authenticate with \<whatever\>`@nrc.no` (any email ending with `nrc.no`)
    - Click `Login with Norwegian Refugee Council`
    - Put any password (there is no password verification for this development oidc-provider)
- Open the browser http://localhost:3001
- Credentials will be generated in `creds/`

## What is my user/password?

You can put any password you want in the fake oidc provider.

# Component Overview

| Component | Address | Purpose | 
|-----------|---------|---------|
Core React App | http://localhost:3000 | Frontend for form management & data collection
Core App Backend| http://localhost:9000 | API for Core React App
Core Admin React App | http://localhost:3001 | Frontend for Core Admin Backend
Core Admin Backend | http://localhost:9001 | OAuth & Identity Provider Management
Core Login Server | http://localhost:9002 | Federates Identity Providers & Password Credentials
Hydra Public| http://localhost:4444 | Provides OIDC Protocol to Core Login Server
Hydra Private | http://localhost:4445 | Private management api for Hydra
Simple OIDC | http://localhost:9005 | Local OIDC provider for development purposes
