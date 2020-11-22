
run_mlr --opprint put '$NEW = $[[3]]'     $indir/abixy
run_mlr --opprint put '$NEW = $[[[3]]]'   $indir/abixy

run_mlr --opprint put '$NEW = $[[11]]'    $indir/abixy
run_mlr --opprint put '$NEW = $[[[11]]]'  $indir/abixy

run_mlr --opprint put '$[[3]]   = "NEW"'  $indir/abixy
run_mlr --opprint put '$[[[3]]] = "NEW"'  $indir/abixy

run_mlr --opprint put '$[[11]]   = "NEW"' $indir/abixy
run_mlr --opprint put '$[[[11]]] = "NEW"' $indir/abixy

run_mlr --opprint put '$[[1]] = $[[2]]' $indir/abixy

run_mlr --opprint put '$a     = $[[2]]; unset $["a"]' $indir/abixy
run_mlr --opprint put '$[[1]] = $b;     unset $[[1]]' $indir/abixy
run_mlr --opprint put '$[[1]] = $[[2]]; unset $["a"]' $indir/abixy

# xxx to do -- there is an old bug here with lack of lhmsmv_unset on the typed overlay at unset
run_mlr --opprint put 'unset $c' $indir/abixy
run_mlr --opprint put 'unset $c; $c="new"' $indir/abixy
run_mlr --opprint put '$c=$a.$b; unset $c; $c="new"' $indir/abixy
run_mlr --opprint put '$c=$a.$b; unset $c' $indir/abixy
