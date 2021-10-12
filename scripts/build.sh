#!/usr/bin/env bash

if ! tsc; then
  exit 1
fi

go build -o ./tmp/main ./cmd/core
