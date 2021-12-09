#!/usr/bin/env bash

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

# Short git commit sha
GIT_COMMIT="$(git rev-parse --short HEAD)"

# git commit timestamp
GIT_TS="$(git show -s --format=%ct)"

# Tag for the core server. Core server contains code for and can run
# - login-frontend
# - forms-api
# - authnz-api
# - authnz-bouncer
CORE_TAG="sandies.azurecr.io/core:${GIT_TS}-${GIT_COMMIT}"

# Tag for the authnz-frontend app
AUTHNZ_FRONTEND_TAG="sandies.azurecr.io/authnz-frontend:${GIT_TS}-${GIT_COMMIT}"

# Tag for the core-frontend app
CORE_FRONTEND_TAG="sandies.azurecr.io/core-frontend:${GIT_TS}-${GIT_COMMIT}"

# Moving to the root source dir
cd "${ROOT_DIR}" || exit 1

# Building core
DOCKER_BUILDKIT=1 docker build . --build-arg git_commit='${GIT_COMMIT}' -f build/package/core.Dockerfile --progress=plain --tag "${CORE_TAG}"

# Building authnz-frontend app
DOCKER_BUILDKIT=1 docker build . --build-arg git_commit='${GIT_COMMIT}' -f build/package/authnz-frontend.Dockerfile --progress=plain --tag "${AUTHNZ_FRONTEND_TAG}"

# Building core-frontend app
DOCKER_BUILDKIT=1 docker build . --build-arg git_commit='${GIT_COMMIT}' -f build/package/core-frontend.Dockerfile --progress=plain --tag "${CORE_FRONTEND_TAG}"

# Pushing images
docker push "${CORE_TAG}"
docker push "${AUTHNZ_FRONTEND_TAG}"
docker push "${CORE_FRONTEND_TAG}"
