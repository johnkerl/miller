run_mlr --csvlite rename -r    '^Date_[0-9].*$,Date'  $indir/date1.csv $indir/date2.csv
run_mlr --csvlite rename -r    '(.*)e(.*),\1EEE\2'    $indir/date1.csv $indir/date2.csv
run_mlr --csvlite rename -r    '"(.*)e(.*)"i,\1EEE\2' $indir/date1.csv $indir/date2.csv
run_mlr --csvlite rename -r -g '"(.*)e(.*)"i,\1EEE\2' $indir/date1.csv $indir/date2.csv
run_mlr --csvlite rename -r    '^(.a.e)(_.*)?,AA\1BB\2CC' $indir/date1.csv
run_mlr --csvlite rename -r    '"e",EEE'              $indir/date1.csv $indir/date2.csv
run_mlr --csvlite rename -r -g '"e",EEE'              $indir/date1.csv $indir/date2.csv
run_mlr --csvlite rename -r    '"e"i,EEE'             $indir/date1.csv $indir/date2.csv
run_mlr --csvlite rename -r -g '"e"i,EEE'             $indir/date1.csv $indir/date2.csv
