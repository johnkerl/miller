mlr_expect_fail --from $indir/ten.dkvp gap
run_mlr --from $indir/ten.dkvp gap -n 4
run_mlr --from $indir/ten.dkvp gap -g a
run_mlr --from $indir/ten.dkvp sort -f a then gap -g a
run_mlr --from $indir/ten.dkvp sort -f a,b then gap -g a,b
