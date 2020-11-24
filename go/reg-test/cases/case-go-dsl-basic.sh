run_mlr --from $indir/s.dkvp --idkvp --opprint put '$j=$i+$i'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$j=$i+$x'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$j=$y+$x'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$j=$y+$i'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$j=$y+$y'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$j=$i+$i'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$y=$x*1e6'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$y=$x+1e6'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$y=$x+1'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$y=FILENAME'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$y=FILENUM'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$y=NF'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$y=NR'

run_mlr --idkvp --opprint put '$y=FNR'       $indir/s.dkvp $indir/t.dkvp
run_mlr --idkvp --opprint put '$y=NR'        $indir/s.dkvp $indir/t.dkvp
run_mlr --icsv  --opprint put '$y=FNR'       $indir/s.csv  $indir/t.csv
run_mlr --idkvp --opprint put '$y=FNR+1'     $indir/s.dkvp $indir/t.dkvp
run_mlr --idkvp --opprint put '$y=FNR+$i'    $indir/s.dkvp $indir/t.dkvp
run_mlr --idkvp --opprint put '$y=FNR+3'     $indir/s.dkvp $indir/t.dkvp
run_mlr --idkvp --opprint put '$y=FNR+3+$i'  $indir/s.dkvp $indir/t.dkvp
run_mlr --idkvp --opprint put '$y=$i+$y'     $indir/s.dkvp $indir/t.dkvp
run_mlr --idkvp --opprint put '$y=$i+$x'     $indir/s.dkvp $indir/t.dkvp
run_mlr --idkvp --opprint put '$z=$x+$y'     $indir/s.dkvp $indir/t.dkvp
run_mlr --idkvp --opprint put '$z=$x+$i'     $indir/s.dkvp $indir/t.dkvp
run_mlr --idkvp --opprint put '$z=NR+$i'     $indir/s.dkvp $indir/t.dkvp
run_mlr --idkvp --opprint put '$z=NR-$i'     $indir/s.dkvp $indir/t.dkvp
run_mlr --idkvp --opprint put '$z=4-1'       $indir/s.dkvp $indir/t.dkvp
run_mlr --idkvp --opprint put '$z=NR'        $indir/s.dkvp $indir/t.dkvp
run_mlr --idkvp --opprint put '$z=$i'        $indir/s.dkvp $indir/t.dkvp
run_mlr --idkvp --opprint put '$z=100*NR-$i' $indir/s.dkvp $indir/t.dkvp
run_mlr --idkvp --opprint put '$z=100*$i+$x' $indir/s.dkvp $indir/t.dkvp

run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z=100*$i+$x'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z=100*$i/$x'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z=NR/$i'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z=100/$i'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z=100//$i'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z=100//$x'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z=100.0//$i'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z=100.0//$i'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z=100.0/$i'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z=100.0'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z=100'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z=100.4'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z=1.2'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z=100.0/$i'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z=100.0//$i'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z=0x7fffffffffffffff  + 0x7fffffffffffffff'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z=0x7fffffffffffffff .+ 0x7fffffffffffffff'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z=0x7fffffffffffffff  * 0x7fffffffffffffff'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z=0x7fffffffffffffff .* 0x7fffffffffffffff'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z= (~ $i) + 1'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z= $i == 2'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z= $i != 2'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z= $i >  2'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z= $i >= 2'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z= $i <  2'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z= $i >= 2'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = $i < 5 ? "low" : "high"'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = $i ** 3'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = $x ** 0.5'
