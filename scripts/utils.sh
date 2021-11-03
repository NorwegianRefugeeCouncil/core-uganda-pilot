#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

# Usage
# createFileIfNotExists some/path somevalue
function createFileIfNotExists() {
  local FILEPATH=$1
  local VALUE=$2
  local DIR
  DIR=$(dirname "${FILEPATH}")
  mkdir -p "${DIR}"
  if [ ! -f "${FILEPATH}" ]; then
    echo "${VALUE}" >"${FILEPATH}"
    echo -n "${VALUE}"
  else
    cat "${FILEPATH}"
  fi

}
