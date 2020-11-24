run_mlr --from $indir/s.dkvp put '
  z = 1;
  if (NR <= 2) {
    z = 2;
  } else {
    z = 3;
  }
  $z = z
'
