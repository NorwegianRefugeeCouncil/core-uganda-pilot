#!/usr/bin/env bash

tsc --build tsconfig.json

if [ ! "$?" ]; then
  exit 1
fi

echo ">>> Transpiled typescript"

go build -o ./tmp/main ./cmd/core

if [ ! "$?" ]; then
  exit 1
fi

echo ">>> Built core"
