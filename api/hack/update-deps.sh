#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
rm -rf "${ROOT}/third-party"
mkdir -p "${ROOT}/third-party"
wget https://github.com/kubernetes/kubernetes/archive/refs/heads/master.zip -O "${ROOT}/third-party/k8s.zip"
unzip "${ROOT}/third-party/k8s.zip" -d "${ROOT}/third-party"
mv "${ROOT}/third-party/kubernetes-master" "${ROOT}/third-party/kubernetes"
rm "${ROOT}/third-party/k8s.zip"
