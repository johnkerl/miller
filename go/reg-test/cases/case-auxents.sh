run_mlr lecat --mono < $indir/line-ending-cr.bin
run_mlr lecat --mono < $indir/line-ending-lf.bin
run_mlr lecat --mono < $indir/line-ending-crlf.bin

$path_to_mlr unhex < $indir/256.txt | run_mlr hex
$path_to_mlr unhex < $indir/256.txt | run_mlr hex -r
$path_to_mlr unhex < $indir/256-ragged.txt | run_mlr hex
$path_to_mlr unhex < $indir/256-ragged.txt | run_mlr hex -r
$path_to_mlr unhex $indir/256.txt | run_mlr hex
$path_to_mlr unhex $indir/256.txt | run_mlr hex -r
$path_to_mlr unhex $indir/256-ragged.txt | run_mlr hex
$path_to_mlr unhex $indir/256-ragged.txt | run_mlr hex -r

$path_to_mlr termcvt --cr2lf   < $indir/line-ending-cr.bin   > $outdir/line-ending-temp-1.bin
$path_to_mlr termcvt --cr2crlf < $indir/line-ending-cr.bin   > $outdir/line-ending-temp-2.bin
$path_to_mlr termcvt --lf2cr   < $indir/line-ending-lf.bin   > $outdir/line-ending-temp-3.bin
$path_to_mlr termcvt --lf2crlf < $indir/line-ending-lf.bin   > $outdir/line-ending-temp-4.bin
$path_to_mlr termcvt --crlf2cr < $indir/line-ending-crlf.bin > $outdir/line-ending-temp-5.bin
$path_to_mlr termcvt --crlf2lf < $indir/line-ending-crlf.bin > $outdir/line-ending-temp-6.bin

run_mlr hex < $outdir/line-ending-temp-1.bin
run_mlr hex < $outdir/line-ending-temp-2.bin
run_mlr hex < $outdir/line-ending-temp-3.bin
run_mlr hex < $outdir/line-ending-temp-4.bin
run_mlr hex < $outdir/line-ending-temp-5.bin
run_mlr hex < $outdir/line-ending-temp-6.bin
