# ----------------------------------------------------------------
announce DSL COMMENTS

run_mlr --from $indir/abixy put '
  $s = 1;
  #$t = 2;
  $u = 3;
'

run_mlr --from $indir/abixy filter '
  NR == 1 ||
  #NR == 2 ||
  NR == 3
'

run_mlr --from $indir/abixy put '
  $s = "here is a pound#sign"; # but this is a comment
  #$t = 2;
  $u = 3;
'
