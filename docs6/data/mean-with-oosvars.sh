mlr --opprint put -q '
  @x_sum += $x;
  @x_count += 1;
  end {
    @x_mean = @x_sum / @x_count;
    emit @x_mean
  }
' data/medium

