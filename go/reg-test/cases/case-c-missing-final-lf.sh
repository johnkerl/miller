# ----------------------------------------------------------------
announce MISSING FINAL LF

run_mlr --csvlite cat $indir/truncated.csv
run_mlr --dkvp    cat $indir/truncated.dkvp
run_mlr --nidx    cat $indir/truncated.nidx
run_mlr --pprint  cat $indir/truncated.pprint
run_mlr --xtab    cat $indir/truncated.xtab-crlf
