run_mlr seqgen --start 1 --stop 5 --step  1
run_mlr seqgen --start 1 --stop 5 --step  2
run_mlr seqgen --start 1 --stop 1 --step  1 -f a
run_mlr seqgen --start 5 --stop 1 --step  1 -f b
run_mlr seqgen --start 5 --stop 1 --step -1 -f c
run_mlr seqgen --start 5 --stop 5 --step -1 -f d
run_mlr --from $indir/abixy cat then seqgen --start 1 --stop 5
