run_mlr         label NEW             $indir/abixy
run_mlr         label a,NEW,c         $indir/abixy
run_mlr         label 1,2,3,4,5,6,7,8 $indir/abixy
run_mlr         label d,x,f           $indir/abixy
mlr_expect_fail label d,x,d           $indir/abixy
