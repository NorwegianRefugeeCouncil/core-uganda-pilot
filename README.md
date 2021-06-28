# core
A case management tool for the humanitarian sector - keep track of beneficiaries and services.

## Architecture

```
├── api                             # Main Core go app
│   ├── artifacts                   # docker files
│   ├── cmd                         # main go package definition
│   │   └── app
│   └── pkg                         
│       ├── apps                    # app modules    
│       │   ├── cms             
│       │   ├── iam
│       │   ├── login
│       │   │   └── templates
│       │   ├── seed
│       │   └── webapp              # client-facing web app
│       │       └── templates
│       ├── auth
│       ├── middleware
│       ├── rest
│       ├── server                  # main server app
│       ├── sessionmanager
│       ├── testhelpers             # unused for now
│       └── testing                 # backend unit tests
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
 go run ./cmd --mongo-database=core --mongo-username=root --mongo-password=example --keycloak-base-url=http://localhost:8080 --keycloak-realm-name=nrc --keycloak-client-id=api --keycloak-client-secret=e6486272-039d-430f-b3c7-47887aa9e206
 
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
