run_mlr put -v '$["a"] = $["b"]; $["x"] = 10 * $["y"]' $indir/abixy
run_mlr --from $indir/abixy put 'while (NF < 256) { $["k".string(NF+1)] = "v".string(NF) }'
