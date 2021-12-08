```
├── build                  # Dockerfiles
├── cmd                    # startup commands
│     ├── clear.go         # clears the database
│     ├── migrate.go       # migrate the database
│     ├── root.go          # root command
│     ├── serve.go         # serve base command
│     ├── serve_admin.go   # serve admin server command
│     ├── serve_login.go   # serve login server command
│     └── serve_public.go  # serve public server command
├── configs                # configuration files
├── deployments            # docker-compose files
├── docs                   # documentation
├── githooks               # git hooks
├── pkg                    # code
│     ├── api              # common api types
│     │     ├── meta       # api meta types
│     │     └── types      # api types
│     ├── client           # go client
│     ├── constants        # shared constants 
│     ├── rest             # HTTP commands
│     ├── server           
│     │     ├── admin      # admin server
│     │     ├── generic    # generic server
│     │     ├── login      # login server
│     │     ├── options    # server options
│     │     └── public     # public server
│     ├── sql
│     │     ├── convert    # sql conversion
│     │     └── schema     # sql schema types
│     ├── store            # store interface
│     ├── utils            # utilities
│     │     ├── pointers   # pointer utilities
│     │     └── sets       # set utilities
│     └── validation       # validation utilities
├── scripts                # various scripts
└── frontend               # frontend lerna monorepo 
    ├── apps               # frontend applications
    |   ├── native         # core react native application
    |   ├── pwa            # core pwa
    |   ├── admin          # identity provider admin application 
    └── packages           # frontend shared packages
        ├── core-api-client# core api client
        ├── core-auth      # core authentication react component
        ├── design-system  # core react native desgin system

```
