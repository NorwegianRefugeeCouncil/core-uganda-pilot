![nrc](https://user-images.githubusercontent.com/55964909/125176016-409c4e80-e19e-11eb-8ef1-8315dc6e3e34.jpg)

# [Core](https://nrc-no.github.io/core/)
A case management tool for the humanitarian sector - keep track of beneficiaries and services.

- [ARCHITECTURE](ARCHITECTURE.md): get to know the repo
- [CONTRIBUTING](CONTRIBUTING.md): contribution guidelines

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
 - mongo dashboard: [http://localhost:8090](http://localhost:8090)
 
 
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
