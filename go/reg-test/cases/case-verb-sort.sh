run_mlr --opprint --from $indir/s.dkvp sort -f nonesuch

run_mlr --opprint --from $indir/s.dkvp sort -f a
run_mlr --opprint --from $indir/s.dkvp sort -r a

run_mlr --opprint --from $indir/s.dkvp sort -f i
run_mlr --opprint --from $indir/s.dkvp sort -r i

run_mlr --opprint --from $indir/s.dkvp sort -nf i
run_mlr --opprint --from $indir/s.dkvp sort -nr i

run_mlr --opprint --from $indir/s.dkvp sort -f x
run_mlr --opprint --from $indir/s.dkvp sort -r x

run_mlr --opprint --from $indir/s.dkvp sort -nf x
run_mlr --opprint --from $indir/s.dkvp sort -nr x

run_mlr --opprint --from $indir/s.dkvp sort -f a,b
run_mlr --opprint --from $indir/s.dkvp sort -r a,b

run_mlr --opprint --from $indir/s.dkvp sort -f x,y
run_mlr --opprint --from $indir/s.dkvp sort -r x,y

run_mlr --opprint --from $indir/s.dkvp sort -nf x,y
run_mlr --opprint --from $indir/s.dkvp sort -nr x,y

run_mlr --opprint --from $indir/s.dkvp sort -f a -f b
run_mlr --opprint --from $indir/s.dkvp sort -f a -nr x
run_mlr --opprint --from $indir/s.dkvp sort -nr y -f a
run_mlr --opprint --from $indir/s.dkvp sort -f a -r b -nf x -nr y

run_mlr --from $indir/sort-het.dkvp sort -f x
run_mlr --from $indir/sort-het.dkvp sort -r x
