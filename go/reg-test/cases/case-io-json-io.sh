run_mlr --ijson --opprint cat $indir/small-non-nested.json
run_mlr --ijson --opprint cat $indir/small-non-nested-wrapped.json
run_mlr --ijson --oxtab   cat $indir/small-nested.json

run_mlr --ojson                                       cat $indir/json-output-options.dkvp
run_mlr --ojson --jvstack                             cat $indir/json-output-options.dkvp
run_mlr --ojson             --jlistwrap               cat $indir/json-output-options.dkvp
run_mlr --ojson --jvstack   --jlistwrap               cat $indir/json-output-options.dkvp
run_mlr --ojson                         --jquoteall   cat $indir/json-output-options.dkvp
run_mlr --ojson --jvstack               --jquoteall   cat $indir/json-output-options.dkvp
run_mlr --ojson             --jlistwrap --jquoteall   cat $indir/json-output-options.dkvp
run_mlr --ojson --jvstack   --jlistwrap --jquoteall   cat $indir/json-output-options.dkvp
run_mlr --ojson                         --jvquoteall  cat $indir/json-output-options.dkvp
run_mlr --ojson --jvstack               --jvquoteall  cat $indir/json-output-options.dkvp
run_mlr --ojson             --jlistwrap --jvquoteall  cat $indir/json-output-options.dkvp
run_mlr --ojson --jvstack   --jlistwrap --jvquoteall  cat $indir/json-output-options.dkvp
run_mlr --ojson                         --jknquoteint cat $indir/json-output-options.dkvp
run_mlr --ojson --jvstack               --jknquoteint cat $indir/json-output-options.dkvp
run_mlr --ojson             --jlistwrap --jknquoteint cat $indir/json-output-options.dkvp
run_mlr --ojson --jvstack   --jlistwrap --jknquoteint cat $indir/json-output-options.dkvp

run_mlr put -q --jvquoteall 'dump $*'                   $indir/json-output-options.dkvp
run_mlr put -q --jvquoteall 'o = $*; o[7] = 8; dump o'  $indir/json-output-options.dkvp
run_mlr put -q --jknquoteint 'dump $*'                  $indir/json-output-options.dkvp
run_mlr put -q --jknquoteint 'o = $*; o[7] = 8; dump o' $indir/json-output-options.dkvp

run_mlr  --ijson --opprint cat $indir/small-non-nested-wrapped.json $indir/small-non-nested-wrapped.json

run_mlr --icsv --ojson --rs lf cat <<EOF
col
"abc ""def"" \ghi"
EOF

run_mlr --icsv --ojson --jvquoteall --rs lf cat <<EOF
col
"abc ""def"" \ghi"
EOF

run_mlr --ijson --oxtab                              cat $indir/arrays.json
run_mlr --ijson --oxtab --json-map-arrays-on-input   cat $indir/arrays.json
run_mlr --ijson --oxtab --json-skip-arrays-on-input  cat $indir/arrays.json
run_mlr --ijson --oxtab --json-fatal-arrays-on-input cat $indir/arrays.json

run_mlr --json cat $indir/escapes.json
