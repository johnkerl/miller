run_mlr --from $indir/abixy put '
  func f(int i) {
    return i+3;
  }
  $c = f($i);
'

mlr_expect_fail --from $indir/abixy put '
  func f(int i) {
    return i+3;
  }
  $c = f($x);
'

mlr_expect_fail --from $indir/abixy put '
  func f(num i): int {
    return i+3.45;
  }
  $c = f($x);
'

mlr_expect_fail --from $indir/abixy put '
  func f(num i): int {
    return i+3.45;
  }
  $c = f($i);
'

mlr_expect_fail --from $indir/abixy put '
  func f(num i): int {
    i = "a";
    return 2;
  }
  $c = f($x);
'


run_mlr --from $indir/abixy put '
  subr s(int i) {
    print i+3;
  }
  call s($i);
'

mlr_expect_fail --from $indir/abixy put '
  subr s(int i) {
    print i+3;
  }
  call s($x);
'

mlr_expect_fail --from $indir/abixy put '
  subr s(num i) {
    i = "a";
    print 2;
  }
  call s($x);
'
