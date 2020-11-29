run_mlr --opprint --from $indir/abixy count-similar -g a
run_mlr --opprint --from $indir/abixy count-similar -g a,b
run_mlr --opprint --from $indir/abixy count-similar -g b,i
run_mlr --opprint --from $indir/abixy count-similar -g a     -o other_name_for_counter
run_mlr --opprint --from $indir/abixy count-similar -g a,b,i -o other_name_for_counter
