mlr --from ../c/s put '
  func f(map m): map {
    m["x"] *= 200;
    return m;
  }
  $* = f({"a": $a, "x": $x});
'
