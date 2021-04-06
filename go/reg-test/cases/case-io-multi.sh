# TODO: needs --ofmt

ofmt=pprint
for ifmt in csv dkvp nidx json; do
  run_mlr --i$ifmt --o$ofmt put '
    $nf=NF;
    $nr=NR;
    $fnr=FNR;
    $filename=FILENAME;
    $filenum=FILENUM;
  ' $indir/s.$ifmt $indir/t.$ifmt
done

ifmt=dkvp
for ofmt in pprint csv dkvp nidx json; do
  run_mlr --i$ifmt --o$ofmt put '
    $nf=NF;
    $nr=NR;
    $fnr=FNR;
    $filename=FILENAME;
    $filenum=FILENUM;
  ' $indir/s.$ifmt $indir/t.$ifmt
done

run_mlr --ocsv    cat $indir/het.dkvp
run_mlr --opprint cat $indir/het.dkvp

run_mlr --opprint cat <<EOF
EOF

run_mlr --opprint cat <<EOF
a=1,b=2,c=3
EOF

run_mlr --opprint cat <<EOF
a=1,b=2,c=3
a=2,b=2,c=3
EOF

run_mlr --opprint cat <<EOF
a=1,b=2,c=3
a=2,b=2,c=3
d=3,e=5,f=6
EOF

run_mlr --opprint cat <<EOF
a=1,b=2,c=3
d=2,e=5,f=6
d=3,e=5,f=6
EOF

run_mlr --opprint --barred cat $indir/s.dkvp
run_mlr --opprint --barred cat $indir/het.dkvp

# To-do: port format-specific default separators from C.
# E.g. NIDX's IFS should default to space.

run_mlr --inidx --oxtab cat <<EOF
a,b,c,d,e,f
EOF
run_mlr --inidx --oxtab cat <<EOF
a b c d e f
EOF

run_mlr --inidx --ifs , --oxtab cat <<EOF
a,b,c,d,e,f
EOF
run_mlr --inidx --ifs , --oxtab cat <<EOF
a b c d e f
EOF
run_mlr --inidx --ifs ' ' --oxtab cat <<EOF
a,b,c,d,e,f
EOF
run_mlr --inidx --ifs ' ' --oxtab cat <<EOF
a b c d e f
EOF

run_mlr --inidx --ifs comma --oxtab cat <<EOF
a,b,c,d,e,f
EOF
run_mlr --inidx --ifs comma --oxtab cat <<EOF
a b c d e f
EOF
run_mlr --inidx --ifs space --oxtab cat <<EOF
a,b,c,d,e,f
EOF
run_mlr --inidx --ifs space --oxtab cat <<EOF
a b c d e f
EOF

run_mlr --itsv --ocsv cat $indir/s.tsv
run_mlr --icsv --otsv cat $indir/s.tsv
run_mlr --icsv --otsv cat $indir/s.csv
run_mlr --c2j cat $indir/s.csv

run_mlr --ocsv cat $indir/het.dkvp
run_mlr --ocsv --headerless-csv-output cat $indir/het.dkvp
run_mlr --icsv --ojson cat $indir/implicit.csv
run_mlr --implicit-csv-header --icsv --ojson cat $indir/implicit.csv

run_mlr --icsvlite --ojson cat $indir/s.csv
run_mlr --icsvlite --implicit-csv-header --ojson cat $indir/implicit.csv

run_mlr --icsvlite --opprint cat $indir/het-a1.csv $indir/het-a2.csv
run_mlr --icsvlite --opprint cat $indir/het-b1.csv $indir/het-b2.csv
run_mlr --icsvlite --opprint cat $indir/het-c1.csv
run_mlr --icsvlite --opprint cat $indir/het-d1.csv

run_mlr --icsvlite --ojson --allow-ragged-csv-input cat $indir/ragged-short.csv
run_mlr --icsvlite --ojson --allow-ragged-csv-input cat $indir/ragged-long.csv
run_mlr --icsv     --ojson --allow-ragged-csv-input cat $indir/ragged-short.csv
run_mlr --icsv     --ojson --allow-ragged-csv-input cat $indir/ragged-long.csv

run_mlr --ixtab --ojson cat $indir/test-1.xtab
run_mlr --ixtab --ojson cat $indir/test-2.xtab
run_mlr --ixtab --ojson cat $indir/test-3.xtab
run_mlr --ixtab --ojson cat $indir/test-1.xtab $indir/test-2.xtab
run_mlr --ixtab --ojson cat $indir/test-2.xtab $indir/test-1.xtab
run_mlr --ixtab --ojson cat $indir/test-1.xtab $indir/test-2.xtab $indir/test-3.xtab

run_mlr --ojson --from $indir/s.dkvp head -n 0
run_mlr --ojson --from $indir/s.dkvp head -n 1
run_mlr --ojson --from $indir/s.dkvp head -n 2
run_mlr --ojson --from $indir/s.dkvp head -n 3

run_mlr --jlistwrap --ojson --from $indir/s.dkvp head -n 0
run_mlr --jlistwrap --ojson --from $indir/s.dkvp head -n 1
run_mlr --jlistwrap --ojson --from $indir/s.dkvp head -n 2
run_mlr --jlistwrap --ojson --from $indir/s.dkvp head -n 3

