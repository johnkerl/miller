run_cat $indir/flatten-input-1.json

run_mlr --ijson --oxtab --from $indir/flatten-input-1.json put '
  map o = {};
  for (k, v in $*) {
    for (k2, v2 in flatten(k, ".", v)) {
      o[k2] = v2
    }
  }
  $* = o;
'

run_mlr --ijson --oxtab --from $indir/flatten-input-1.json put '$* = flatten("", ".", $*)'

run_mlr --ijson --oxtab --from $indir/flatten-input-1.json put '$* = flatten($*, ".")'

run_mlr --ijson --ojson --no-auto-unflatten --from $indir/flatten-input-1.json put '$a = flatten("a", ".", $a)'

run_mlr --ijson --ojson --no-auto-unflatten --from $indir/flatten-input-1.json put '$b = flatten("b", ".", $b)'

run_mlr --ijson --oxtab --from $indir/flatten-input-2.json put '$* = flatten("", ".", $*)'

run_mlr --ixtab --ojson --no-auto-unflatten --from $indir/unflatten-input.xtab put '$* = unflatten($*, ".")'
