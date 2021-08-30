---
layout: page
title: "Architecture"
permalink: /architecture/
---
Core is written in go and uses server-side rendering via [go templates](https://golang.org/pkg/text/template/).


```
├── api                             # core lives here
|   ├── Makefile                    # used to store shorthand commands for build/run/test
|   ├── tsconfig.json               
│   ├── certs                       
│   ├── cmd                         # this is the app's entry point
│   │   └── app
├── actions                         
│   ├── git-hooks                   # contains the scripts that run on commit/push
│   └── pkg                         
│       ├── apps
│       │   ├── attachments         # deals with objects that have been 'attached' to individuals/cases
│       │   ├── cms                 # case management system (cases, case types, comments)
│       │   ├── iam                 # identity and access management (persons, relationships, teams, etc)
│       │   ├── login               # login front-end
│       │   ├── seeder              # static hard-coded data (users, countries, casetypes, etc)
│       │   └── webapp              # main web frontend
│       │       └── static          # typescript files
│       ├── auth                    # helper functions for authentication 
│       ├── form                    # front-end form element abstractions
│       ├── generic                 # server struct abstractions
│       │   ├── pagination         
│       │   └── server              # generalizes server structs and constructors used in apps
│       ├── middleware              
│       ├── registrationctrl        # a basic workflow engine that handles registration steps
│       ├── rest                    # REST client
│       ├── server                  # main server
│       ├── sessionmanager          # redis session manager
│       ├── teamstatusctrl          # controller for managing team statuses
│       ├── testutils
│       ├── utils
│       └── validation              # validation tools and abstractions

├── artifacts                       # docker-compose configs
├── docs                            # dev documentation
└── e2e
    ├── cypress.json                # cypress config
    └── cypress
        ├── fixtures                # common data (ex. json)
        ├── helpers                 # reusable constants, functions
        ├── integration             # tests
        ├── pages                   # wrappers over page-specific cypress operations
        ├── support                 # setup/teardown type operations
```
