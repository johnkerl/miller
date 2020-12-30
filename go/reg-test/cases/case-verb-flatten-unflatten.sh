run_mlr --ijson --oxtab flatten        $indir/flatten-input-2.json
run_mlr --ijson --oxtab flatten -s :   $indir/flatten-input-2.json
run_mlr --ijson --oxtab flatten -s .   $indir/flatten-input-2.json

run_mlr --ixtab --ojson unflatten      $indir/unflatten-input.xtab
run_mlr --ixtab --ojson unflatten -s : $indir/unflatten-input.xtab
run_mlr --ixtab --ojson unflatten -s . $indir/unflatten-input.xtab
