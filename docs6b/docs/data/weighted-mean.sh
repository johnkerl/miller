mlr --from data/medium put -q '
  # Using the y field for weighting in this example
  weight = $y;

  # Using the a field for weighted aggregation in this example
  @sumwx[$a] += weight * $i;
  @sumw[$a] += weight;

  @sumx[$a] += $i;
  @sumn[$a] += 1;

  end {
    map wmean = {};
    map mean  = {};
    for (a in @sumwx) {
      wmean[a] = @sumwx[a] / @sumw[a]
    }
    for (a in @sumx) {
      mean[a] = @sumx[a] / @sumn[a]
    }
    #emit wmean, "a";
    #emit mean, "a";
    emit (wmean, mean), "a";
  }'
