run_mlr --csv  --prepipe 'cat'   cat   $indir/rfc-csv/simple.csv-crlf
run_mlr --dkvp --prepipe 'cat'   cat   $indir/abixy
run_mlr --csv  --prepipe 'cat'   cat < $indir/rfc-csv/simple.csv-crlf
run_mlr --dkvp --prepipe 'cat'   cat < $indir/abixy
