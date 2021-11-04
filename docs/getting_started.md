# Run the server

1. Install dependencies<br>`make install-all`
2. Start docker resources && bootstrap<br>`make up`
3. Start the server<br>`make serve`
4. Start the frontend <br>`make serve-pwa`
5. Open the browser<br>http://localhost:3000
6. Credentials will be generated in
   <br>`creds/`
7. Configuration/Env files will be templated from credentials in `/creds` such as
   <br>`deployments/hydra.custom.yaml`
   <br>`deployments/postgres.env`
   <br>`deployments/redis.env`
   <br>`deployments/oidc.config.json`
   <br>`deployments/oidc.users.json`
   <br>`configs/config.custom.yaml`

# Component Overview

| Component | Address | Purpose | 
|-----------|---------|---------|
Core React App | http://localhost:3000 | Frontend for form management & data collection
Core App Backend| http://localhost:9000 | API for Core React App
Core Admin Backend | http://localhost:9001 | OAuth & Identity Provider Management
Core Admin React App | http://localhost:3001 | Frontend for Core Admin Backend
Core Login Server | http://localhost:9002 | Federates Identity Providers & Password Credentials
Hydra Public| http://localhost:4444 | Provides OIDC Protocol to Core Login Server
Hydra Private | http://localhost:4445 | Private management api for Hydra
Simple OIDC | http://localhost:9005 | Local OIDC provider for development purposes
