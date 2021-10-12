#!/usr/bin/env bash

go run "$(go env GOROOT)"/src/crypto/tls/generate_cert.go --rsa-bits 2048 --host 127.0.0.1,::1,localhost,hydra --ca --start-date "Jan 1 00:00:00 1970" --duration=1000000h

mkdir -p certs
rm -rf certs/*
mv -- *.pem certs/
