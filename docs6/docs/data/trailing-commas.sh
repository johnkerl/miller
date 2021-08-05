mlr --csvlite --from data/a.csv put '
  func f(
    num a,
    num b,
  ): num {
    return a**2 + b**2;
  }
  $* = {
    "s": $a + $b,
    "t": $a - $b,
    "u": f(
      $a,
      $b,
    ),
    "v": NR,
  }
'
