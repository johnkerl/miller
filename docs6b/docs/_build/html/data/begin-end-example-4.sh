mlr put -q '
  @x_count += 1;
  @x_sum += $x;
  end {
    emit @x_count;
    emit @x_sum;
  }
' ../data/small
