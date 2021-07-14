---
layout: page
title: "Architecture"
permalink: /architecture/
---
Core is written in go and uses server-side rendering via [go templates](https://golang.org/pkg/text/template/).


```
├── api                             # core lives here
│   ├── certs                       
│   ├── cmd                         # the package's root function
│   │   └── app
│   └── pkg                         
│       ├── apps
│       │   ├── cms                 # content management system (cases, case types, comments)
│       │   ├── iam                 # identity and access management (persons, relationships, teams, etc)
│       │   ├── login               # login front-end
│       │   │   └── templates
│       │   ├── seeder              # DB seeding
│       │   └── webapp              # main web frontend
│       │       └── templates
│       ├── auth                    # helper functions for authentication 
│       ├── generic                 # server struct abstractions
│       │   └── server
│       ├── middleware              
│       ├── rest                    # REST client
│       ├── server                  # main amd secondary servers
│       └── sessionmanager
├── artifacts                       # docker-compose configs
├── docs
└── e2e
    └── cypress
        ├── fixtures
        ├── helpers
        ├── integration             # e2e tests
        ├── pages
        ├── plugins
        ├── screenshots
        ├── support
        └── videos
```
