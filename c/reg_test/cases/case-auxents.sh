run_mlr lecat --mono < $indir/line-ending-cr.bin
run_mlr lecat --mono < $indir/line-ending-lf.bin
run_mlr lecat --mono < $indir/line-ending-crlf.bin

run_mlr_no_output termcvt --cr2lf   < $indir/line-ending-cr.bin   > $outdir/line-ending-temp-1.bin
run_mlr_no_output termcvt --cr2crlf < $indir/line-ending-cr.bin   > $outdir/line-ending-temp-2.bin
run_mlr_no_output termcvt --lf2cr   < $indir/line-ending-lf.bin   > $outdir/line-ending-temp-3.bin
run_mlr_no_output termcvt --lf2crlf < $indir/line-ending-lf.bin   > $outdir/line-ending-temp-4.bin
run_mlr_no_output termcvt --crlf2cr < $indir/line-ending-crlf.bin > $outdir/line-ending-temp-5.bin
run_mlr_no_output termcvt --crlf2lf < $indir/line-ending-crlf.bin > $outdir/line-ending-temp-6.bin

run_mlr hex < $outdir/line-ending-temp-1.bin
run_mlr hex < $outdir/line-ending-temp-2.bin
run_mlr hex < $outdir/line-ending-temp-3.bin
run_mlr hex < $outdir/line-ending-temp-4.bin
run_mlr hex < $outdir/line-ending-temp-5.bin
run_mlr hex < $outdir/line-ending-temp-6.bin
