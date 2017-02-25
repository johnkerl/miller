mlr --from data/small --opprint put '
  sum = 0;
  for (k,v in $*) {
    if (is_numeric(v)) {
      sum += $[k];
    }
  }
  $sum = sum
'
