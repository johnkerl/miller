#!/bin/bash

# Build everything
# go:
#   go build 
#   go test
# man6:
#   creates manpage mlr.1 and manpage.txt using mlr from the $PATH
#   copies the latter to docs6/docs
# docs6: 
#   turn *.md.in into *.md (live code samples), using mlr from the $PATH
#   turn *.md into docs6/site HTML and CSS files

set -euo pipefail

cd go
go fmt ./...
./build

cd ../man6
make

cd ../docs6
./regen.sh
