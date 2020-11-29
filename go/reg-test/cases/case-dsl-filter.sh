run_mlr --from $indir/s.dkvp --opprint put       'filter NR > 2'
run_mlr --from $indir/s.dkvp --opprint put -x    'filter NR > 2'
run_mlr --from $indir/s.dkvp --opprint put       'NR > 2'
run_mlr --from $indir/s.dkvp --opprint put -x    'NR > 2'
run_mlr --from $indir/s.dkvp --opprint filter    'NR > 2'
run_mlr --from $indir/s.dkvp --opprint filter -x 'NR > 2'

# The bare-boolean filter condition needn't be the last statement.
run_mlr --from $indir/abixy --opprint filter '$u=1; NR > 3; $v=2'
run_mlr --from $indir/abixy --opprint put    '$u=1; NR > 3; $v=2'
