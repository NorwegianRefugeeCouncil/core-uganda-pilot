#!/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

cd "${SCRIPT_DIR}/.." || exit 1

# go generate ./...

controller-gen object paths=$(pwd)/...
conversion-gen --go-header-file "./hack/boilerplate.go.txt" --input-dirs "./testing/testscheme" 	--output-base "." --output-file-base="zz_generated.conversion" --skip-unsafe=true

find . -type f -name "*.go" -print0 | xargs -0 sed -i '' -e 's/k8s.io\/apimachinery\/pkg\/runtime/github.com\/nrc-no\/core\/apps\/api\/pkg\/runtime/g'
find . -type f -name "*.go" -print0 | xargs -0 sed -i '' -e 's/k8s.io\/apimachinery\/pkg\/conversion/github.com\/nrc-no\/core\/apps\/api\/pkg\/conversion/g'



