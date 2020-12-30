run_mlr --ijson --oxtab flatten        $indir/flatten-input-2.json
run_mlr --ijson --oxtab flatten -s :   $indir/flatten-input-2.json
run_mlr --ijson --oxtab flatten -s .   $indir/flatten-input-2.json

run_mlr --oflatsep @ --from $indir/flatten-input-2.json --ijson --oxtab flatten
run_mlr --oflatsep @ --from $indir/flatten-input-2.json --ijson --oxtab flatten -s %

run_mlr --ixtab --ojson unflatten      $indir/unflatten-input.xtab
run_mlr --ixtab --ojson unflatten -s : $indir/unflatten-input.xtab
run_mlr --ixtab --ojson unflatten -s . $indir/unflatten-input.xtab

run_mlr --ixtab --ojson --iflatsep @ unflatten $indir/unflatten-input-2.xtab

run_mlr --xtab --iflatsep . --oflatsep @ unflatten then flatten $indir/unflatten-input.xtab
