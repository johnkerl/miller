mlr put '
  begin {
    @i_cumu = 0;
  }

  @i_cumu += $i;
  $* = {
    "z": $x + y,
    "KEYFIELD": $a,
    "i": @i_cumu,
    "b": $b,
    "y": $x,
    "x": $y,
  };
' data/small
