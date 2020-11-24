cp $indir/abixy $outdir/abixy.temp1
cp $indir/abixy $outdir/abixy.temp2
run_cat $outdir/abixy.temp1
run_cat $outdir/abixy.temp2
run_mlr -I --opprint head -n 2 $outdir/abixy.temp1 $outdir/abixy.temp2
run_cat $outdir/abixy.temp1
run_cat $outdir/abixy.temp2

mlr_expect_fail -I --opprint head -n 2 < $outdir/abixy.temp1
mlr_expect_fail -I --opprint -n head -n 2 $outdir/abixy.temp1

cp $indir/abixy $outdir/abixy.temp1
cp $indir/abixy $outdir/abixy.temp2
run_cat $outdir/abixy.temp1
run_cat $outdir/abixy.temp2
run_mlr -I --opprint rename a,AYE,b,BEE $outdir/abixy.temp1 $outdir/abixy.temp2
run_cat $outdir/abixy.temp1
run_cat $outdir/abixy.temp2
