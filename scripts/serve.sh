#!/usr/bin/env bash

function cleanup() {
	# kill all child processes
	pkill -P $$
}
# cleanup on exit
for sig in INT QUIT HUP TERM; do
	trap "
	cleanup
	trap - $sig EXIT
	kill -s $sig "'"$$"' "$sig"
done
trap cleanup EXIT

flags="--mongo-database=core \
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
       	--base-url=http://localhost:9000"

if [ ! -f tmp/main ]; then
  ./scripts/build.sh
fi

# shellcheck disable=SC2086
./tmp/main $flags
