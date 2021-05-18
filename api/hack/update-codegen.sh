#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)

CODEGEN_PKG=${CODEGEN_PKG:-$(
  cd "${SCRIPT_ROOT}"
  ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator
)}

bash "${CODEGEN_PKG}/generate-groups.sh" all \
  github.com/nrc-no/coreapi/pkg/generated github.com/nrc-no/coreapi/pkg/apis \
  "core:v1" \
  --go-header-file "${SCRIPT_ROOT}"/hack/boilerplate.go.txt

bash "${CODEGEN_PKG}/generate-internal-groups.sh" "deepcopy,defaulter,conversion,openapi" \
  github.com/nrc-no/coreapi/pkg/generated github.com/nrc-no/coreapi/pkg/apis github.com/nrc-no/coreapi/pkg/apis \
  "core:v1" \
  --go-header-file "${SCRIPT_ROOT}/hack/boilerplate.go.txt"
