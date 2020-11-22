run_mlr --from $indir/abixy put '
  $* = {
    "a": $a,
    "x": $x,
  }
'

run_mlr         --from $indir/xyz345 put 'func f(): int { return 999 } $y=f()'
mlr_expect_fail --from $indir/xyz345 put 'func f(): int { return 999 } $y=f(,)'

run_mlr         --from $indir/xyz345 put 'func f(a,):  int { return a*2 } $y=f(NR)'
mlr_expect_fail --from $indir/xyz345 put 'func f(a,,): int { return a*2 } $y=f(NR)'

run_mlr         --from $indir/xyz345 put 'func f(int a,):  int { return a*2 } $y=f(NR)'
mlr_expect_fail --from $indir/xyz345 put 'func f(int a,,): int { return a*2 } $y=f(NR)'

run_mlr         --from $indir/xyz345 put 'subr s()      { print 999  } call s()'
mlr_expect_fail --from $indir/xyz345 put 'subr s()      { print 999  } call s(,)'

run_mlr         --from $indir/xyz345 put 'subr s(a,)   { print a*2 } call s(NR)'
mlr_expect_fail --from $indir/xyz345 put 'subr s(a,,)  { print a*2 } call s(NR)'

run_mlr         --from $indir/xyz345 put 'subr s(int a,)   { print a*2 } call s(NR)'
mlr_expect_fail --from $indir/xyz345 put 'subr s(int a,,)  { print a*2 } call s(NR)'

run_mlr         --from $indir/xyz345 put '$y=log10($x)'
run_mlr         --from $indir/xyz345 put '$y=log10($x,)'
mlr_expect_fail --from $indir/xyz345 put '$y=log10($x,,)'
mlr_expect_fail --from $indir/xyz345 put '$y=log10(,$x)'
run_mlr         --from $indir/xyz345 put '$y=pow($x,2)'
run_mlr         --from $indir/xyz345 put '$y=pow($x,2,)'
