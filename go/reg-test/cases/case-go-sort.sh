run_mlr --opprint --from $indir/s.dkvp sort -f nonesuch
run_mlr --opprint --from $indir/s.dkvp sort -f a
run_mlr --opprint --from $indir/s.dkvp sort -f a,b
run_mlr --opprint --from $indir/s.dkvp sort -r a
run_mlr --opprint --from $indir/s.dkvp sort -r a,b
run_mlr --opprint --from $indir/s.dkvp sort -f a -r b
run_mlr --opprint --from $indir/s.dkvp sort -f b -n i
run_mlr --opprint --from $indir/s.dkvp sort -f b -nr i
