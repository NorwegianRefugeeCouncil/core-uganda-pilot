#!/bin/bash

ROOT=$(dirname "$(dirname "${BASH_SOURCE[0]}")/../")

cd "${ROOT}/.."

## DEEPCOPY GEN

controller-gen object paths="./pkg/..."
defaulter-gen --go-header-file "./hack/boilerplate.go.txt" --input-dirs "./pkg/apis/apiextensions,./pkg/apis/apiextensions/v1" --output-base "." --output-file-base="zz_generated.defaults"

array=()
while IFS= read -r -d ''; do
  array+=("$REPLY")
done < <(find . -type f -name "zz_generated.deepcopy.go" -print0)

while IFS= read -r -d ''; do
  array+=("$REPLY")
done < <(find . -type f -name "zz_generated.defaults.go" -print0)

for item in "${array[@]}"; do
  package=$(basename $(dirname "${item}"))
  expr2="s/[[:space:]]+pkg${package} \".*\"//g"
  expr="s/[[:space:]]+${package} \".*\"//g"

  sed -i -e 's/\"k8s.io\/apimachinery\/pkg\/runtime\"/\"github.com\/nrc-no\/core\/apps\/api\/pkg\/runtime\"/g' "${item}"

  sed -i -E "${expr2}" "${item}"
  sed -i -E "s/pkg${package}.//g" "${item}"
  sed -i -E "${expr}" "${item}"
  sed -i -E "s/${package}.//g" "${item}"
done

#### CONVERSION GEN

conversion-gen -v 4 \
  --go-header-file "./hack/boilerplate.go.txt" \
  --input-dirs "./pkg/apis/runtime/,./pkg/apis/meta/v1/,./pkg/apis/meta/internalversion/,./pkg/apis/example/v1/,./pkg/apis/example/v2/,./pkg/apis/example/,./pkg/apis/core/,./pkg/apis/core/v1/,./pkg/apis/apiextensions,./pkg/apis/apiextensions/v1" \
  --output-base "." \
  --output-file-base="zz_generated.conversion"

find . -type f -name "zz_generated.conversion.go" -print0 | xargs -0 sed -i '' -e 's/\"k8s.io\/apimachinery\/pkg\/runtime\"/\"github.com\/nrc-no\/core\/apps\/api\/pkg\/runtime\"/g'
find . -type f -name "zz_generated.conversion.go" -print0 | xargs -0 sed -i '' -e 's/\"k8s.io\/apimachinery\/pkg\/conversion\"/\"github.com\/nrc-no\/core\/apps\/api\/pkg\/conversion\"/g'

#array=()
#while IFS= read -r -d ''; do
#  array+=("$REPLY")
#done < <(find . -type f -name "zz_generated.conversion.go" -print0)
#
#for item in "${array[@]}"; do
#
#  package=$(basename $(dirname "${item}"))
#  expr2="s/[[:space:]]+pkg${package} \".*\"//g"
#  expr="s/[[:space:]]+${package} \".*\"//g"
#
#  sed -i -e 's/\"k8s.io\/apimachinery\/pkg\/runtime\"/\"github.com\/nrc-no\/core\/apps\/api\/pkg\/runtime\"/g' "${item}"
#
#  sed -i -E "${expr2}" "${item}"
#  sed -i -E "s/pkg${package}.//g" "${item}"
#  sed -i -E "${expr}" "${item}"
#  sed -i -E "s/${package}.//g" "${item}"
#done
