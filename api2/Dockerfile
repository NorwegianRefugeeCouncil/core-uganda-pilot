FROM golang:1.16.5-alpine

RUN mkdir /app
COPY . /app/
WORKDIR /app

ENTRYPOINT go run ./cmd --mongo-database=core --mongo-username=root --mongo-password=example --keycloak-base-url=http://localhost:8080 --keycloak-realm-name=nrc --keycloak-client-id=api --keycloak-client-secret=e6486272-039d-430f-b3c7-47887aa9e206
