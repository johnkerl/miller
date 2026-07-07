#!/bin/bash

# ================================================================
# Builds a single-file epub of the Miller documentation from the generated
# Markdown files in docs/src, for offline reading (issue #1835).
#
# Usage: docs/build-epub.sh [output-path]
# Default output is ./miller-docs.epub.
#
# Requires pandoc (https://pandoc.org). This script is invoked on Read the
# Docs (see .readthedocs.yaml) to publish an epub download alongside the HTML
# docs; Read the Docs' built-in `formats:` support only produces epub/PDF for
# Sphinx projects, not MkDocs ones, hence this script. It can also be run
# locally by anyone with pandoc installed. It is not part of `make dev` or
# any other default build path.
# ================================================================

set -euo pipefail

output="${1:-miller-docs.epub}"
# Make the output path absolute, since we cd around below.
case "$output" in
  /*) ;;
  *) output="$PWD/$output" ;;
esac

docs_dir=$(cd "$(dirname "$0")" && pwd)
src_dir="$docs_dir/src"

if ! command -v pandoc > /dev/null 2>&1; then
  echo "$0: pandoc not found; please install it (https://pandoc.org)." 1>&2
  exit 1
fi

# Chapter ordering is the nav order in mkdocs.yml. Nav entries look like
#   - "Miller in 10 minutes": "10min.md"
# and these are the only lines in mkdocs.yml ending with a quoted .md name.
chapters=$(sed -n 's/^ *-.*: *"\(.*\.md\)" *$/\1/p' "$docs_dir/mkdocs.yml")
if [ -z "$chapters" ]; then
  echo "$0: could not extract any nav entries from $docs_dir/mkdocs.yml." 1>&2
  exit 1
fi

tmp_dir=$(mktemp -d)
trap 'rm -rf "$tmp_dir"' EXIT

# Each generated page starts with a quicklinks navigation block -- raw HTML,
# useful on the website but not in an epub -- which we strip here. It is the
# only <div>...</div> pair at column zero in each page.
inputs=()
for chapter in $chapters; do
  if [ ! -f "$src_dir/$chapter" ]; then
    echo "$0: $src_dir/$chapter is listed in mkdocs.yml nav but does not exist." 1>&2
    echo "$0: perhaps you need to run: make -C $src_dir genmds" 1>&2
    exit 1
  fi
  sed '/^<div>$/,/^<\/div>$/d' "$src_dir/$chapter" > "$tmp_dir/$chapter"
  inputs+=("$tmp_dir/$chapter")
done

# Run from src_dir so relative image paths (pix/*.png) resolve.
cd "$src_dir"
pandoc \
  --toc \
  --toc-depth=2 \
  --split-level=1 \
  --resource-path="$src_dir" \
  --metadata title="Miller Documentation" \
  --metadata author="John Kerl" \
  --metadata lang="en" \
  --output "$output" \
  "${inputs[@]}"

echo "Wrote $output"
