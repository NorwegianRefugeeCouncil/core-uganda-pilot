SHELL := bash
.ONESHELL: # ensures each Make recipe is ran as one single shell session, rather than one new shell per line
.SHELLFLAGS := -eu -o pipefail -c # fail on errors

scripts = ./scripts
flags = --mongo-database=core \
	--mongo-username=root \
	--mongo-password=example \
	--mongo-hosts=localhost:27017 \
	--redis-address=localhost:6379 \
	--environment=Development \
	--fresh=true \
	--seed=true \
	--tls-disable=true \
	--hydra-admin-url=http://localhost:4445 \
	--hydra-public-url=http://localhost:4444 \
	--login-templates-directory=web/templates/login \
	--login-client-id=login \
	--login-client-name=login \
	--login-client-secret=somesecret \
	--login-iam-host=localhost:9000 \
	--login-iam-scheme=http \
	--web-templates-directory=web/templates \
	--web-static-directory=web/static \
	--web-client-id=webapp \
	--web-client-secret=webapp \
	--web-client-name=webapp \
	--web-iam-host=localhost:9000 \
	--web-iam-scheme=http \
	--web-cms-host=localhost:9000 \
	--web-cms-scheme=http \
	--listen-address=:9000 \
	--base-url=http://localhost:9000

up:
	@docker-compose -f deployments/docker-compose.yaml up --build -d

down:
	@docker-compose -f deployments/docker-compose.yaml down --remove-orphans

gen-certs:
	@$(scripts)/gen-certs.sh

test-integration:
	@go test github.com/nrc-no/core/pkg/iam github.com/nrc-no/core/pkg/cms github.com/nrc-no/core/pkg/attachments

test-e2e:
	@cd test/e2e && npm i && npm run run

build:
	@tsc && go build -o ./tmp/main ./cmd/core

serve: build
	@./tmp/main $(flags)

watch:
	@air -c configs/.air.toml

docker-build:
	@docker build -t digitalbackbonecr.azurecr.io/nrc-no/core:latest -f build/package/Dockerfile .

docker-push: docker-build
	@docker push digitalbackbonecr.azurecr.io/nrc-no/core:latest

.PHONY: up down serve build
