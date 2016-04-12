mlr put -q '
  @v[$a][$b]["x_sum"] += 1;
  @v[$a][$b]["x_count"] += $x;
  end {
    emit @v, "a", "b";
  }
' ../data/small
