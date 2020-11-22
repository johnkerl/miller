# ----------------------------------------------------------------
announce STATS1/STEP INT/FLOAT

run_mlr --opprint step      -a rsum,delta,counter  -f x,y,z $indir/int-float.dkvp
run_mlr --opprint step   -F -a rsum,delta,counter  -f x,y,z $indir/int-float.dkvp
run_mlr --oxtab   stats1    -a min,max,sum,count   -f x,y,z $indir/int-float.dkvp
run_mlr --oxtab   stats1 -F -a min,max,sum,count   -f x,y,z $indir/int-float.dkvp
