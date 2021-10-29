#!/usr/bin/env bash
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
DOC_DIR="$( cd "${SCRIPT_DIR}"/../docs && pwd )"

echo Rendering dot graphs...

find "${DOC_DIR}" -type f -name "*.dot" -exec sh -c 'echo "Rendering ${0}" && dot -Tpng "${0}" -o "${0%.*}.png"' {} \;

echo Rendering dot graphs...Done!
