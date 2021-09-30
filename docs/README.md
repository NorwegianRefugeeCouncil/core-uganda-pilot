![nrc](https://user-images.githubusercontent.com/55964909/125176016-409c4e80-e19e-11eb-8ef1-8315dc6e3e34.jpg)

# [Core](https://nrc-no.github.io/core/)
A case management tool for the humanitarian sector - keep track of beneficiaries and services.

- [ARCHITECTURE](ARCHITECTURE.md): get to know the repo
- [GLOSSARY](GLOSSARY.md): understand commonly used terms
- [CONTRIBUTING](CONTRIBUTING.md): contribution guidelines

## Getting started:

__nb__ Use your IDE of choice, but be aware this repository includes convenience scripts for the GoLand IDE.

### Setup

#### 1. Install dependencies

(GoLand should do this automatically the first time you open the directory)
```bash
cd api
go get ./cmd
cd ../e2e
npm install
```
You will also need to install `air`: [https://github.com/cosmtrek/air](https://github.com/cosmtrek/air). Make sure `air` is available in `bash`'s `PATH

And Typescript: `npm install -g typescript`. -> make sure `tsc` is available from bash

**note**, all the following commands are run from the `api` directory, unless otherwise specified


#### 2. Generate certificates (you only have to do this once)

```bash
make gen-certs
```

#### 3. Install git hooks

Navigate to the `actions` directory
```bash
./install.sh
```

### Launch

#### 1. Start the backend services

With the `docker` daemon running in the background, the following command reads from a docker-compose config and launches mongoDB, hydra (auth) and redis.

```bash
make up
```


#### 2. Start the dev server
 
 
 ```bash
 # to launch the server
 make serve
 
 # to launch the server in watch (hot reload) mode
 make watch
 ```
 
 These are also available as run configs from Goland

### Endpoints:
 - webapp: [http://localhost:9000](http://localhost:9000)
 - mongo express: [http://localhost:8090](http://localhost:8090)
 

### Login:

To log into the webapp, you may use the following credentials:
- username: __courtney.lare@email.com__
- password: __password__


## Running tests

With the git hook installed, tests will run locally automatically before every push. If any tests fail, the push operation will be aborted.

You can run tests manually: with `make test-iam`, `make test-cms`, `make test-attachments` and `make test-e2e`

Or via the run configurations available in Goland.

