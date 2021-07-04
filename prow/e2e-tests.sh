#!/usr/bin/env bash

#!/bin/bash
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

echo RUNNING UNIT TESTS
cd "${SCRIPT_DIR}/../api" || exit 1
go test ./...

echo GENERATING SELF SIGNED CERTIFICATE

openssl genrsa -out "${SCRIPT_DIR}/ca.key" 2048
openssl req -new -x509 -days 365 -key "${SCRIPT_DIR}/ca.key" -subj "/C=CN/ST=GD/L=SZ/O=Acme, Inc./CN=Acme Root CA" -out "${SCRIPT_DIR}/ca.crt"
openssl req -newkey rsa:2048 -nodes -keyout "${SCRIPT_DIR}/server.key" -subj "/C=CN/ST=GD/L=SZ/O=Acme, Inc./CN=*.example.com" -out "${SCRIPT_DIR}/server.csr"
openssl x509 -req -extfile <(printf "subjectAltName=DNS:example.com,DNS:www.example.com") -days 365 -in "${SCRIPT_DIR}/server.csr" -CA "${SCRIPT_DIR}/ca.crt" -CAkey "${SCRIPT_DIR}/ca.key" -CAcreateserial -out "${SCRIPT_DIR}/server.crt"

token=$(cat /var/run/secrets/kubernetes.io/serviceaccount)

export BUILD_ID=123
export JOB_NAME=e2e
export REPO_OWNER=nrc-no
export REPO_NAME=core
export PROW_JOB_ID=prow123
export PULL_BASE_REF=ref123
export PULL_BASE_SHA=sha123
export TLS_KEY=$(cat "${SCRIPT_DIR}/server.key" | base64 -w0)
export TLS_CERT=$(cat "${SCRIPT_DIR}/server.crt" | base64 -w0)
export ID=${REPO_OWNER}-${REPO_NAME}-${JOB_NAME}-$(tr </dev/urandom -dc a-z0-9 | head -c${1:-8})

export HYDRA_DB="${ID}-hydra-db"
export MONGO_HOST="${ID}-mongo"

export POSTGRES_USERNAME_RAW=$(tr </dev/urandom -dc A-Za-z0-9 | head -c${1:-8})
export POSTGRES_USERNAME=$(echo -n "${POSTGRES_USERNAME_RAW}" | base64 -w0)
export POSTGRES_PASSWORD_RAW=$(tr </dev/urandom -dc A-Za-z0-9 | head -c${1:-64})
export POSTGRES_PASSWORD=$(echo -n "${POSTGRES_PASSWORD_RAW}" | base64 -w0)
export POSTGRES_DSN=$(echo -n "postgres://${POSTGRES_USERNAME_RAW}:${POSTGRES_PASSWORD_RAW}@${HYDRA_DB}:5432/hydra?sslmode=disable&max_conns=20&max_idle_conns=4" | base64 -w0)
export MONGO_PASSWORD=$(tr </dev/urandom -dc A-Za-z0-9 | head -c${1:-64} | base64 -w0)
export MONGO_USERNAME=$(tr </dev/urandom -dc A-Za-z0-9 | head -c${1:-8} | base64 -w0)
export HYDRA_SECRET=$(tr </dev/urandom -dc A-Za-z0-9 | head -c${1:-32} | base64 -w0)

echo GENERATING KUBERNETES MANIFESTS

envsubst <"${SCRIPT_DIR}/e2e.yaml" >"${SCRIPT_DIR}/test.yaml"

kubectl apply -f "${SCRIPT_DIR}/test.yaml" --namespace test-pods
sleep 10

# echo RUNNING INTEGRATION TESTS
cd "${SCRIPT_DIR}/../api" || exit 1
go test ./... -coverprofile "${SCRIPT_DIR}/cover.out" --tags=integration

# echo RUNNING E2E TESTS
# cd "${SCRIPT_DIR}/../e2e" || exit 1 && npm i && npm run run

kubectl delete -f "${SCRIPT_DIR}/test.yaml" --namespace test-pods

