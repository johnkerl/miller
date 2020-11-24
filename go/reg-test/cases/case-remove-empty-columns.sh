run_mlr --csv remove-empty-columns $indir/remove-empty-columns.csv

run_mlr --from $indir/s.csv                    --icsv --opprint remove-empty-columns
run_mlr --from $indir/remove-empty-columns.csv --icsv --opprint cat
run_mlr --from $indir/remove-empty-columns.csv --icsv --opprint remove-empty-columns
run_mlr --icsv --opprint fill-down -f z       $indir/remove-empty-columns.csv
run_mlr --icsv --opprint fill-down -f a,b,c,d $indir/remove-empty-columns.csv
