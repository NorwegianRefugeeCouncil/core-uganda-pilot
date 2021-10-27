SHELL := bash
.ONESHELL: # ensures each Make recipe is ran as one single shell session, rather than one new shell per line
.SHELLFLAGS := -eu -o pipefail -c # fail on errors

up:
	@./scripts/up.sh

down:
	@./scripts/down.sh

migrate:
	@./scripts/migrate.sh

build:
	@./scripts/build.sh

serve:
	@./scripts/serve.sh

frontend:
	@./scripts/frontend.sh

watch:
	@./scripts/watch.sh

package:
	@./scripts/package.sh

pwa:
	@./scripts/pwa.sh

.PHONY: up down serve build test frontend
