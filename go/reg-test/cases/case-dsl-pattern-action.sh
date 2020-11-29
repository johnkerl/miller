run_mlr --from $indir/s.dkvp put 'NR == 2 { $z = 100 }'
run_mlr --from $indir/s.dkvp put 'NR != 2 { $z = 100 }'
