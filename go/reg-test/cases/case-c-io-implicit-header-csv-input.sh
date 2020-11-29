run_mlr --irs crlf                       --icsvlite --ifs ,  --opprint cut -x -f b/ $indir/multi-sep.csv-crlf
run_mlr --irs crlf --implicit-csv-header --icsvlite --ifs ,  --opprint cut -x -f 2  $indir/multi-sep.csv-crlf

run_mlr --irs crlf                       --icsvlite --ifs /, --opprint cut -x -f b  $indir/multi-sep.csv-crlf
run_mlr --irs crlf --implicit-csv-header --icsvlite --ifs /, --opprint cut -x -f 2  $indir/multi-sep.csv-crlf

run_mlr                       --icsv --ifs ,  --opprint cut -x -f b/ $indir/multi-sep.csv-crlf
run_mlr --implicit-csv-header --icsv --ifs ,  --opprint cut -x -f 2  $indir/multi-sep.csv-crlf

run_mlr                       --icsv --ifs /, --opprint cut -x -f b  $indir/multi-sep.csv-crlf
run_mlr --implicit-csv-header --icsv --ifs /, --opprint cut -x -f 2  $indir/multi-sep.csv-crlf

run_mlr --csv -N reorder -f 1,3,2,5,4 $indir/multi-sep.csv-crlf
run_mlr --csv -N reorder -f 5,4,3,2,1 $indir/multi-sep.csv-crlf

run_mlr --icsv --otsv -N cat <<EOF
a,b,c
1,2,3
4,5,6
EOF

run_mlr --ixtab --opprint -N cat <<EOF
a 1
b 2
c 3

a 4
b 5
c 6
EOF

run_mlr --icsv --pprint -N cat <<EOF
1,2,3
4,5,6
EOF
