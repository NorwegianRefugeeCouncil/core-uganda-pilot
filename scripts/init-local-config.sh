#!/bin/bash

set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"

CUSTOM="${SCRIPT_DIR}/../configs/config.custom.yaml"

if [ ! -f "${CUSTOM}" ]; then
  touch "${CUSTOM}"
  echo "serve:" > "${CUSTOM}"
  echo "  admin:" >> "${CUSTOM}"
  echo "    secrets:" >> "${CUSTOM}"
  echo "      hash:" >> "${CUSTOM}"
  echo "      - $(openssl rand -hex 64)" >> "${CUSTOM}"
  echo "      block:" >> "${CUSTOM}"
  echo "      - $(openssl rand -hex 32)" >> "${CUSTOM}"
  echo "  public:" >> "${CUSTOM}"
  echo "    secrets:" >> "${CUSTOM}"
  echo "      hash:" >> "${CUSTOM}"
  echo "      - $(openssl rand -hex 64)" >> "${CUSTOM}"
  echo "      block:" >> "${CUSTOM}"
  echo "      - $(openssl rand -hex 32)" >> "${CUSTOM}"
  echo "  login:" >> "${CUSTOM}"
  echo "    secrets:" >> "${CUSTOM}"
  echo "      hash:" >> "${CUSTOM}"
  echo "      - $(openssl rand -hex 64)" >> "${CUSTOM}"
  echo "      block:" >> "${CUSTOM}"
  echo "      - $(openssl rand -hex 32)" >> "${CUSTOM}"
fi
