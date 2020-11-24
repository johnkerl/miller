run_mlr --opprint --from $indir/abixy put '
  func f() {
    var a = 1;
    if (NR > 5) {
      a = 2;
    }
    return a;
  }
  func g() {
    var b = 1;
    if (NR > 5) {
      var b = 2;
    }
    return b;
  }
  func h() {
    var a = 1;
    if (NR > 5) {
      a = 2;
      return a;
    }
    return a;
  }
  func i() {
    var b = 1;
    if (NR > 5) {
      var b = 2;
      return b;
    }
    return b;
  }
  $of = f();
  $og = g();
  $oh = h();
  $oi = i();
 '

# test fencing at function entry
run_mlr --opprint --from $indir/abixy put '
  func f() {
    var a = 2;
    var b = 2;
    if (NR > 5) {
      var a = 3;
      b = 3;
    }
    return 1000 * a + b;
  }
  a = 1;
  b = 1;
  $ab = f();
  $oa = a;
  $ob = b;
 '
