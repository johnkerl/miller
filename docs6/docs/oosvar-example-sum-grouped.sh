mlr --opprint put -q '
  @x_sum[$b] += $x;
  end {
    emit @x_sum, "b"
  }
' data/medium
