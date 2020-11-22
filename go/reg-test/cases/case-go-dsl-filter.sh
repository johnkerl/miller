run_mlr --from $indir/s.dkvp --opprint put 'filter NR > 2'
run_mlr --from $indir/s.dkvp --opprint filter 'NR > 2'
run_mlr --from $indir/s.dkvp --opprint filter -x 'NR > 2'
