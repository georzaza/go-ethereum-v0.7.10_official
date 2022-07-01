#!/bin/bash

set -e

TEST_DEPS=$(go list -f '{{.Imports}} {{.TestImports}} {{.XTestImports}}' github.com/georzaza/go-ethereum-v0.7.10_official/... | sed -e 's/\[//g' | sed -e 's/\]//g' | sed -e 's/C //g')
if [ "$TEST_DEPS" ]; then
  go get -race $TEST_DEPS
fi
