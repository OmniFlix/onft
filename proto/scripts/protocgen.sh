#!/usr/bin/env bash

set -eo pipefail
echo "Generating gogo proto code"
(cd proto; buf generate --template buf.gen.gogo.yaml)

# move proto files to the right places
cp -r github.com/OmniFlix/onft/types ./
rm -rf github.com