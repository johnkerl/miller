
run_mlr --opprint --from $indir/abixy put -f $indir/put-script-piece-1
run_mlr --opprint --from $indir/abixy put -f $indir/put-script-piece-1 -f $indir/put-script-piece-2
run_mlr --opprint --from $indir/abixy put -f $indir/put-script-piece-1 -f $indir/put-script-piece-2 -f $indir/put-script-piece-3

run_mlr --opprint --from $indir/abixy put -e '$xy = $x**2 + $y**2'
run_mlr --opprint --from $indir/abixy filter -e 'NR == 7'

run_mlr --opprint --from $indir/abixy put -e 'print "PRE";' -f $indir/put-script-piece-1 -f $indir/put-script-piece-2 -f $indir/put-script-piece-3 -e 'print "POST"'

run_mlr --opprint --from $indir/abixy filter -f $indir/filter-script-piece-1 -f $indir/filter-script-piece-2
