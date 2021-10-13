SHELL := bash
.ONESHELL: # ensures each Make recipe is ran as one single shell session, rather than one new shell per line
.SHELLFLAGS := -eu -o pipefail -c # fail on errors

up:
	@docker-compose -f deployments/docker-compose.yaml up --build -d

down:
	@docker-compose -f deployments/docker-compose.yaml down --remove-orphans

up-all:
	@docker-compose -f deployments/docker-compose-all.yaml up --build -d

down-all:
	@docker-compose -f deployments/docker-compose-all.yaml down --remove-orphans

gen-certs:
	@./scripts/gen-certs.sh

test-integration:
	@./scripts/test-integration.sh

test-e2e:
	@./scripts/test-e2e.sh

test: test-integration test-e2e

build:
	@./scripts/build.sh

serve:
	@./scripts/serve.sh

watch:
	@./scripts/watch.sh

docker-build:
	@docker build -t digitalbackbonecr.azurecr.io/nrc-no/core:latest -f build/package/Dockerfile .

docker-push: docker-build
	@docker push digitalbackbonecr.azurecr.io/nrc-no/core:latest

.PHONY: up down serve build test
