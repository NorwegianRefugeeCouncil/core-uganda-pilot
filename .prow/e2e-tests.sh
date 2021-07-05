#!/usr/bin/env bash

set -euxo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

echo ">> Installing dependencies"
apk update &&
  apk add --no-cache --virtual \
    .build-deps \
    npm \
    nodejs \
    curl \
    go \
    gcc \
    musl-dev \
    openssl \
    ca-certificates &&
  update-ca-certificates

echo ">> Installing go 1.16.5"
wget -O go.tgz https://dl.google.com/go/go1.16.5.src.tar.gz
tar -C /usr/local -xzf go.tgz
cd /usr/local/go/src/
./make.bash
export PATH="/usr/local/go/bin:$PATH"
export GOPATH=/opt/go/
export PATH=$PATH:$GOPATH/bin
apk del .build-deps

echo ">> Installing kubectl"
apk add curl
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x kubectl
mv kubectl /bin/kubectl

echo ">> Running unit tests"
cd "${SCRIPT_DIR}/../api" || exit 1
go test ./...

echo ">> Generating self-signed certificate"
openssl genrsa -out "${SCRIPT_DIR}/ca.key" 2048
openssl req -new -x509 -days 365 -key "${SCRIPT_DIR}/ca.key" -subj "/C=CN/ST=GD/L=SZ/O=Acme, Inc./CN=Acme Root CA" -out "${SCRIPT_DIR}/ca.crt"
openssl req -newkey rsa:2048 -nodes -keyout "${SCRIPT_DIR}/server.key" -subj "/C=CN/ST=GD/L=SZ/O=Acme, Inc./CN=*.example.com" -out "${SCRIPT_DIR}/server.csr"
openssl x509 -req -extfile <(printf "subjectAltName=DNS:example.com,DNS:www.example.com") -days 365 -in "${SCRIPT_DIR}/server.csr" -CA "${SCRIPT_DIR}/ca.crt" -CAkey "${SCRIPT_DIR}/ca.key" -CAcreateserial -out "${SCRIPT_DIR}/server.crt"

echo ">> Preparing environment variables"
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
export HYDRA="${ID}-hydra"
export MONGO_HOST="${ID}-mongo"
export POSTGRES_USERNAME_RAW=$(tr </dev/urandom -dc A-Za-z0-9 | head -c${1:-8})
export POSTGRES_USERNAME=$(echo -n "${POSTGRES_USERNAME_RAW}" | base64 -w0)
export POSTGRES_PASSWORD_RAW=$(tr </dev/urandom -dc A-Za-z0-9 | head -c${1:-64})
export POSTGRES_PASSWORD=$(echo -n "${POSTGRES_PASSWORD_RAW}" | base64 -w0)
export POSTGRES_DSN=$(echo -n "postgres://${POSTGRES_USERNAME_RAW}:${POSTGRES_PASSWORD_RAW}@${HYDRA_DB}:5432/hydra?sslmode=disable&max_conns=20&max_idle_conns=4" | base64 -w0)
export MONGO_PASSWORD_RAW=$(tr </dev/urandom -dc A-Za-z0-9 | head -c${1:-64})
export MONGO_PASSWORD=$(echo "${MONGO_PASSWORD_RAW}" | base64 -w0)
export MONGO_USERNAME_RAW=$(tr </dev/urandom -dc A-Za-z0-9 | head -c${1:-8})
export MONGO_USERNAME=$(echo "${MONGO_USERNAME_RAW}" | base64 -w0)
export HYDRA_SECRET=$(tr </dev/urandom -dc A-Za-z0-9 | head -c${1:-32} | base64 -w0)

echo ">> Generating k8s manifests"
envsubst <"${SCRIPT_DIR}/e2e.yaml.tpl" >"${SCRIPT_DIR}/test.yaml"

echo ">> Starting Services"
kubectl apply -f "${SCRIPT_DIR}/test.yaml" --namespace test-pods

cleanup() {
  rv=$?
  echo ">> Cleaning up"
  kubectl delete -f "${SCRIPT_DIR}/test.yaml" --namespace test-pods
  exit $rv
}
trap "error" TERM
trap "cleanup" EXIT

echo ">> Waiting for services"
sleep 60

echo ">> Running integration tests"
cd "${SCRIPT_DIR}/../api" || exit 1
go test ./... -coverprofile "${SCRIPT_DIR}/cover.out" --tags=integration

echo ">> Launching Server for e2e tests"
cd "${SCRIPT_DIR}/../api"
go run ./cmd \
  --mongo-database=e2e \
  --mongo-username=root \
  --mongo-password="${MONGO_PASSWORD_RAW}" \
  --mongo-hosts="${MONGO_HOST}:27017" \
  --environment=Development \
  --fresh=true \
  --seed=true \
  --hydra-admin-url="https://${HYDRA}:4445" \
  --hydra-public-url="https://${HYDRA}:4444" \
  --login-templates-directory=pkg/apps/login/templates \
  --login-client-id=login \
  --login-client-name=login \
  --login-client-secret=somesecret \
  --login-iam-host=localhost:9000 \
  --login-iam-scheme=https \
  --web-templates-directory=pkg/apps/webapp/templates \
  --web-client-id=webapp \
  --web-client-secret=webapp \
  --web-client-name=webapp \
  --web-iam-host=localhost:9000 \
  --web-iam-scheme=https \
  --web-cms-host=localhost:9000 \
  --web-cms-scheme=https \
  --listen-address=:9000 \
  --base-url=https://localhost:9000 \
  --tls-cert-path="${SCRIPT_DIR}/server.pem" \
  --tls-key-path="${SCRIPT_DIR}/server.key" &

echo ">> Running e2e tests"
cd "${SCRIPT_DIR}/../e2e" || exit 1
npm i
npm run run
