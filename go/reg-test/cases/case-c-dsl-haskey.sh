run_mlr --from $indir/abixy-het put '$haskeya = haskey($*, "a")'
run_mlr --from $indir/abixy-het put '$haskey3 = haskey($*, 3)'

run_mlr --from $indir/xyz2 put '$haskeya = haskey({3:4}, "a")'
run_mlr --from $indir/xyz2 put '$haskey3 = haskey({3:4}, 3)'
run_mlr --from $indir/xyz2 put '$haskey3 = haskey({3:4}, 4)'

run_mlr --from $indir/xyz2 put 'o = {3:4}; $haskeya = haskey(o, "a")'
run_mlr --from $indir/xyz2 put 'o = {3:4}; $haskey3 = haskey(o, 3)'
run_mlr --from $indir/xyz2 put 'o = {3:4}; $haskey3 = haskey(o, 4)'

run_mlr --from $indir/xyz2 put '@o = {3:4}; $haskeya = haskey(@o, "a")'
run_mlr --from $indir/xyz2 put '@o = {3:4}; $haskey3 = haskey(@o, 3)'
run_mlr --from $indir/xyz2 put '@o = {3:4}; $haskey3 = haskey(@o, 4)'

run_mlr --from $indir/xyz2 put 'o = "3:4"; $haskeya = haskey(@o, "a")'
run_mlr --from $indir/xyz2 put 'o = "3:4"; $haskey3 = haskey(@o, 3)'
run_mlr --from $indir/xyz2 put 'o = "3:4"; $haskey3 = haskey(@o, 4)'
