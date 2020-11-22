run_mlr --icsv --ojson cat                                $indir/clean-whitespace.csv
run_mlr --icsv --ojson put '$a = lstrip($a)'              $indir/clean-whitespace.csv
run_mlr --icsv --ojson put '$a = rstrip($a)'              $indir/clean-whitespace.csv
run_mlr --icsv --ojson put '$a = strip($a)'               $indir/clean-whitespace.csv
run_mlr --icsv --ojson put '$a = collapse_whitespace($a)' $indir/clean-whitespace.csv
run_mlr --icsv --ojson put '$a = clean_whitespace($a)'    $indir/clean-whitespace.csv

run_mlr --icsv --ojson clean-whitespace -k $indir/clean-whitespace.csv
run_mlr --icsv --ojson clean-whitespace -v $indir/clean-whitespace.csv
run_mlr --icsv --ojson clean-whitespace    $indir/clean-whitespace.csv
