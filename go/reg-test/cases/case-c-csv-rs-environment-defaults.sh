# ----------------------------------------------------------------
announce CSV/RS ENVIRONMENT DEFAULTS

run_mlr --csv cut -f a $indir/rfc-csv/simple.csv-crlf
run_mlr --csv --rs crlf cut -f a $indir/rfc-csv/simple.csv-crlf
mlr_expect_fail --csv --rs lf cut -f a $indir/rfc-csv/simple.csv-crlf
