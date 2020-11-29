run_mlr --opprint --from $indir/freq.dkvp most-frequent -f a -n 3
run_mlr --opprint --from $indir/freq.dkvp most-frequent -f a,b -n 3
run_mlr --opprint --from $indir/freq.dkvp most-frequent -f a,b -n 3 -b
run_mlr --opprint --from $indir/freq.dkvp most-frequent -f nonesuch -n 3

run_mlr --opprint --from $indir/freq.dkvp least-frequent -f a -n 3
run_mlr --opprint --from $indir/freq.dkvp least-frequent -f a,b -n 3
run_mlr --opprint --from $indir/freq.dkvp least-frequent -f a,b -n 3 -b
run_mlr --opprint --from $indir/freq.dkvp least-frequent -f nonesuch -n 3

run_mlr --opprint --from $indir/freq.dkvp most-frequent -f a -n 3 -o foo
run_mlr --opprint --from $indir/freq.dkvp most-frequent -f a,b -n 3 -o foo
run_mlr --opprint --from $indir/freq.dkvp most-frequent -f a,b -n 3 -b -o foo
run_mlr --opprint --from $indir/freq.dkvp most-frequent -f nonesuch -n 3 -o foo

run_mlr --opprint --from $indir/freq.dkvp least-frequent -f a -n 3 -o foo
run_mlr --opprint --from $indir/freq.dkvp least-frequent -f a,b -n 3 -o foo
run_mlr --opprint --from $indir/freq.dkvp least-frequent -f a,b -n 3 -b -o foo
run_mlr --opprint --from $indir/freq.dkvp least-frequent -f nonesuch -n 3 -o foo
