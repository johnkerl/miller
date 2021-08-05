mlr -n put --jknquoteint -q '
  begin {
    @myvar = {
      1: 2,
      3: { 4 : 5 },
      6: { 7: { 8: 9 } }
    }
  }
  end {
    for ((k1, k2), v in @myvar) {
      print
        "key1=" . k1 .
        ",key2=" . k2 .
        ",valuetype=" . typeof(v);
    }
  }
'
