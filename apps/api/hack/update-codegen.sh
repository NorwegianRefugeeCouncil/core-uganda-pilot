#!/bin/bash

ROOT=$(dirname "$(dirname "${BASH_SOURCE[0]}")/../")

cd "${ROOT}/.."

controller-gen object paths="./pkg/apis/..."

find . -type f -name "zz_generated.deepcopy.go" -print0 | xargs -0 sed -i '' -e 's/\"k8s.io\/apimachinery\/pkg\/runtime\"/\"github.com\/nrc-no\/core\/apps\/api\/pkg\/runtime\"/g'

conversion-gen -v 5 \
  --go-header-file "./hack/boilerplate.go.txt" \
  --input-dirs "./pkg/apis/meta/v1/,./pkg/apis/example/v1/,./pkg/apis/example/v2/,./pkg/apis/example/,./pkg/apis/core/,./pkg/apis/core/v1/,./pkg/apis/core/internalversion/" \
  --output-base "." \
  --output-file-base="zz_generated.conversion"

find . -type f -name "zz_generated.conversion.go" -print0 | xargs -0 sed -i '' -e 's/\"k8s.io\/apimachinery\/pkg\/runtime\"/\"github.com\/nrc-no\/core\/apps\/api\/pkg\/runtime\"/g'
find . -type f -name "zz_generated.conversion.go" -print0 | xargs -0 sed -i '' -e 's/\"k8s.io\/apimachinery\/pkg\/conversion\"/\"github.com\/nrc-no\/core\/apps\/api\/pkg\/conversion\"/g'

array=()
while IFS= read -r -d ''; do
  array+=("$REPLY")
done < <(find . -type f -name "zz_generated.conversion.go" -print0)

for item in "${array[@]}"; do
  echo "$item"
  package=$(basename $(dirname "${item}"))
  expr="s/[[:space:]]+${package} \".*\"//g"
  sed -i -E "${expr}" "${item}"
  sed -i -E "s/${package}.//g" "${item}"
done

