
# mention LHS value on first record should result in ZYX for process creation
export indir; run_mlr --from $indir/abixy put -q 'ENV["ZYX"]="CBA".NR; print | ENV["indir"]."/env-assign.sh" , "a is " . $a'
