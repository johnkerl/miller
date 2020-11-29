run_mlr --from $indir/s.dkvp put -q '@sum += $x; @count += 1; dump'
run_mlr --from $indir/s.dkvp put -q '@sum += $x; @count += 1; dump @sum'
run_mlr --from $indir/s.dkvp put -q '@sum += $x; @count += 1; dump @sum, @count'

run_mlr --from $indir/s.dkvp put -q '@sum += $x; @count += 1; print'
run_mlr --from $indir/s.dkvp put -q '@sum += $x; @count += 1; print @sum'
run_mlr --from $indir/s.dkvp put -q '@sum += $x; @count += 1; print @sum, @count'

run_mlr --from $indir/s.dkvp put -q 'print'
run_mlr --from $indir/s.dkvp put -q 'print $x'
run_mlr --from $indir/s.dkvp put -q 'print $x,$y'
run_mlr --from $indir/s.dkvp put -q 'print $x,$y,$i'
