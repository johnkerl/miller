mention empty for-srec
run_mlr --from $indir/abixy put -v 'for(k,v in $*) { }'

mention for-srec without boundvars
run_mlr --from $indir/abixy put -v 'for(k,v in $*) {$nr= NR}'

mention for-srec modifying the srec
run_mlr --from $indir/abixy put -v 'for(k,v in $*) {unset $[k]}; $j = NR'
run_mlr --from $indir/abixy put -v 'for(k,v in $*) {if (k != "x") {unset $[k]}}; $j = NR'
run_mlr --from $indir/abixy --opprint put -v 'for(k,v in $*) {$[k."_orig"]=v; $[k] = "other"}'
run_mlr --from $indir/abixy put -v 'for(k,v in $*) {$[string(v)]=k}'

run_mlr --from $indir/abixy put -v '
  $sum = 0;
  for(k,v in $*) {
    if (k =~ "^[xy]$") {
      $sum += $[k]
    }
  }'

run_mlr --from $indir/abixy put -v '
  $sum = float(0);
  for(k,v in $*) {
    if (k =~ "^[xy]$") {
      $sum += float($[k])
    }
  }'
