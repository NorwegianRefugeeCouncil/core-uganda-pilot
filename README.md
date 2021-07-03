# core
A case management tool for the humanitarian sector - keep track of beneficiaries and services.

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
 go run ./cmd \
  --fresh=true \
  --seed=true \
  --mongo-database=core \
  --mongo-username=root \
  --mongo-password=example \
  --mongo-hosts=localhost:27017 \
  --environment=Development \
  --hydra-admin-url=http://localhost:4445 \
  --hydra-public-url=http://localhost:4444 \
  --login-client-id=login \
  --login-client-name=login \
  --login-client-secret=somesecret \
  --login-iam-host=localhost:9000 \
  --login-iam-scheme=http \
  --web-client-id=webapp \
  --web-client-secret=webapp \
  --web-client-name=webapp \
  --web-iam-host=localhost:9000 \
  --web-iam-scheme=http \
  --web-cms-host=localhost:9000 \
  --web-cms-scheme=http \
  --listen-address=:9000 \
  --base-url=http://localhost:9000
 
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
