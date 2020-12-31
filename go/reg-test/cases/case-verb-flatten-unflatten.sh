run_mlr --ijson --oxtab flatten      $indir/flatten-input-2.json
run_mlr --ijson --oxtab flatten -s : $indir/flatten-input-2.json
run_mlr --ijson --oxtab flatten -s . $indir/flatten-input-2.json

run_mlr --json flatten -f req $indir/flatten-input-2.json
run_mlr --json flatten -f res $indir/flatten-input-2.json

run_mlr --oflatsep @ --from $indir/flatten-input-2.json --ijson --oxtab flatten
run_mlr --oflatsep @ --from $indir/flatten-input-2.json --ijson --oxtab flatten -s %

run_mlr --ixtab --ojson unflatten      $indir/unflatten-input.xtab
run_mlr --ixtab --ojson unflatten -s : $indir/unflatten-input.xtab
run_mlr --ixtab --ojson unflatten -s . $indir/unflatten-input.xtab

run_mlr --ixtab --ojson --iflatsep @ unflatten $indir/unflatten-input-2.xtab

run_mlr --xtab --iflatsep . --oflatsep @ unflatten then flatten $indir/unflatten-input.xtab

run_mlr --ixtab --ojson unflatten -s . -f req $indir/unflatten-input.xtab
run_mlr --ixtab --ojson unflatten -s . -f res $indir/unflatten-input.xtab
run_mlr --ixtab --ojson unflatten -s . -f req,res $indir/unflatten-input.xtab
run_mlr --ixtab --ojson unflatten -s . -f nonesuch,res $indir/unflatten-input.xtab

# auto-flatten / auto-unflatten
run_mlr --j2x cat $indir/flatten-input-2.json
$path_to_mlr --j2x cat $indir/flatten-input-2.json | run_mlr --x2j cat
