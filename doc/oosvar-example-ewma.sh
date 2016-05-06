mlr --opprint put '
  begin{ @a=0.1 };
  $e = NR==1 ? $x : @a * $x + (1 - @a) * @e;
  @e=$e
' data/small
