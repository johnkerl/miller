# ----------------------------------------------------------------
announce SPACE-PADDING

run_mlr --idkvp    --odkvp --ifs space --repifs cat $indir/space-pad.dkvp
run_mlr --inidx    --odkvp --ifs space --repifs cat $indir/space-pad.nidx
run_mlr --icsvlite --odkvp --ifs space --repifs cat $indir/space-pad.pprint
