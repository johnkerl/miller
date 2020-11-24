run_mlr --from $indir/ten.dkvp --opprint put '$z=sec2gmt($i)'
run_mlr --from $indir/ten.dkvp --opprint put '$z=sec2gmt($i, $i-1)'
run_mlr --from $indir/ten.dkvp --opprint put '$z=sec2gmt($i+0.123456789)'
run_mlr --from $indir/ten.dkvp --opprint put '$z=sec2gmt($i+0.123456789,$i-1)'
