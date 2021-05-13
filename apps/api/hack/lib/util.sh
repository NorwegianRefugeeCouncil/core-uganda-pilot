#!/usr/bin/env bash

util::group-version-to-pkg-path() {
  local group_version="$1"

  # "v1" is the API GroupVersion
  if [[ "${group_version}" == "v1" ]]; then
    echo "vendor/k8s.io/api/core/v1"
    return
  fi

  case "${group_version}" in
  # both group and version are "", this occurs when we generate deep copies for internal objects of the legacy v1 API.
  __internal)
    echo "pkg/apis/core"
    ;;
  meta/v1)
    echo "pkg/apis/meta/v1"
    ;;
  meta/v1beta1)
    echo "vendor/k8s.io/apimachinery/pkg/apis/meta/v1beta1"
    ;;
  internal.apiserver.k8s.io/v1alpha1)
    echo "vendor/k8s.io/api/apiserverinternal/v1alpha1"
    ;;
  *.k8s.io)
    echo "pkg/apis/${group_version%.*k8s.io}"
    ;;
  *.k8s.io/*)
    echo "pkg/apis/${group_version/.*k8s.io/}"
    ;;
  *)
    echo "pkg/apis/${group_version%__internal}"
    ;;
  esac

}

function codegen::join() {
  local IFS="$1"
  shift
  echo "$*"
}
