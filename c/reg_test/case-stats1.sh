run_mlr --opprint stats1    -a mean,sum,count,min,max,antimode,mode     -f i,x,y $indir/abixy
run_mlr --opprint stats1    -a min,p10,p50,median,antimode,mode,p90,max -f i,x,y $indir/abixy
run_mlr --opprint stats1    -a mean,meaneb,stddev                       -f i,x,y $indir/abixy
run_mlr --oxtab   stats1 -s -a mean,sum,count,min,max,antimode,mode     -f i,x,y $indir/abixy

run_mlr --opprint stats1    -a mean,sum,count,min,max,antimode,mode     -f i,x,y -g a $indir/abixy
run_mlr --opprint stats1    -a min,p10,p50,median,antimode,mode,p90,max -f i,x,y -g a $indir/abixy
run_mlr --opprint stats1    -a mean,meaneb,stddev                       -f i,x,y -g a $indir/abixy
run_mlr --oxtab   stats1 -s -a mean,sum,count,min,max,antimode,mode     -f i,x,y -g a $indir/abixy

run_mlr --opprint stats1    -a mean,sum,count,min,max,antimode,mode     -f i,x,y -g a,b $indir/abixy
run_mlr --opprint stats1    -a min,p10,p50,median,antimode,mode,p90,max -f i,x,y -g a,b $indir/abixy
run_mlr --opprint stats1    -a mean,meaneb,stddev                       -f i,x,y -g a,b $indir/abixy
run_mlr --oxtab   stats1 -s -a mean,sum,count,min,max,antimode,mode     -f i,x,y -g a,b $indir/abixy

run_mlr --oxtab stats1 -a min,p0,p50,p100,max -f x,y,z $indir/string-numeric-ordering.dkvp

run_mlr --oxtab   stats1 -a mean -f x      $indir/abixy-het
run_mlr --oxtab   stats1 -a mean -f x -g a $indir/abixy-het

run_mlr --oxtab   stats1 -a p0,p50,p100 -f x,y    $indir/near-ovf.dkvp
run_mlr --oxtab   stats1 -a p0,p50,p100 -f x,y -F $indir/near-ovf.dkvp
