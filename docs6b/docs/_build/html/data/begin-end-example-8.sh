mlr put '
  begin {
    @num_total = 0;
    @num_positive = 0;
  };
  @num_total += 1;
  $x > 0.0 {
    @num_positive += 1;
    $y = log10($x); $z = sqrt($y)
  };
  end {
    emitf @num_total, @num_positive
  }
' data/put-gating-example-1.dkvp
