mlr put -q '
  @x_count[$a][$b] += 1;
  @x_sum[$a][$b] += $x;
  end {
    emit @x_count, "a", "b";
    emit @x_sum, "a", "b";
  }
' ../data/small
