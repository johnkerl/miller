
run_mlr put    -f $indir/put-example.dsl $indir/abixy
run_mlr filter -f $indir/filter-example.dsl $indir/abixy

run_mlr --from $indir/abixy put    -f $indir/put-example.dsl
run_mlr --from $indir/abixy filter -f $indir/filter-example.dsl

run_mlr --from $indir/abixy --from $indir/abixy-het put    -f $indir/put-example.dsl
run_mlr --from $indir/abixy --from $indir/abixy-het filter -f $indir/filter-example.dsl
