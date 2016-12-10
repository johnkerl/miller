mlr --from data/small --opprint put -q '
  # These are executed once per record, which is the first pass.
  # The key is to use NR to index an out-of-stream variable to
  # retain all the x-field values.
  @x_min = min($x, @x_min);
  @x_max = max($x, @x_max);
  @x[NR] = $x;

  # The second pass is in a for-loop in an end-block.
  end {
    for (nr, x in @x) {
      @x_pct[nr] = 100 * (x - @x_min) / (@x_max - @x_min);
    }
    emit (@x, @x_pct), "NR"
  }
'
