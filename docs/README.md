# core
A case management tool for the humanitarian sector - keep track of beneficiaries and services.

- [ARCHITECTURE](./ARCHITECTURE.md): get to know the repo
- [CONTRIBUTING](./CONTRIBUTING.md): contribution guidelines

## Getting started:

__nb__ Use your IDE of choice, but be aware this repository includes convenience scripts for the GoLand IDE.

### 1. Install dependencies

(GoLand should do this automatically the first time you open the directory)
```bash
cd api
go get ./cmd
cd ../e2e
npm install
```

### 2. Start the backend services
```bash
cd api
make up
```

### 3. Generate certificates (you only have to do this once)
```bash
make gen-certs
```

### 4. Start the development server
In Goland, choose the `start` configuration from the top-right menu and press the green arrow.

or
 
 ```bash
 make serve
 ```

### 5. Endpoints:
 - webapp: [https://localhost:9000](http://localhost:9000)
 - mongo dashboard: [https://localhost:8090](http://localhost:8090)
 
 
## Running tests
 
### Integration tests
 
In Goland, choose the `iam test` configuration from the top-right menu and press the green arrow.

or

```
cd api
make test
```

### E2E tests


In Goland, choose the `e2e (browser)` configuration from the top-right menu and press the green arrow.

```bash
cd e2e
npm run open
```
=======
## Architecture

```
├── api                             # Main Core go app
│   ├── artifacts                   # docker files
│   ├── cmd                         # main go package definition
│   │   └── app
│   └── pkg                         
│       ├── apps                    # app modules    
│       │   ├── cms             
│       │   ├── iam
│       │   ├── login
│       │   │   └── templates
│       │   ├── seed
│       │   └── webapp              # client-facing web app
│       │       └── templates
│       ├── auth
│       ├── middleware
│       ├── rest
│       ├── server                  # main server app
│       ├── sessionmanager
│       ├── testhelpers             # unused for now
│       └── testing                 # backend unit tests
└── e2e                 
    └── cypress
        ├── fixtures
        ├── integration             # e2e tests
        ├── plugins
        └── support
```

## Getting started:

1. Start the backend services:
```bash
cd core/api

# then

docker-compose -f artifacts/docker-compose.yaml up -d

# or

make spinup
```


2. Start the development server (3 options):
 - Within GoLand IDE, open the `api` project directory as a root workspace (ie. not the `core` directory). Choose the `start` configuration from the top-right menu and press the green arrow.
 
 ```bash

 cd api
 
 # Generate certificate + key for development purpose
 make gen-certs

 go run ./cmd \
  --mongo-database=core \
  --mongo-username=root \
  --mongo-password=example \
  --mongo-hosts=localhost:27017 \
  --environment=Development \
  --fresh=true \
  --seed=true \
  --hydra-admin-url=https://localhost:4445 \
  --hydra-public-url=https://localhost:4444 \
  --login-templates-directory=pkg/apps/login/templates \
  --login-client-id=login \
  --login-client-name=login \
  --login-client-secret=somesecret \
  --login-iam-host=localhost:9000 \
  --login-iam-scheme=https \
  --web-templates-directory=pkg/apps/webapp/templates \
  --web-client-id=webapp \
  --web-client-secret=webapp \
  --web-client-name=webapp \
  --web-iam-host=localhost:9000 \
  --web-iam-scheme=https \
  --web-cms-host=localhost:9000 \
  --web-cms-scheme=https \
  --listen-address=:9000 \
  --base-url=https://localhost:9000 \
  --tls-cert-path=certs/cert.pem \
  --tls-key-path=certs/key.pem
 
 # or
 
 make serve
 ```
 
 
3. Endpoints:
 - webapp: [http://localhost:9000](http://localhost:9000)
 - mongo express: [http://localhost:8090](http://localhost:8090)
 
 
 ## Running tests
 
 ### Unit tests
 
 Goland: open the `api` workspace, navigate to the `testing folder` open any file and click on the two green arrows in the gutter at the top of the editor.

  or
  
 ```
 cd core/api
 make test
 ```
 
 ### E2E tests
 
 ```bash
 cd core/e2e
 npm install
 npm run open
 ```
