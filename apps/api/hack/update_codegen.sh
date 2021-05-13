#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
ROOT="${DIR}/.."

source "${ROOT}/hack/lib/init.sh"
source "${ROOT}/hack/lib/util.sh"

# Install required dependencies
export GO111MODULE=off
go get k8s.io/code-generator/cmd/{defaulter-gen,client-gen,lister-gen,informer-gen,deepcopy-gen,openapi-gen}
go install k8s.io/code-generator/cmd/{defaulter-gen,client-gen,lister-gen,informer-gen,deepcopy-gen,openapi-gen}
export GO111MODULE=on

# Find go binaries
# Go installs the above commands to get installed in $GOBIN if defined, and $GOPATH/bin otherwise:
GOBIN="$(go env GOBIN)"
gobin="${GOBIN:-$(go env GOPATH)/bin}"

# Finding list of group versions to use for code generation
IFS=" " read -r -a GROUP_VERSIONS <<<"meta/v1 ${AVAILABLE_GROUP_VERSIONS}"

# enumerate group versions
PKGS=() # e.g. k8s.io/api/apps/v1
DIRS=() # e.g. k8s.io/api/apps/v1
for GVs in "${GROUP_VERSIONS[@]}"; do
  PKGS+=("github.com/nrc-no/core/apps/api/$(util::group-version-to-pkg-path "${GVs}")")
  DIRS+=("$(util::group-version-to-pkg-path "${GVs}")")
done


echo "Generating openapi files with binary '${gobin}/openapi-gen' for $(codegen::join , "${PKGS[@]}")"
"${gobin}/openapi-gen" \
  --input-dirs "$(codegen::join , "${PKGS[@]}")" \
  --output-file-base zz_generated.openapi \
  --output-base "${ROOT}/../../../../.." \
  --go-header-file "${ROOT}/hack/boilerplate.go.txt" \
  --output-package "github.com/nrc-no/core/apps/api/openapi" \
  -v 10

exit 0

echo "Generating deepcopy funcs with binary '${gobin}/deepcopy-gen' for $(codegen::join , "${PKGS[@]}")"
"${gobin}/deepcopy-gen" \
  --input-dirs "$(codegen::join , "${PKGS[@]}")" \
  --output-file-base zz_generated.deepcopy \
  --output-base "${ROOT}" \
  --go-header-file "${ROOT}/hack/boilerplate.go.txt"

echo "Generating conversion funcs with binary '${gobin}/conversion-gen' for $(codegen::join , "${PKGS[@]}")"
"${gobin}/conversion-gen" \
  --input-dirs "$(codegen::join , "${PKGS[@]}")" \
  --output-file-base zz_generated.conversion \
  --output-base "${ROOT}" \
  --go-header-file "${ROOT}/hack/boilerplate.go.txt" \
  -v 10

echo "Replacing deepcopy file packages"
deepcopy_files=()
while IFS= read -r -d ''; do
  deepcopy_files+=("$REPLY")
done < <(find . -type f -name "zz_generated.deepcopy.go" -print0)

for item in "${deepcopy_files[@]}"; do
  package=$(basename "$(dirname "${item}")")
  expr2="s/[[:space:]]+pkg${package} \".*\"//g"
  expr="s/[[:space:]]+${package} \".*\"//g"
  sed -i -e 's/\"k8s.io\/apimachinery\/pkg\/runtime\"/\"github.com\/nrc-no\/core\/apps\/api\/pkg\/runtime\"/g' "${item}"
  sed -i -E "${expr2}" "${item}"
  sed -i -E "s/pkg${package}.//g" "${item}"
  sed -i -E "${expr}" "${item}"
  sed -i -E "s/${package}.//g" "${item}"
done

echo "Replacing conversion file packages"
conversion_files=()
while IFS= read -r -d ''; do
  conversion_files+=("$REPLY")
done < <(find . -type f -name "zz_generated.conversion.go" -print0)

for item in "${conversion_files[@]}"; do
  sed -i -e 's/k8s.io\/apimachinery\/pkg/github.com\/nrc-no\/core\/apps\/api\/pkg/g' "${item}"

  package=$(basename "$(dirname "${item}")")
  expr2="s/[[:space:]]+pkg${package} \".*\"//g"
  expr="s/[[:space:]]+${package} \".*\"//g"

  sed -i -E "${expr2}" "${item}"
  sed -i -E "s/pkg${package}.//g" "${item}"
  sed -i -E "${expr}" "${item}"
  sed -i -E "s/${package}.//g" "${item}"
done

exit 0
