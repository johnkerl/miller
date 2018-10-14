mlr --from data/small put '
  func f(map m): map {
    m["x"] *= 200;
    return m;
  }
  $* = f({"a": $a, "x": $x});
'
