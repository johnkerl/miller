#!/bin/bash

set -euo pipefail

pushd docs
./genmds
popd
mkdocs build
