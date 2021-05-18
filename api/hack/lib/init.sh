#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

AVAILABLE_GROUP_VERSIONS="${AVAILABLE_GROUP_VERSIONS:-\
apiextensions/v1 \
core/v1 \
}"
