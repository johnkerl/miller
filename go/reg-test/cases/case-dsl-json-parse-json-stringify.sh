# TODO: maybe just git-checkin the intermediates

run_mlr --ijson --ojson --from $indir/flatten-input-2.json put '$req=json_stringify($req)'
run_mlr --ijson --ojson --from $indir/flatten-input-2.json put '$req=json_stringify($req, false)'
run_mlr --ijson --ojson --from $indir/flatten-input-2.json put '$req=json_stringify($req, true)'

$path_to_mlr --ijson --oxtab --from $indir/flatten-input-2.json put '$req=json_stringify($req)' then flatten \
  | run_mlr --ixtab --ojson cat

$path_to_mlr --ijson --oxtab --from $indir/flatten-input-2.json put '$req=json_stringify($req)' then flatten \
  | run_mlr --ixtab --ojson put '$req = json_parse($req)'
