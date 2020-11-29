run_mlr --icsvlite --opprint put '$langue = toupper($langue)' $indir/utf8-1.csv
run_mlr --icsvlite --opprint put '$nom    = toupper($nom)'    $indir/utf8-1.csv
run_mlr --icsvlite --opprint put '$jour   = toupper($jour)'   $indir/utf8-1.csv

run_mlr --icsvlite --opprint put '$langue = capitalize($langue)' $indir/utf8-1.csv
run_mlr --icsvlite --opprint put '$nom    = capitalize($nom)'    $indir/utf8-1.csv
run_mlr --icsvlite --opprint put '$jour   = capitalize($jour)'   $indir/utf8-1.csv
