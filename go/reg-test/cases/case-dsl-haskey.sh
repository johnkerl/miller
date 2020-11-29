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

run_mlr --oxtab --from $indir/s.dkvp head -n 1 then put '
  $hk01 = haskey($x, $a);
  $hk02 = haskey($nonesuch, $a);
  $hk03 = haskey($*, 7);
  $hk04 = haskey($*, "a");
  $hk05 = haskey($*, "nonesuch");
  $hk06 = haskey([10,20,30], 0);
  $hk07 = haskey([10,20,30], 1);
  $hk08 = haskey([10,20,30], 2);
  $hk09 = haskey([10,20,30], 3);
  $hk10 = haskey([10,20,30], 4);
  $hk11 = haskey([10,20,30], -4);
  $hk12 = haskey([10,20,30], -3);
  $hk13 = haskey([10,20,30], -2);
  $hk14 = haskey([10,20,30], -1);
  $hk15 = haskey([10,20,30], "a");
'
