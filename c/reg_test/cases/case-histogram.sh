run_mlr --opprint histogram -f x,y --lo 0 --hi 1 --nbins 20 $indir/small
run_mlr --opprint histogram -f x,y --lo 0 --hi 1 --nbins 20 -o foo_ $indir/small

run_mlr --opprint histogram --nbins 9 --auto -f x,y $indir/ints.dkvp
run_mlr --opprint histogram --nbins 9 --auto -f x,y -o foo_ $indir/ints.dkvp
