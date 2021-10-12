#!/usr/bin/env bash

tsc --build configs/tsconfig.json

if [ ! "$?" ]; then
  exit 1
fi

echo ">>> Transpiled .ts"

go build -o ./tmp/main ./cmd/core

if [ ! "$?" ]; then
  exit 1
fi

echo ">>> Built core"
