#!/bin/bash

"${GOROOT}"/bin/go test -coverprofile=coverage.out --tags=integration
"${GOROOT}"/bin/go tool -func=coverage.out
"${GOROOT}"/bin/go tool cover -html=coverage.out
