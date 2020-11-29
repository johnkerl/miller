mention key-only fors

run_mlr --from $indir/abixy-het put -q '
  ab = $a . "_" . $b;
  xy = $x . "_" . $y;
  map o = {};
  o[ab] = xy;
  for (k in o) {
    print "k is " . k;
  }
'

run_mlr --from $indir/abixy-het put -q '
  ab = $a . "_" . $b;
  xy = $x . "_" . $y;
  unset @o;
  @o[ab] = xy;
  for (k in @o) {
    print "k is " . k;
  }
'

run_mlr --from $indir/abixy-het put -q '
  ab = $a . "_" . $b;
  xy = $x . "_" . $y;
  for (k in {ab : xy}) {
    print "k is " . k;
  }
'

run_mlr --from $indir/abixy-het put -q '
  func f(a, b, x, y): map {
    ab = $a . "_" . $b;
    xy = $x . "_" . $y;
    return {ab : xy};
  }
  for (k in f($a, $b, $x, $y)) {
    print "k is " . k;
  }
'

mention key-value fors with scalar values

run_mlr --from $indir/abixy-het put -q '
  ab = $a . "_" . $b;
  xy = $x . "_" . $y;
  map o = {};
  o[ab] = xy;
  for (k, v in o) {
    print "k is " . k . "  v is " . v;
  }
'

run_mlr --from $indir/abixy-het put -q '
  ab = $a . "_" . $b;
  xy = $x . "_" . $y;
  unset @o;
  @o[ab] = xy;
  for (k, v in @o) {
    print "k is " . k . "  v is " . v;
  }
'

run_mlr --from $indir/abixy-het put -q '
  ab = $a . "_" . $b;
  xy = $x . "_" . $y;
  for (k, v in {ab : xy}) {
    print "k is " . k . "  v is " . v;
  }
'

run_mlr --from $indir/abixy-het put -q '
  func f(a, b, x, y): map {
    ab = $a . "_" . $b;
    xy = $x . "_" . $y;
    return {ab : xy};
  }
  for (k, v in f($a, $b, $x, $y)) {
    print "k is " . k . "  v is " . v;
  }
'


mention key-value fors with map values

run_mlr --from $indir/abixy-het put -q '
  ab = $a . "_" . $b;
  xy = $x . "_" . $y;
  map o = {};
  o[ab] = {"foo": xy};
  for (k, v in o) {
    print "k is " . k . "  v is ";
    dump v;
  }
'

run_mlr --from $indir/abixy-het put -q '
  ab = $a . "_" . $b;
  xy = $x . "_" . $y;
  unset @o;
  @o[ab]["foo"] = xy;
  for (k, v in @o) {
    print "k is " . k . "  v is ";
    dump v;
  }
'

run_mlr --from $indir/abixy-het put -q '
  ab = $a . "_" . $b;
  xy = $x . "_" . $y;
  for (k, v in {ab : {"foo": xy}}) {
    print "k is " . k . "  v is ";
    dump v;
  }
'

run_mlr --from $indir/abixy-het put -q '
  func f(a, b, x, y): map {
    ab = $a . "_" . $b;
    xy = $x . "_" . $y;
    return {ab : {"foo" : xy}};
  }
  for (k, v in f($a, $b, $x, $y)) {
    print "k is " . k . "  v is ";
    dump v;
  }
'
