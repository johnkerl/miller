run_mlr --c2p --from $indir/mod.csv put '
  $add = madd($a, $b, $m);
  $sub = msub($a, $b, $m);
  $mul = mmul($a, $b, $m);
  $exp = mexp($a, $b, $m);
'
