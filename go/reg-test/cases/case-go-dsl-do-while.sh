run_mlr --from $indir/s.dkvp put '$z = 0; while ($z < $i) {$z += 1}'
run_mlr --from $indir/s.dkvp put '$z = 0; do {$z += 1} while ($z < $i)'
run_mlr --from $indir/s.dkvp put '$z = 10; while ($z < $i) {$z += 1}'
run_mlr --from $indir/s.dkvp put '$z = 10; do {$z += 1} while ($z < $i)'
