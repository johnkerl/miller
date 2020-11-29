run_mlr --opprint step -a rsum,shift,delta,counter -f x,y $indir/abixy
run_mlr --opprint step -a rsum,shift,delta,counter -f x,y -g a $indir/abixy
run_mlr --opprint step -a ewma -d 0.1,0.9 -f x,y -g a $indir/abixy
run_mlr --opprint step -a ewma -d 0.1,0.9 -o smooth,rough -f x,y -g a $indir/abixy

run_mlr --odkvp   step -a rsum,shift,delta,counter -f x,y      $indir/abixy-het
run_mlr --odkvp   step -a rsum,shift,delta,counter -f x,y -g a $indir/abixy-het
run_mlr --opprint step -a ewma -d 0.1,0.9 -f x,y -g a $indir/abixy-het
run_mlr --opprint step -a ewma -d 0.1,0.9 -o smooth,rough -f x,y -g a $indir/abixy-het

run_mlr --icsvlite --opprint step -a from-first -f x      $indir/from-first.csv
run_mlr --icsvlite --opprint step -a from-first -f x -g g $indir/from-first.csv
