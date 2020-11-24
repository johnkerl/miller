run_mlr fraction -f x,y              $indir/abixy-het
run_mlr fraction -f x,y -g a         $indir/abixy-het
run_mlr fraction -f x,y -g a,b       $indir/abixy-het

run_mlr fraction -f x,y        -p    $indir/abixy-het
run_mlr fraction -f x,y -g a   -p    $indir/abixy-het
run_mlr fraction -f x,y -g a,b -p    $indir/abixy-het

run_mlr fraction -f x,y           -c $indir/abixy-het
run_mlr fraction -f x,y -g a      -c $indir/abixy-het
run_mlr fraction -f x,y -g a,b    -c $indir/abixy-het

run_mlr fraction -f x,y        -p -c $indir/abixy-het
run_mlr fraction -f x,y -g a   -p -c $indir/abixy-het
run_mlr fraction -f x,y -g a,b -p -c $indir/abixy-het
