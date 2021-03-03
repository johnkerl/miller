run_mlr --from $indir/abixy put '
  func f(m) {
    dump m;
    sum = 0;
    for (k, v in m) {
      sum += int(k)
    }
    return sum
  }
  @v[$i] = $a;
  $y = f(@v)
'

run_mlr --from $indir/abixy put '
  subr s(m) {
    dump m;
    sum = 0;
    for (k, v in m) {
      sum += int(k)
    }
    @sum = sum;
  }
  @v[$i] = $a;
  call s(@v);
  $y = @sum;
'

run_mlr --from $indir/abixy-het put    'func f(x) {return {"a":x,"b":x**2}}; map o = f($x); $* = o'
run_mlr --from $indir/abixy-het put -q 'func f(x) {return x**2}; var z = f($x); dump z'
run_mlr --from $indir/abixy-het put -q 'func f(x) {map m = {NR:x};return m}; z = f($y); dump z'

mlr_expect_fail --from $indir/abixy put '
  func f(int x): map {
    if (NR==2) {
      return 2
    } else {
      return {}
    }
  }
  $y=f($x)
'

run_mlr --from $indir/abixy put '
  func f(int x): map {
    if (NR==200) {
      return 2
    } else {
      return {}
    }
  }
  $y=f($i)
'

mlr_expect_fail --from $indir/abixy put '
  func f(int x): map {
    if (NR==200) {
      return 2
    } else {
      return {}
    }
  }
  $y=f($x)
'

run_mlr --from $indir/abixy put '
  func f(int x): var {
    if (NR==2) {
      return 2
    } else {
      return {}
    }
  }
  $y=f($i)
'

mlr_expect_fail --from $indir/abixy put '
  func f(int x): var {
    if (NR==2) {
      return 2
    } else {
      return {}
    }
  }
  $y=f($x)
'

mlr_expect_fail --from $indir/abixy put '
  func f(x): int {
    # fall-through return value is absent-null
  }
  $y=f($x)
'

mlr_expect_fail --from $indir/abixy put '
  int a = 1;
  var b = a[2]; # cannot index localvar declared non-map
'

# This one is intended to particularly look at freeing, e.g.  with './reg-test/run --valgrind'.
run_mlr --oxtab --from $indir/abixy-het put '
 $* = mapdiff(
   mapsum($*, {"a": "newval"}),
   {"b": "nonesuch"},
 )
'
