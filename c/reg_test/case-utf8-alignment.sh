# ----------------------------------------------------------------
announce UTF-8 alignment

run_mlr --icsvlite --opprint cat $indir/utf8-1.csv
run_mlr --icsvlite --opprint cat $indir/utf8-2.csv
run_mlr --icsvlite --oxtab   cat $indir/utf8-1.csv
run_mlr --icsvlite --oxtab   cat $indir/utf8-2.csv

run_mlr --inidx --ifs space --opprint         cat $indir/utf8-align.nidx
run_mlr --inidx --ifs space --opprint --right cat $indir/utf8-align.nidx
run_mlr --oxtab cat $indir/utf8-align.dkvp

run_mlr --inidx --ifs space --oxtab --xvright cat $indir/utf8-align.nidx
