# Intended to be invoked by "." from reg-test/run

run_mlr cut -f a,x $indir/abixy
run_mlr cut --complement -f a,x $indir/abixy

run_mlr cut -r    -f c,e         $indir/having-fields-regex.dkvp
run_mlr cut -r    -f '"C","E"'   $indir/having-fields-regex.dkvp
run_mlr cut -r    -f '"c"i,"e"'  $indir/having-fields-regex.dkvp
run_mlr cut -r    -f '"C"i,"E"'  $indir/having-fields-regex.dkvp
run_mlr cut -r -x -f c,e         $indir/having-fields-regex.dkvp
run_mlr cut -r -x -f '"C","E"'   $indir/having-fields-regex.dkvp
run_mlr cut -r -x -f '"C","E"i'  $indir/having-fields-regex.dkvp
run_mlr cut -r -x -f '"c","e"i'  $indir/having-fields-regex.dkvp

run_mlr --csvlite cut -r -f '^Name$,^Date_[0-9].*$' $indir/date1.csv $indir/date2.csv
