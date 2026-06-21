#!/bin/bash
# Usage: plot-timings.sh timings.dat
# Requires: mlr, pgr
# Run from the Miller repo root.

set -e

datafile="${1:?Usage: $0 timings.dat}"
stem="${datafile%.dat}"

mlr --d2p --from "$datafile" \
  grep cat then reshape -s desc,seconds \
  | sed '1s/^/#/' \
  | pgr -cat -ymin 0 -flabels -lur -lp -ms 5 -o "${stem}-cats.png" &

mlr --d2p --from "$datafile" \
  grep chain then reshape -s desc,seconds \
  | sed '1s/^/#/' \
  | pgr -cat -ymin 0 -flabels -lur -lp -ms 5 -o "${stem}-chains.png" &

mlr --d2p --from "$datafile" \
  grep -v cat then grep -v chain then reshape -s desc,seconds \
  | sed '1s/^/#/' \
  | pgr -cat -ymin 0 -flabels -lur -lp -ms 5 -o "${stem}-verbs.png" &

wait
echo "Plots written:"
echo "  ${stem}-cats.png"
echo "  ${stem}-chains.png"
echo "  ${stem}-verbs.png"
