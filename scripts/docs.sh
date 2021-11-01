#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

echo Building documentation...

"${SCRIPT_DIR}/render-dot-graphs.sh"


echo Building documentation...Done!
