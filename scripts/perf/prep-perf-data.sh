#!/bin/bash

# Invoke from Miller repo base.

set -x

indir=./test/input
outdir=~/data
mkdir -p $outdir

mlr --csv \
  --from docs/src/example.csv \
  repeat -n 100000 \
  then shuffle \
  then put '
    begin{@index=1}
    $k = NR;
    @index += urandint(2,10);
    $index=@index;
    $quantity=fmtnum(urandrange(50,100),"%.4f");
    $rate=fmtnum(urandrange(1,10),"%.4f");
  ' \
> $outdir/big.csv

mlr --csv head -n  1000 $outdir/big.csv > $outdir/small.csv
mlr --csv head -n 10000 $outdir/big.csv > $outdir/medium.csv

for kind in small medium big; do
  mlr --c2d cat $outdir/$kind.csv > $outdir/$kind.dkvp
  mlr --c2n cat $outdir/$kind.csv > $outdir/$kind.nidx
  mlr --c2x cat $outdir/$kind.csv > $outdir/$kind.xtab
  mlr --c2j cat $outdir/$kind.csv > $outdir/$kind.json
done
