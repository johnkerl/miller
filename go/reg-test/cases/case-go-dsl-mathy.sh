run_mlr --from $indir/s.dkvp --opprint put '$z = max($x,$y)'
run_mlr --from $indir/s.dkvp --opprint put '$z = min($x,$y,$i)'
run_mlr --from $indir/s.dkvp --opprint put '$z = abs($x)'
run_mlr --from $indir/s.dkvp --opprint put '$c = cos(2*M_PI*NR/32); $s = sin(2*M_PI*NR/32)'

run_mlr --from $indir/ten.dkvp --opprint put '$si = sgn($i-5); $sy = sgn($y); $t = atan2($y, $x); $p = $x ** $y; $q = pow($x, $y)'
run_mlr --opprint --from $indir/ten.dkvp put '$q = qnorm(-5 + $i); $r = 5 + invqnorm($q)'
run_mlr --opprint --from $indir/ten.dkvp put '
  $r2 = roundm($i + $x, 2.0);
  $r4 = roundm($i + $x, 4.0);
'
run_mlr --opprint --from $indir/ten.dkvp put '$z=logifit($i,$x,$y)'

run_mlr --from $indir/ten.dkvp --opprint put '$nx = bitcount($x); $ni = bitcount($i)'
