run_mlr --from $indir/abixy put '
  for (int i = 0; i < $i; i += 1) {
    $c = i * 10;
  }
'

mlr_expect_fail --from $indir/abixy put '
  for (float i = 0; i < $i; i += 1) {
    $c = i * 10;
  }
'

run_mlr --from $indir/abixy put '
  for (int i = 0; i < $i; i += 1) {
    i += 2;
    $c = i;
  }
'

mlr_expect_fail --from $indir/abixy put '
  for (int i = 0; i < $i; i += 1) {
    i += 1.5;
    $c = i;
  }
'

mlr_expect_fail --from $indir/abixy put '
  for (int i = 0; i < $i; i += 1) {
    i += 1.0;
    $c = i;
  }
'

run_mlr --from $indir/abixy put '
  for (num i = 0; i < $i; i += 1) {
    i += 1.0;
    $c = i;
  }
'
