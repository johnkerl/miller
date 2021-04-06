# TODO: pending bugfix

run_mlr --from $indir/s.dkvp put -q '@sum += $x; dump'
run_mlr --from $indir/s.dkvp put -q '@sum[$a] += $x; dump'
run_mlr --from $indir/s.dkvp put -q 'begin{@sum=0} @sum += $x; end{dump}'
run_mlr --from $indir/s.dkvp put -q 'begin{@sum={}} @sum[$a] += $x; end{dump}'

run_mlr --from $indir/s.dkvp put -q 'begin{@sum=[3,4]} @sum[1+NR%2] += $x; end{dump}'
run_mlr --from $indir/s.dkvp put -q 'begin{@sum=[0,0]} @sum[1+NR%2] += $x; end{dump}'

# TODO: fix these two
run_mlr --from $indir/s.dkvp put -q 'begin{@sum=[]} @sum[1+NR%2] += $x; end{dump}'
run_mlr --from $indir/s.dkvp put -q 'begin{} @sum[1+(NR%2)] += $x; end{dump}'

run_mlr --from $indir/s.dkvp put 'nr=NR; $nr=nr'

