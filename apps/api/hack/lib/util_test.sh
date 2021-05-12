#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

source "${SCRIPT_DIR}/util.sh"

util::group-version-to-pkg-path meta/v1
util::group-version-to-pkg-path extensions/v1
util::group-version-to-pkg-path extensions/v1
