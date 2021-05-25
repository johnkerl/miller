mlr --opprint put -q '
  @x_sum[$a][$b] += $x;
  @x_count[$a][$b] += 1;
  end{
    for ((a, b), v in @x_sum) {
      @x_mean[a][b] = @x_sum[a][b] / @x_count[a][b];
    }
    emit @x_mean, "a", "b"
  }
' data/medium

