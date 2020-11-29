run_mlr --j2p --from $indir/typecast.json put '
  $t       = typeof($a);
  $string  = string($a);
  $int     = int($a);
  $float   = float($a);
  $boolean = boolean($a);
' then reorder -f t,a
run_mlr --idkvp --opprint --from $indir/s.dkvp put '
  for (k, v in $*) {
    $["t".k] = typeof(v)
  }
  $tnonesuch = typeof($nonesuch)
'

run_mlr --idkvp --opprint --from $indir/s.dkvp put '
  for (k, v in $*) {
    $["s".k] = string(v)
  }
  $snonesuch = string($nonesuch)
'
