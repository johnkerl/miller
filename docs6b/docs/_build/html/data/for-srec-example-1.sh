mlr --pprint --from data/for-srec-example.tbl put '
  $sum1 = $f1 + $f2 + $f3;
  $sum2 = 0;
  $sum3 = 0;
  for (key, value in $*) {
    if (key =~ "^f[0-9]+") {
      $sum2 += value;
      $sum3 += $[key];
    }
  }
'
