mlr --from data/medium put -q '
  @x_count[$a][$b] += 1;
  @x_sum[$a][$b] += $x;
  end {
    emit (@x_count, @x_sum), "a", "b";
  }
'
