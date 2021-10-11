#!/bin/bash

#if [ ! -d "$dir" ]; then
#  echo "Usage: ./integration-coverage-report.sh PATH"
#fi

cd "$(dirname "$0")" || exit
go test -coverprofile=coverage.out github.com/nrc-no/core/pkg/iam github.com/nrc-no/core/pkg/cms github.com/nrc-no/core/pkg/attachments
go tool cover -html=coverage.out
#go tool cover -func=coverage.out
#go tool cover -html=coverage.out
rm coverage.out
