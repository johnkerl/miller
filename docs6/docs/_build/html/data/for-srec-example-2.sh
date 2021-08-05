mlr --from data/small --opprint put '
  $sum1 = 0;
  $sum2 = 0;
  for (k,v in $*) {
    if (is_numeric(v)) {
      $sum1 +=v;
      $sum2 += $[k];
    }
  }
'
