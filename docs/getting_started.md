# Run the server

1. Install dependencies `make install-all`
2. Start docker resources `make up`
3. Create initial configuration `make bootstrap`
4. Start the server `make serve`
5. Start the frontend `make serve-pwa`
6. Open the browser https://localhost:3000
7. Credentials will be generated in `creds/`
8. **Add the generated CA Certificate to your trust store**
   1. Fedora: `sudo cp creds/ca/tls.cert /etc/pki/ca-trust/source/anchors/nrc_core_dev.pem && sudo update-ca-trust`
   2. Ubuntu: `sudo cp creds/ca/tls.cert /usr/local/share/ca-certificates/nrc_core_dev.pem && sudo update-ca-certificates`
   2. Mac: `sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain creds/ca/tls.cert`

# Component Overview

| Component | Address | Purpose | 
|-----------|---------|---------|
Core React App | https://localhost:3000 | Frontend for form management & data collection
Core App Backend| https://localhost:9000 | API for Core React App
Core Admin Backend | https://localhost:9001 | OAuth & Identity Provider Management
Core Admin React App | https://localhost:3001 | Frontend for Core Admin Backend
Core Login Server | https://localhost:9002 | Federates Identity Providers & Password Credentials
Hydra Public| https://localhost:4444 | Provides OIDC Protocol to Core Login Server
Hydra Private | https://localhost:4445 | Private management api for Hydra
Simple OIDC | https://localhost:9005 | Local OIDC provider for development purposes
