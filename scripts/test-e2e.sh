#!/usr/bin/env bash

cd "$(dirname "$0")" || exit

cd ../test/e2e || exit

npm i && npm run run
