#!/bin/bash

# Invoke from Miller repo base.

indir=./test/input
outdir=~/data
mkdir -p $outdir

mlr --csv --from $indir/example.csv repeat -n 1000   then shuffle > $outdir/small.csv
mlr --csv --from $indir/example.csv repeat -n 10000  then shuffle > $outdir/medium.csv
mlr --csv --from $indir/example.csv repeat -n 100000 then shuffle > $outdir/big.csv

for kind in small medium big; do
  mlr --c2d cat $outdir/$kind.csv > $outdir/$kind.dkvp
  mlr --c2n cat $outdir/$kind.csv > $outdir/$kind.nidx
  mlr --c2x cat $outdir/$kind.csv > $outdir/$kind.xtab
  mlr --c2j cat $outdir/$kind.csv > $outdir/$kind.json
done
