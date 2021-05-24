[live showcase](https://nrc-no.github.io/core/)

# NRC Core

This project delivers the NRC Core Web Application

### Project Layout

```
├── api                                 # api server source code
│   ├── cmd                             # launch commands
│   │   └── server                      # starts the api server
│   ├── docker                          # docker files for dev/build like docker-compose.yaml
│   ├── examples                        # example files/payloads
│   │   ├── customresourcedefinitions 
│   │   └── formdefinitions
│   ├── hack                            # misc scripts for codegen and so on
│   ├── pkg                             # main source directory
│   │   ├── apis                        # api types definitions
│   │   ├── client                      # go client library for the core api server
│   │   ├── controllers                 # controllers/reconcilers
│   │   ├── customresource              # contains code for handling custom resource definitions
│   │   ├── e2e                         # e2e tests
│   │   ├── endpoints                   # http endpoints
│   │   ├── fields                      # related to query selectors
│   │   ├── generated                   # generated openapi stuff
│   │   ├── openapi                     # related to handling jsonSchema/openapi structures
│   │   ├── registry                    # REST storage registry configurations
│   │   ├── server                      # HTTP server config and startup
│   │   └── store                       # mongo storage implementation
│   └── third-party                     # third party dependencies that can't be obtained with go.mod
│       └── kubernetes                   
├── docs                                # documentation app
└── ui                                  # ui related files
    ├── apps                            # apps
    │   ├── core-web-app                # core app
    │   ├── core-web-app-e2e            # core app e2e tests (cypress)
    │   ├── showcase                    # showcase app
    │   └── showcase-e2e                # showcase app e2e tests (cypress)
    ├── dist                            # build/dist folder
    │   ├── apps
    │   └── libs
    ├── libs
    │   ├── core-ui                     # core ui components
    │   ├── formbuilder                 # form builder
    │   ├── formrenderer                # form renderer
    │   ├── openapi                     # openapi definitions
    │   └── shared                      # shared stuff
    │       ├── api-client              # typescript client for interacting with the api
    │       ├── bootstrap               # customized bootstrap theme
    │       └── ui-toolkit              # bootstrap components wrapped in react components
    └── tools                           # nwrl nx related commands
        ├── executors     
        └── generators



```

## Features

The main features of this project include,

- API Server for Custom Resource Definition & Management
- JSONSchema Forms
- PWA for Offline Data Capture

# Getting Started

The instructions here will get your development environment setup for this project.

## Prerequisites

TODO

## IDE

TODO

## Installing

TODO

## Configuration / Environment

TODO

## Building

### UI

Install dependencies,

    $ npm install

Build bootstrap theme,

    $ npm run bootstrap:build

## Testing

ToDo

## Tests

TODO

## Executing

TODO

## Code style check

Code style checks are performed with tslint settings in _tslint.json_ and the IDE preferences.

## Deployment

This project can be deployed by ....

# Versioning

We use [Semantic Versioning](https://semver.org/). For the versions available, see the [tags on this repository].

# Authors

- [Ludovic Cleroux](https://github.com/ludydoo) ludovic.cleroux@nrc.no
- [Robert Focke](https://github.com/shinroo) robert.focke@nrc.no
- [Senyao Hou](https://github.com/senyaoh) senyao.hou@nrc.no
- [Nicolas Epstein](https://github.com/nilueps) nilueps@gmail.com
