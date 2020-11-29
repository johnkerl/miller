run_mlr put 'filter NR > 2' $indir/s.dkvp
mlr_expect_fail filter 'filter NR > 2' $indir/s.dkvp
