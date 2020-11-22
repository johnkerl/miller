run_mlr uniq    -g a   $indir/abixy-het
run_mlr uniq    -g a,b $indir/abixy-het

run_mlr uniq    -f a   $indir/abixy-het
run_mlr uniq    -f a,b $indir/abixy-het

run_mlr uniq -c -g a   $indir/abixy-het
run_mlr uniq -c -g a,b $indir/abixy-het

run_mlr uniq    -g a   -o foo $indir/abixy-het
run_mlr uniq    -g a,b -o foo $indir/abixy-het

run_mlr uniq    -f a   -o foo $indir/abixy-het
run_mlr uniq    -f a,b -o foo $indir/abixy-het

run_mlr uniq -c -g a   -o foo $indir/abixy-het
run_mlr uniq -c -g a,b -o foo $indir/abixy-het

run_mlr uniq -a           $indir/repeats.dkvp
run_mlr uniq -a -c        $indir/repeats.dkvp
run_mlr uniq -a -c -o foo $indir/repeats.dkvp
run_mlr uniq -a -n        $indir/repeats.dkvp
run_mlr uniq -a -n -o bar $indir/repeats.dkvp
