
mention non-top-level begin/end
mlr_expect_fail put -v 'begin{begin{@x=1}}'
mlr_expect_fail put -v 'true{begin{@x=1}}'
mlr_expect_fail put -v 'end{end{@x=1}}'
mlr_expect_fail put -v 'true{end{@x=1}}'

mention srecs in begin/end
mlr_expect_fail put -v 'begin{$x=1}'
mlr_expect_fail put -v 'begin{@x=$y}'
mlr_expect_fail put -v 'end{$x=1}'
mlr_expect_fail put -v 'end{@x=$y}'
mlr_expect_fail put -v 'begin{@v=$*}'
mlr_expect_fail put -v 'end{$*=@v}'

mlr_expect_fail put -v 'begin{unset $x}'
mlr_expect_fail put -v 'end{unset $x}'
mlr_expect_fail put -v 'begin{unset $*}'
mlr_expect_fail put -v 'end{unset $*}'

mention break/continue outside loop
mlr_expect_fail put -v 'break'
mlr_expect_fail put -v 'continue'

mention oosvars etc. in mlr filter
mlr_expect_fail filter -v 'break'
mlr_expect_fail filter -v 'continue'

mention expanded filter

run_mlr --from $indir/abixy filter '
  begin {
    @avoid = 3
  }
  NR != @avoid
'

run_mlr --from $indir/abixy filter -x '
  begin {
    @avoid = 3
  }
  NR != @avoid
'

run_mlr --from $indir/abixy filter '
  func f(n) {
    return n - 1
  }
  f(NR) == 5
'

run_mlr --from $indir/abixy filter '
  subr s(n) {
    print "NR is ".n
  }
  call s(NR);
  false
'

run_mlr --from $indir/abixy filter '
  int a = 5;
  int b = 7;
  a <= NR && NR <= b
'

mlr_expect_fail --from $indir/abixy filter 'filter false'
mlr_expect_fail --from $indir/abixy filter 'filter false; true'
