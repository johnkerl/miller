
run_mlr --icsvlite --odkvp cat $indir/null-fields.csv
run_mlr --inidx --ifs comma --odkvp cat $indir/null-fields.nidx
run_mlr --idkvp --oxtab cat $indir/missings.dkvp

run_mlr --oxtab stats1 -a sum,min,max,antimode,mode -f x          $indir/nullvals.dkvp
run_mlr --oxtab stats1 -a sum,min,max,antimode,mode -f y          $indir/nullvals.dkvp
run_mlr --oxtab stats1 -a sum,min,max,antimode,mode -f z          $indir/nullvals.dkvp
run_mlr --oxtab stats1 -a sum,min,max,antimode,mode -f x,y,z      $indir/nullvals.dkvp
run_mlr --oxtab stats1 -a sum,min,max,antimode,mode -f x     -g a $indir/nullvals.dkvp
run_mlr --oxtab stats1 -a sum,min,max,antimode,mode -f y     -g a $indir/nullvals.dkvp
run_mlr --oxtab stats1 -a sum,min,max,antimode,mode -f z     -g a $indir/nullvals.dkvp
run_mlr --oxtab stats1 -a sum,min,max,antimode,mode -f x,y,z -g a $indir/nullvals.dkvp

run_mlr --opprint merge-fields -a sum,min,max,antimode,mode -f x,y,z -o xyz $indir/nullvals.dkvp
run_mlr --opprint merge-fields -a sum,min,max,antimode,mode -r x,y,z -o xyz $indir/nullvals.dkvp
run_mlr --opprint merge-fields -a sum,min,max,antimode,mode -c x,y,z        $indir/nullvals.dkvp

run_mlr --oxtab stats2 -a cov -f x,y        $indir/nullvals.dkvp
run_mlr --oxtab stats2 -a cov -f x,z        $indir/nullvals.dkvp
run_mlr --oxtab stats2 -a cov -f y,z        $indir/nullvals.dkvp
run_mlr --oxtab stats2 -a cov -f x,y   -g a $indir/nullvals.dkvp
run_mlr --oxtab stats2 -a cov -f x,z   -g a $indir/nullvals.dkvp
run_mlr --oxtab stats2 -a cov -f y,z   -g a $indir/nullvals.dkvp

run_mlr --opprint top    -n 5 -f x          $indir/nullvals.dkvp
run_mlr --opprint top    -n 5 -f y          $indir/nullvals.dkvp
run_mlr --opprint top    -n 5 -f z          $indir/nullvals.dkvp
run_mlr --opprint top    -n 5 -f x,y,z      $indir/nullvals.dkvp
run_mlr --opprint top    -n 5 -f x     -g a $indir/nullvals.dkvp
run_mlr --opprint top    -n 5 -f y     -g a $indir/nullvals.dkvp
run_mlr --opprint top    -n 5 -f z     -g a $indir/nullvals.dkvp
run_mlr --opprint top    -n 5 -f x,y,z -g a $indir/nullvals.dkvp

run_mlr --opprint top -a -n 5 -f x          $indir/nullvals.dkvp
run_mlr --opprint top -a -n 5 -f y          $indir/nullvals.dkvp
run_mlr --opprint top -a -n 5 -f z          $indir/nullvals.dkvp
run_mlr --opprint top -a -n 5 -f x     -g a $indir/nullvals.dkvp
run_mlr --opprint top -a -n 5 -f y     -g a $indir/nullvals.dkvp
run_mlr --opprint top -a -n 5 -f z     -g a $indir/nullvals.dkvp

run_mlr --opprint top    -n 5 -f x          -o foo $indir/nullvals.dkvp
run_mlr --opprint top    -n 5 -f x     -g a -o foo $indir/nullvals.dkvp

run_mlr --opprint top -a -n 5 -f x -o foo      $indir/nullvals.dkvp
run_mlr --opprint top -a -n 5 -f z -o foo -g a $indir/nullvals.dkvp

run_mlr --opprint step -a counter,rsum -f x          $indir/nullvals.dkvp
run_mlr --opprint step -a counter,rsum -f y          $indir/nullvals.dkvp
run_mlr --opprint step -a counter,rsum -f z          $indir/nullvals.dkvp
run_mlr --opprint step -a counter,rsum -f x,y,z      $indir/nullvals.dkvp

run_mlr --opprint step -a counter,rsum -f x     -g a $indir/nullvals.dkvp
run_mlr --opprint step -a counter,rsum -f y     -g a $indir/nullvals.dkvp
run_mlr --opprint step -a counter,rsum -f z     -g a $indir/nullvals.dkvp
run_mlr --opprint step -a counter,rsum -f x,y,z -g a $indir/nullvals.dkvp
