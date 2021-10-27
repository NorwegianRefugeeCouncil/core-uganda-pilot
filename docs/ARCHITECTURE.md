---
layout: page title: "Architecture"
permalink: /architecture/
---
Core is written in go and uses server-side rendering via [go templates](https://golang.org/pkg/text/template/).

```
├── build
│     └── package
├── cmd
├── configs
├── deployments
├── docs
├── githooks
├── pkg
│     ├── api
│     │     └── meta
│     ├── client
│     ├── constants
│     ├── options
│     ├── pointers
│     ├── rest
│     ├── server
│     │     ├── admin
│     │     ├── generic
│     │     ├── login
│     │     └── public
│     ├── sets
│     ├── sqlconvert
│     ├── sqlschema
│     ├── storage
│     ├── store
│     ├── types
│     ├── utils
│     └── validation
├── scripts
└── web
    ├── app
    │     ├── client
    │     └── frontend
    └── pwa
        ├── node_modules
        ├── public
        └── src

```

# Data flow

![](overview.png)

1. The browser requests a webpage from `webapp`
1. If the user is not authenticated, the user is redirected to the `login` app
	1. the user provides credentials
	1. `login` app performs login flow with hydra. Performs oAuth redirect to `webapp`
	1. `webapp` stores session information in `redis`, populates session cookie
1. `webapp` verifies the authenticity of the access-token present in the session by calling `hydra`
1. `webapp` makes HTTP requests to any of `iam`, `cms` or `attachments` apps, forwarding the access-token using
   the `Authorization` header
1. `iam`, `cms`, `attachments` also verify the `Authorization` header passed by `webapp`
1. `iam`, `cms`, `attachments`, `login` can put/get data from mongo
1. `attachments` stores blob files in `s3` (photos, videos, pdfs, etc.)
