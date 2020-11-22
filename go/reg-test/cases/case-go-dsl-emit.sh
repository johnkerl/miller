# run_mlr --from $indir/s.dkvp --opprint put -q '@sum += $i; emit {"sum": @sum}'
run_mlr --from $indir/s.dkvp --opprint put -q '@sum[$a] += $i; emit {"sum": @sum}'
