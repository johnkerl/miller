#!/usr/bin/env bash
# Run each cut implementation N times and print mean execution time.
# Usage: ./run-cut-benchmark.sh [step] [fields] [file ...]
#   step: 1-6 (default 6 = full pipeline)
#   fields: comma-separated (default "a,x")
#   file: input file(s); default stdin (you can pipe or pass a file)
# Example: ./run-cut-benchmark.sh 6 a,x data.dkvp
# Example: ./run-cut-benchmark.sh 6 ccode,milex,year,cinc ../c/nmc1.dkvp

set -e
N=5
STEP="${1:-6}"
FIELDS="${2:-a,x}"
shift 2 2>/dev/null || true
FILES=("$@")
if [ ${#FILES[@]} -eq 0 ]; then
  FILES=("-")
fi

run_one() {
  local bin="$1"
  if [ ! -x "$bin" ]; then
    echo "$bin: binary not found (build with make)"
    return
  fi
  local sum=0
  local i
  for i in $(seq 1 "$N"); do
    local start end t
    start=$(python3 -c 'import time; print(time.time())')
    "$bin" "$STEP" "$FIELDS" "${FILES[@]}" >/dev/null
    end=$(python3 -c 'import time; print(time.time())')
    t=$(echo "$end - $start" | bc -l)
    sum=$(echo "$sum + $t" | bc -l)
  done
  local mean
  mean=$(echo "scale=4; $sum / $N" | bc -l)
  printf "%-8s  %s s (mean of %d runs)\n" "$bin" "$mean" "$N"
}

echo "Cut benchmark: step=$STEP fields=$FIELDS files=${FILES[*]} (N=$N)"
echo "---"
run_one "./cutgo"
run_one "./cutcpp"
run_one "./cutnim"
run_one "./cutrs"
run_one "./cutzig"
