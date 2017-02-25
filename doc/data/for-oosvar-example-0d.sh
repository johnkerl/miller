mlr --from data/small head -n 2 then put -q '
  begin {
    @myvar["nesting-is-too-shallow"] = 1;
    @myvar["nesting-is"]["just-right"] = 2;
    @myvar["nesting-is"]["also-just-right"] = 3;
    @myvar["nesting"]["is"]["deeper"] = 4;
  }
  end {
    for ((k1, k2), v in @*) {
      @terminal[k1][k2] = v
    }
    emit @terminal, "basename", "index1"
  }
'
