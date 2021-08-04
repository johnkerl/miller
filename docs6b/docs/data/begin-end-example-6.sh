mlr put -q '
  @x_count[$a] += 1;
  @x_sum[$a] += $x;
  end {
    emit @x_count, "a";
    emit @x_sum, "a";
  }
' ./data/small
