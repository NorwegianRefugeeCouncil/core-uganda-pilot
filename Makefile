SHELL := bash
.ONESHELL: # ensures each Make recipe is ran as one single shell session, rather than one new shell per line
.SHELLFLAGS := -eu -o pipefail -c # fail on errors

.PHONY: up
up:
	@./scripts/up.sh

.PHONY: down
down:
	@./scripts/down.sh

.PHONY: migrate
migrate:
	@./scripts/migrate.sh

.PHONY: clear-db
clear-db:
	@./scripts/clear-db.sh

.PHONY: reset-db
reset-db: clear-db migrate

.PHONY: build
build:
	@./scripts/build.sh

.PHONY: serve
serve:
	@./scripts/serve.sh

.PHONY: frontend
frontend:
	@./scripts/frontend.sh

.PHONY: watch
watch:
	@./scripts/watch.sh

.PHONY: package
package:
	@./scripts/package.sh

.PHONY: serve-pwa
serve-pwa:
	@./scripts/serve-pwa.sh

.PHONY: docs
docs:
	@./scripts/render-dot-graphs.sh

.PHONY: open-all
open-all:
	@open http://localhost:3000
	@open http://localhost:9000
	@open http://localhost:9001

