#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"
ROOT="${DIR}/.."

source "${ROOT}/hack/lib/init.sh"
source "${ROOT}/hack/lib/util.sh"
source "${ROOT}/hack/lib/swagger.sh"

cd "${ROOT}"

IFS=" " read -r -a GROUP_VERSIONS <<<"meta/v1 ${AVAILABLE_GROUP_VERSIONS}"

# To avoid compile errors, remove the currently existing files.
for group_version in "${GROUP_VERSIONS[@]}"; do
  rm -f "$(util::group-version-to-pkg-path "${group_version}")/types_swagger_doc_generated.go"
done

export GO111MODULE=off
go get k8s.io/kubernetes/cmd/genswaggertypedocs
go install k8s.io/kubernetes/cmd/genswaggertypedocs
export GO111MODULE=on

for group_version in "${GROUP_VERSIONS[@]}"; do
  swagger::gen_types_swagger_doc "${group_version}" "$(util::group-version-to-pkg-path "${group_version}")"
done
