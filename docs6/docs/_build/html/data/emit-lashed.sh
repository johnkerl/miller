mlr --from data/medium --opprint put -q '
  @x_count[$a][$b] += 1;
  @x_sum[$a][$b] += $x;
  end {
      for ((a, b), _ in @x_count) {
          @x_mean[a][b] = @x_sum[a][b] / @x_count[a][b]
      }
      emit (@x_sum, @x_count, @x_mean), "a", "b"
  }
'
