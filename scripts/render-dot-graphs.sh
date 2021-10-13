#!/usr/bin/env bash
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
dot -Tpng "${SCRIPT_DIR}/../docs/overview.dot" -o "${SCRIPT_DIR}/../docs/overview.png"


