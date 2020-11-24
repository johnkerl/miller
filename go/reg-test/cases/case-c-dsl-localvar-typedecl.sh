run_mlr --from $indir/abixy put '
  str a = $a;
  a = "a:".NR;
  $c = a;
'

mlr_expect_fail --from $indir/abixy put '
  str a = $a;
  a = NR;
  $c = a;
'

mlr_expect_fail --from $indir/abixy put '
  int a = $a;
  a = NR;
  $c = a;
'

run_mlr --from $indir/abixy put '
  int i = $i;
  i = NR;
  $c = a;
'
