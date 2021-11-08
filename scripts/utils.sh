#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

# Usage
# createFileIfNotExists some/path somevalue
function createFileIfNotExists() {
  local FILEPATH=$1
  local DIR
  DIR=$(dirname "${FILEPATH}")
  mkdir -p "${DIR}"
  if [ ! -f "${FILEPATH}" ]; then
    eval "$2" >"${FILEPATH}"
    cat "${FILEPATH}"
  else
    cat "${FILEPATH}"
  fi

}
