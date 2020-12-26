run_mlr --from $indir/abixy --ojson head -n 1 then put '$v[ ["req","host","header"]] = 4'

run_mlr --from $indir/abixy --ojson head -n 1 then put '
  $u[ [1,2,3] ] = 4;
  $v[ [1,2,3] ] = $u[ [1,2,3]];
  $w[ [1,2,3] ] = $u[ [1,2]];
'
