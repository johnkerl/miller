run_mlr count-distinct -f a      $indir/small $indir/abixy
run_mlr count-distinct -f a,b    $indir/small $indir/abixy
run_mlr count-distinct -f a,b -u $indir/small $indir/abixy

run_mlr count-distinct -f a   -n $indir/small $indir/abixy
run_mlr count-distinct -f a,b -n $indir/small $indir/abixy

run_mlr count-distinct -f a   -o foo $indir/small $indir/abixy
run_mlr count-distinct -f a,b -o foo $indir/small $indir/abixy

run_mlr count-distinct -f a   -n -o foo $indir/small $indir/abixy
run_mlr count-distinct -f a,b -n -o foo $indir/small $indir/abixy
