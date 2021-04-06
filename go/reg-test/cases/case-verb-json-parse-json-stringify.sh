# TODO: git-checkin the intermediates

run_mlr --ijson --ojson --from $indir/flatten-input-2.json json-stringify
run_mlr --ijson --oxtab --from $indir/flatten-input-2.json json-stringify then flatten

run_mlr --ijson --ojson --from $indir/flatten-input-2.json json-stringify -f req
run_mlr --ijson --oxtab --from $indir/flatten-input-2.json json-stringify -f req then flatten

run_mlr --ijson --ojson --from $indir/flatten-input-2.json json-stringify -f req --jvstack
run_mlr --ijson --ojson --from $indir/flatten-input-2.json json-stringify -f req --no-jvstack

$path_to_mlr --ijson --ojson --from $indir/flatten-input-2.json json-stringify \
  | run_mlr --j2x flatten
$path_to_mlr --ijson --ojson --from $indir/flatten-input-2.json json-stringify \
  | run_mlr --j2x json-parse then flatten

$path_to_mlr --ijson --ojson --from $indir/flatten-input-2.json json-stringify --jvstack \
  | run_mlr --j2x flatten
$path_to_mlr --ijson --ojson --from $indir/flatten-input-2.json json-stringify --jvstack \
  | run_mlr --j2x json-parse then flatten

$path_to_mlr --ijson --ojson --from $indir/flatten-input-2.json json-stringify --no-jvstack \
  | run_mlr --j2x flatten
$path_to_mlr --ijson --ojson --from $indir/flatten-input-2.json json-stringify --no-jvstack \
  | run_mlr --j2x json-parse then flatten

$path_to_mlr --ijson --ojson --from $indir/flatten-input-2.json json-stringify -f req \
  | run_mlr --ijson --oxtab json-parse -f req then flatten
$path_to_mlr --ijson --oxtab --from $indir/flatten-input-2.json json-stringify -f req then flatten \
  | run_mlr --xtab json-parse -f req then flatten

$path_to_mlr --ijson --ojson --from $indir/flatten-input-2.json json-stringify -f req --jvstack \
  | run_mlr --ijson --oxtab json-parse -f req then flatten

$path_to_mlr --ijson --ojson --from $indir/flatten-input-2.json json-stringify -f req --no-jvstack \
  | run_mlr --ijson --oxtab json-parse -f req then flatten
