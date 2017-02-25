mlr -n put --jknquoteint -q '
  begin {
    @myvar = {
      1: 2,
      3: { 4 : 5 },
      6: { 7: { 8: 9 } }
    }
  }
  end {
    for (k, v in @myvar) {
      print
        "key=" . k .
        ",valuetype=" . typeof(v);
    }
  }
'
