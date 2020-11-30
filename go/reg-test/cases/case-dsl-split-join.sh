# ----------------------------------------------------------------
run_mlr -n put 'end { print joink({}, ",") }'
run_mlr -n put 'end { print joinv({}, ",") }'
run_mlr -n put 'end { print joinkv({}, "=", ",") }'

run_mlr -n put 'end { print joink([], ",") }'
run_mlr -n put 'end { print joinv([], ",") }'
run_mlr -n put 'end { print joinkv([], "=", ",") }'

run_mlr -n put 'end {print joink([1,2,3], ",")}'
run_mlr -n put 'end {print joink({"a":3,"b":4,"c":5}, ",")}'

run_mlr -n put 'end {print joinv([3,4,5], ",")}'
run_mlr -n put 'end {print joinv({"a":3,"b":4,"c":5}, ",")}'

run_mlr -n put 'end {print joinkv([3,4,5], "=", ",")}'
run_mlr -n put 'end {print joinkv({"a":3,"b":4,"c":5}, "=", ",")}'

run_mlr -n put 'end {print splitkv("a=3,b=4,c=5", "=", ",")}'
run_mlr -n put 'end {print splitkvx("a=3,b=4,c=5", "=", ",")}'
run_mlr -n put 'end {print splitnv("3,4,5", ",")}'
run_mlr -n put 'end {print splitnvx("3,4,5", ",")}'

run_mlr -n put 'end {print splitkv("a,b,c", "=", ",")}'
run_mlr -n put 'end {print splitkvx("a,b,c", "=", ",")}'
run_mlr -n put 'end {print splitnv("a,b,c", ",")}'
run_mlr -n put 'end {print splitnvx("a,b,c", ",")}'

run_mlr -n put 'end {print splita("3,4,5", ",")}'
run_mlr -n put 'end {print splitax("3,4,5", ",")}'

run_mlr --ojson --from $indir/s.dkvp put '$keys   = get_keys($*)'
run_mlr --ojson --from $indir/s.dkvp put '$values = get_values($*)'
run_mlr --ojson --from $indir/s.dkvp put '$keys   = get_keys([7,8,9])'
run_mlr --ojson --from $indir/s.dkvp put '$values = get_values([7,8,9])'

run_mlr --ojson --from $indir/s.dkvp put 'begin{@v=[]} @v = append(@v, NR); $v=@v'

