SHELL=bash
.ONESHELL: # ensures each Make recipe is ran as one single shell session, rather than one new shell per line
.SHELLFLAGS := -eu -o pipefail -c # fail on errors

args = `arg="$(filter-out $@,$(MAKECMDGOALS))" && echo $${arg:-${1}}`

%:
    @:

.PHONY: init-secrets
init-secrets:
	@./scripts/init_secrets.sh

.PHONY: up
up: init-secrets
	@./scripts/up.sh

.PHONY: migrate
migrate:
	@./scripts/migrate.sh

.PHONY: bootstrap
bootstrap: migrate
	@./scripts/bootstrap.sh

.PHONY: down
down:
	@./scripts/down.sh

.PHONY: clear-db
clear-db:
	@./scripts/clear-db.sh

.PHONY: reset-db
reset-db: clear-db migrate


.PHONY: build-frontend 
build-frontend:
	@./scripts/build-frontend.sh

.PHONY: prepare-frontend 
prepare-frontend:
	@./scripts/prepare-frontend.sh

.PHONY: serve-native
serve-native:
	@./scripts/serve-native.sh

.PHONY: serve-pwa
serve-pwa:
	@./scripts/serve-pwa.sh

.PHONY: serve-admin
serve-admin:
	@./scripts/serve-admin.sh

.PHONY: serve-storybook
serve-storybook:
	@./scripts/serve-storybook.sh

.PHONY: serve
serve:
	@./scripts/serve.sh

.PHONY: watch
watch:
	@./scripts/watch.sh

.PHONY: package
package:
	@./scripts/package.sh

.PHONY: docs
docs:
	@./scripts/render-dot-graphs.sh

.PHONY: install-all
install-all:
	@./scripts/install-all.sh

.PHONY: restart-all
restart-all:
	@./scripts/restart.sh

.PHONY: restart-proxy
restart-proxy:
	@./scripts/restart.sh proxy

.PHONY: restart-oidc
restart-oidc:
	@./scripts/restart.sh oidc

.PHONY: restart-redis
restart-redis:
	@./scripts/restart.sh redis

.PHONY: restart-hydra
restart-hydra:
	@./scripts/restart.sh hydra

.PHONY: restart-db
restart-db:
	@./scripts/restart.sh db

.PHONY: logs-all
logs-all:
	@./scripts/logs.sh

.PHONY: logs-proxy
logs-proxy:
	@./scripts/logs.sh proxy

.PHONY: logs-oidc
logs-oidc:
	@./scripts/logs.sh oidc

.PHONY: logs-redis
logs-redis:
	@./scripts/logs.sh redis

.PHONY: logs-hydra
logs-hydra:
	@./scripts/logs.sh hydra

.PHONY: logs-db
logs-db:
	@./scripts/logs.sh db

.PHONY: spawn
spawn: down up _sleep bootstrap build

.PHONY: tunnels
tunnels:
	@./scripts/tunnels.sh

.PHONY: proxy-local
proxy-local:
	@./scripts/proxy-local.sh

_sleep:
	sleep 20

.PHONY: respawn
respawn: down spawn

.PHONY: gen
gen:
	@./scripts/gen.sh
