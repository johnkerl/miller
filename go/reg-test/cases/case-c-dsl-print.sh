
run_mlr put -q 'print  1; print  "two"; print  $a; print;  print  $i < 4; print  "y is ".string($y); print ""' $indir/abixy
run_mlr put -q 'printn 1; printn "two"; printn $a; printn; printn $i < 4; printn "y is ".string($y); print ""' $indir/abixy

run_mlr put -q 'print  $*; print  $*; print  {}; print' $indir/abixy
run_mlr put -q 'printn $*; printn $*; printn {}; print' $indir/abixy
