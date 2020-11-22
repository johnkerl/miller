run_mlr top -f x,y -n 2        $indir/abixy-het
run_mlr top -f x,y -n 2 -g a   $indir/abixy-het
run_mlr top -f x,y -n 2 -g a,b $indir/abixy-het

run_mlr top -f x,y -n 2        $indir/ints.dkvp
run_mlr top -f x,y -n 2 -F     $indir/ints.dkvp

run_mlr top    -n 4 -f x        $indir/abixy-wide
run_mlr top    -n 1 -f x,y      $indir/abixy-wide
run_mlr top    -n 4 -f x   -g a $indir/abixy-wide
run_mlr top    -n 1 -f x,y -g a $indir/abixy-wide
run_mlr top -a -n 4 -f x        $indir/abixy-wide
run_mlr top -a -n 4 -f x   -g a $indir/abixy-wide

run_mlr top    -n 3 -f x,y       $indir/near-ovf.dkvp
run_mlr top    -n 3 -f x,y --min $indir/near-ovf.dkvp
run_mlr top -F -n 3 -f x,y       $indir/near-ovf.dkvp
run_mlr top -F -n 3 -f x,y --min $indir/near-ovf.dkvp
