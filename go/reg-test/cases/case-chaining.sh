run_mlr cat then cat $indir/short
run_mlr cat then tac $indir/short
run_mlr tac then cat $indir/short
run_mlr tac then tac $indir/short

run_mlr cat then cat then cat $indir/short
run_mlr cat then cat then tac $indir/short
run_mlr cat then tac then cat $indir/short
run_mlr cat then tac then tac $indir/short
run_mlr tac then cat then cat $indir/short
run_mlr tac then cat then tac $indir/short
run_mlr tac then tac then cat $indir/short
run_mlr tac then tac then tac $indir/short

# Test allowing then-chains to start with an initial 'then'
run_mlr \
    then cat \
    then head -n 2 -g a,b \
    then tac \
    $indir/abixy-het
