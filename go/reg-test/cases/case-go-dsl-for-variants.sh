run_mlr --from $indir/s.dkvp put 'for (@i = 0; @i < NR; @i += 1) { $i += @i }'
run_mlr --from $indir/s.dkvp put 'i=999; for (i = 0; i < NR; i += 1) { $i += i }'
run_mlr --from $indir/s.dkvp put -v 'j = 2; for (i = 0; i < NR; i += 1) { $i += i }'
