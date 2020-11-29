run_mlr --csvlite --opprint merge-fields    -a p0,min,p29,max,p100,sum -c _in,_out $indir/merge-fields-in-out.csv
run_mlr --csvlite --opprint merge-fields -k -a p0,min,p29,max,p100,sum -c _in,_out $indir/merge-fields-in-out.csv

run_mlr --csvlite --opprint merge-fields -i    -a p0,min,p29,max,p100,sum -c _in,_out $indir/merge-fields-in-out.csv
run_mlr --csvlite --opprint merge-fields -i -k -a p0,min,p29,max,p100,sum -c _in,_out $indir/merge-fields-in-out.csv

run_mlr --csvlite --opprint merge-fields    -a p0,min,p29,max,p100 -c _in,_out $indir/merge-fields-in-out-mixed.csv
run_mlr --csvlite --opprint merge-fields -k -a p0,min,p29,max,p100 -c _in,_out $indir/merge-fields-in-out-mixed.csv

run_mlr --oxtab merge-fields -k -a p0,min,p29,max,p100,sum,count -f a_in_x,a_out_x -o foo $indir/merge-fields-abxy.dkvp
run_mlr --oxtab merge-fields -k -a p0,min,p29,max,p100,sum,count -r in_,out_       -o bar $indir/merge-fields-abxy.dkvp
run_mlr --oxtab merge-fields -k -a p0,min,p29,max,p100,sum,count -c in_,out_              $indir/merge-fields-abxy.dkvp

run_mlr --oxtab merge-fields -i -k -a p0,min,p29,max,p100,sum,count -f a_in_x,a_out_x -o foo $indir/merge-fields-abxy.dkvp
run_mlr --oxtab merge-fields -i -k -a p0,min,p29,max,p100,sum,count -r in_,out_       -o bar $indir/merge-fields-abxy.dkvp
run_mlr --oxtab merge-fields -i -k -a p0,min,p29,max,p100,sum,count -c in_,out_              $indir/merge-fields-abxy.dkvp
