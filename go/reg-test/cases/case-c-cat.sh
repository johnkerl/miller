# Intended to be invoked by "." from reg-test/run

run_mlr cat $indir/abixy
run_mlr cat /dev/null

run_mlr cat -n $indir/abixy-het
run_mlr cat -N counter $indir/abixy-het

run_mlr cat -g a,b $indir/abixy-het
run_mlr cat -g a,b $indir/abixy-het

run_mlr cat -g a,b -n $indir/abixy-het
run_mlr cat -g a,b -N counter $indir/abixy-het

echo a,b,c,d,e,f   | run_mlr cat
echo a,b,c,d,e,f,g | run_mlr cat
