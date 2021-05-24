#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)

CODEGEN_PKG=${CODEGEN_PKG:-$(
  cd "${SCRIPT_ROOT}"
  ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator
)}

bash "${CODEGEN_PKG}/generate-groups.sh" deepcopy \
  github.com/nrc-no/core/api/pkg/generated github.com/nrc-no/core/api/pkg/apis \
  "core:v1 discovery:v1 meta:v1" \
  --go-header-file "${SCRIPT_ROOT}"/hack/boilerplate.go.txt

openapi_xtra_pkgs=()
#openapi_xtra_pkgs+=("k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1")
#openapi_xtra_pkgs+=("k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1")

function join_by() {
  local d=${1-} f=${2-}
  if shift 2; then printf %s "$f" "${@/#/$d}"; fi
}

OPENAPI_EXTRA_PACKAGES=$(join_by , "${openapi_xtra_pkgs[@]}") bash "${CODEGEN_PKG}/generate-internal-groups.sh" "deepcopy,defaulter,conversion,openapi" \
  github.com/nrc-no/core/api/pkg/generated github.com/nrc-no/core/api/pkg/apis github.com/nrc-no/core/api/pkg/apis \
  "core:v1 discovery:v1 meta:v1" \
  --go-header-file "${SCRIPT_ROOT}/hack/boilerplate.go.txt"
