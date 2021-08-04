mlr put -q '
  @x_sum += $x;
  end { emit @x_sum }
' ./data/small
