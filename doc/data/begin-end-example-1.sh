mlr put '
  begin { @sum = 0 };
  @x_sum += $x;
  end { emit @x_sum }
' ../data/small
