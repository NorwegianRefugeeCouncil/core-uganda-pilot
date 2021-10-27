SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)"

set -e

# cleanup on exit
trap 'jobs -p | xargs -r kill' EXIT

(cd "${SCRIPT_DIR}/.." && docker build -t digitalbackbonecr.azurecr.io/nrc-no/core:latest -f build/package/core.Dockerfile .)
