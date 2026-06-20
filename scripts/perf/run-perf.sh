#!/bin/bash
# Usage: run-perf.sh mlr-executable [mlr-executable ...]
# Run from scripts/perf/.

set -e

dir="$(cd "$(dirname "$0")" && pwd)"

if [ $# -eq 0 ]; then
  echo "Usage: $0 mlr-executable [mlr-executable ...]" >&2
  echo "Example: $0 ~/bin/mlr-6.18.1 ~/bin/mlr-6.19.0" >&2
  exit 1
fi

datfile="timings-$(date +%Y-%m-%d).dat"

echo "Collecting timings into $datfile ..."
python "$dir/time-verbs.py" "$@" | tee "$datfile"

echo ""
echo "Generating plots ..."
bash "$dir/plot-timings.sh" "$datfile"
