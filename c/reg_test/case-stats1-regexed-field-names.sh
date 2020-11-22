run_mlr --from $indir/abixy --opprint stats1 -a sum  -g  a,b      -f  i,x,y

run_mlr --from $indir/abixy --opprint stats1 -a sum --gr '^[a-h]$' --fr '^[i-z]$'
run_mlr --from $indir/abixy --opprint stats1 -a sum  -g  a,b       --fr '^[i-z]$'
run_mlr --from $indir/abixy --opprint stats1 -a sum --gr '^[a-h]$'  -f  i,x,y

run_mlr --from $indir/abixy --opprint stats1 -a sum --gx '^[i-z]$' --fx '^[a-h]$'
run_mlr --from $indir/abixy --opprint stats1 -a sum  -g  a,b       --fx '^[a-h]$'
run_mlr --from $indir/abixy --opprint stats1 -a sum --gx '^[i-z]$'  -f  i,x,y
