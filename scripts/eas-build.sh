#!/usr/bin/env bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

set -e

PLATFORM=android
PROFILE=$1

if [ "$PROFILE" != "development" ] && [ "$PROFILE" != "staging" ] && [ "$PROFILE" != "production" ]
then
  echo "Invalid profile"
  exit 1
fi

if ! command -v eas &> /dev/null
then
    echo "eas could not be found"
    echo "Run: yarn global add eas-cli"
    exit 1
fi


WHOAMI=$(eas whoami || true)
if [ "$WHOAMI" == "Not logged in" ]
then
    echo "You are not logged in"
    echo "Run: eas login"
    exit 1
fi 

cd "${ROOT_DIR}/frontend/apps/core-app"

eas build --profile $PROFILE --platform $PLATFORM 