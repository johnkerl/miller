mlr --opprint put -q '
  @xes[NR] = $x;
  @all_fields[NR] = $*;
  end {
    dump;

    # Simple loop over key-value pairs in an out-of-stream variable:
    for (key, value in @xes) {
      @item1[key] = value
    }
    emitp @item1, "key";

    # Simple loop over key-value pairs in a subhashmap of an out-of-stream variable:
    for (key, value in @all_fields[4]) {
      @item2[key] = value
    }
    emitp @item2, "key";

    # Loop over multilevel key-value pairs:
    for ((nr, key), value in @all_fields) {
      @item3[nr][key] = value
    }
    emitp @item3, "NR", "key";

    # Loop over all out-of-stream variables with depth two
    for ((k1, k2), v in @*) {
      @item4[string(k1)."_".string(k2)] = v;
    }
    emitp @item4, "key", "value"

  }
' data/small
