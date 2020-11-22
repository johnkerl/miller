# Intended to be invoked by "." from reg_test/run

run_mlr sort -f a   $indir/abixy
run_mlr sort -r a   $indir/abixy
run_mlr sort -f x   $indir/abixy
run_mlr sort -r x   $indir/abixy
run_mlr sort -nf x  $indir/abixy
run_mlr sort -nr x  $indir/abixy

run_mlr sort -f a,b   $indir/abixy
run_mlr sort -r a,b   $indir/abixy
run_mlr sort -f x,y   $indir/abixy
run_mlr sort -r x,y   $indir/abixy
run_mlr sort -nf x,y  $indir/abixy
run_mlr sort -nr x,y  $indir/abixy

run_mlr sort -f a -nr x $indir/abixy
run_mlr sort -nr y -f a $indir/abixy
run_mlr sort -f a -r b -nf x -nr y $indir/abixy

run_mlr sort -f x $indir/sort-het.dkvp
run_mlr sort -r x $indir/sort-het.dkvp
