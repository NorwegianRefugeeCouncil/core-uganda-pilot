#!/usr/bin/env bash

# Update node packages
cd web/app/client || exit
npm install

cd ../designSystem || exit
npm install

cd ../../../

echo ">>> Updated node dependencies"

echo ">>> Generated typescript types"

tsc --build web/app/client/tsconfig.json

if [ ! "$?" ]; then
  exit 1
fi

tsc --build web/app/designSystem/tsconfig.json

if [ ! "$?" ]; then
  exit 1
fi

echo ">>> Transpiled typescript files"

cd web/app/frontend || exit
rm -rf node_modules/
yarn install

if [ ! "$?" ]; then
  exit 1
fi

echo ">>> Updated frontend"

cd ../../../
go build -o ./tmp/main ./cmd/core

if [ ! "$?" ]; then
  exit 1
fi

echo ">>> Built core"
