#!/bin/bash

# Build everything
# go:
#   go build 
#   go test
# man6:
#   creates manpage mlr.1 and manpage.txt using mlr from the $PATH
#   copies the latter to docs6/src
# docs6: 
#   turn *.md.in into *.md (live code samples), using mlr from the $PATH
#   turn *.md into docs6/site HTML and CSS files

set -euo pipefail

cd go
go fmt ./...
gofmt -s -w $(find . -name \*.go | grep -v src/parsing)
./build

cd ../man6
make maybeinstallhome

cd ../docs6
./regen.sh

echo
echo DONE
