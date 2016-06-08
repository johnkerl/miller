mlr --from data/small --opprint put '
  @sum = 0;
  for (k,v in $*) {
    if (isnumeric(v)) {
      @sum += $[k];
    }
  }
  $sum = @sum
'
