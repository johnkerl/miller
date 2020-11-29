run_mlr --csv cat $indir/rfc-csv/simple.csv-crlf
run_mlr --csv cat $indir/rfc-csv/simple-truncated.csv
run_mlr --csv cat $indir/rfc-csv/narrow.csv
run_mlr --csv cat $indir/rfc-csv/narrow-truncated.csv
run_mlr --csv cat $indir/rfc-csv/quoted-comma.csv
run_mlr --csv cat $indir/rfc-csv/quoted-comma-truncated.csv
run_mlr --csv cat $indir/rfc-csv/quoted-crlf.csv
run_mlr --csv cat $indir/rfc-csv/quoted-crlf-truncated.csv
run_mlr --csv cat $indir/rfc-csv/simple-truncated.csv $indir/rfc-csv/simple.csv-crlf
run_mlr --csv --ifs semicolon --ofs pipe --irs lf --ors lflf cut -x -f b $indir/rfc-csv/modify-defaults.csv
run_mlr --csv --rs lf --quote-original cut -o -f c,b,a $indir/quote-original.csv

run_mlr --icsv --oxtab cat $indir/comma-at-eof.csv

run_mlr --csv --quote-all      cat $indir/rfc-csv/simple.csv-crlf
run_mlr --csv --quote-original cat $indir/rfc-csv/simple.csv-crlf

run_mlr --itsv --rs lf --oxtab cat $indir/simple.tsv

run_mlr --iusv --oxtab cat $indir/example.usv
