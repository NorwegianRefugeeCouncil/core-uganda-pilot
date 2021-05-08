#!/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

cd "${SCRIPT_DIR}/.." || exit 1

# go generate ./...

controller-gen object paths=$(pwd)/pkg/apis/meta/v1
controller-gen object paths=$(pwd)/pkg/apis/core/v1


find . -type f -name "*.go" -print0 | xargs -0 sed -i '' -e 's/k8s.io\/apimachinery\/pkg\/runtime/github.com\/nrc-no\/core\/apps\/api\/pkg\/runtime/g'


