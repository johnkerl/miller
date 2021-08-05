mlr --from data/small put -q '
  begin {
    @o = {
      "nrec": 0,
      "nkey": {"numeric":0, "non-numeric":0},
    };
  }
  @o["nrec"] += 1;
  for (k, v in $*) {
    if (is_numeric(v)) {
      @o["nkey"]["numeric"] += 1;
    } else {
      @o["nkey"]["non-numeric"] += 1;
    }
  }
  end {
    dump @o;
  }
'
