#!/usr/bin/env bash

# Update node packages
cd web/app/client || exit
npm install

cd ../designSystem || exit
npm install

cd ../../../test/e2e || exit
npm install

cd ../../

echo ">>> Updated node dependencies"

# Generate typescript types for go types
go build -o ./tmp/gotypes2ts ./tools/gotypes2ts
./tmp/gotypes2ts web/app/client/src/types/models.ts >/dev/null

# remove problematic line
# FIXME this shouldn't be necessary
grep -v '?: PartyTypeRule;' web/app/client/src/types/models.ts >tmp/scratch
cat tmp/scratch >web/app/client/src/types/models.ts

if [ ! "$?" ]; then
  exit 1
fi

echo ">>> Generated typescript types"

# Transpile typescript
tsc --build tsconfig.json

if [ ! "$?" ]; then
  exit 1
fi

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
