#!/usr/bin/env bash

echo ">>> Testing package: iam"
go test github.com/nrc-no/core/pkg/iam
if [ ! "$?" ]; then
  exit 1
fi
printf "\n"

echo ">>> Testing package: cms"
go test github.com/nrc-no/core/pkg/cms
if [ ! "$?" ]; then
  exit 1
fi
printf "\n"

echo ">>> Testing package: attachments"
go test github.com/nrc-no/core/pkg/attachments
if [ ! "$?" ]; then
  exit 1
fi
printf "\n"
