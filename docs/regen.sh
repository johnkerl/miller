#!/bin/bash

set -euo pipefail

pushd src
./genmds
popd
mkdocs build
