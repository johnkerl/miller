mlr --opprint put '
  $* = {
    "a": $i,
    "i": $a,
    "y": $y * 10,
  }
' data/small
