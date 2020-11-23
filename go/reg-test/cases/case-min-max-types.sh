# NUMERICS < BOOL < VOID < STRING

# TODO: cmp-matrices need to be fixed to follow the advertised rule for mixed types.

run_mlr --ojson put '
  $min["n"]["n"] = min($n,$n);
  $min["n"]["b"] = min($n,$b);
  $min["n"]["v"] = min($n,$v);
  $min["n"]["s"] = min($n,$s);

  $min["b"]["n"] = min($b,$n);
  $min["b"]["b"] = min($b,$b);
  $min["b"]["v"] = min($b,$v);
  $min["b"]["s"] = min($b,$s);

  $min["v"]["n"] = min($v,$n);
  $min["v"]["b"] = min($v,$b);
  $min["v"]["v"] = min($v,$v);
  $min["v"]["s"] = min($v,$s);

  $min["s"]["n"] = min($s,$n);
  $min["s"]["b"] = min($s,$b);
  $min["s"]["v"] = min($s,$v);
  $min["s"]["s"] = min($s,$s);
' <<EOF
n=1,b=true,v=,s=abc
EOF

run_mlr --ojson put '
  $max["n"]["n"] = max($n,$n);
  $max["n"]["b"] = max($n,$b);
  $max["n"]["v"] = max($n,$v);
  $max["n"]["s"] = max($n,$s);

  $max["b"]["n"] = max($b,$n);
  $max["b"]["b"] = max($b,$b);
  $max["b"]["v"] = max($b,$v);
  $max["b"]["s"] = max($b,$s);

  $max["v"]["n"] = max($v,$n);
  $max["v"]["b"] = max($v,$b);
  $max["v"]["v"] = max($v,$v);
  $max["v"]["s"] = max($v,$s);

  $max["s"]["n"] = max($s,$n);
  $max["s"]["b"] = max($s,$b);
  $max["s"]["v"] = max($s,$v);
  $max["s"]["s"] = max($s,$s);
' <<EOF
n=1,b=true,v=,s=abc
EOF

run_mlr --ojson put '
  $le["n"]["n"] = $n <= $n;
  $le["n"]["b"] = $n <= $b;
  $le["n"]["v"] = $n <= $v;
  $le["n"]["s"] = $n <= $s;

  $le["b"]["n"] = $b <= $n;
  $le["b"]["b"] = $b <= $b;
  $le["b"]["v"] = $b <= $v;
  $le["b"]["s"] = $b <= $s;

  $le["v"]["n"] = $v <= $n;
  $le["v"]["b"] = $v <= $b;
  $le["v"]["v"] = $v <= $v;
  $le["v"]["s"] = $v <= $s;

  $le["s"]["n"] = $s <= $n;
  $le["s"]["b"] = $s <= $b;
  $le["s"]["v"] = $s <= $v;
  $le["s"]["s"] = $s <= $s;
' <<EOF
n=1,b=true,v=,s=abc
EOF

run_mlr --ojson put '
  $ge["n"]["n"] = $n >= $n;
  $ge["n"]["b"] = $n >= $b;
  $ge["n"]["v"] = $n >= $v;
  $ge["n"]["s"] = $n >= $s;

  $ge["b"]["n"] = $b >= $n;
  $ge["b"]["b"] = $b >= $b;
  $ge["b"]["v"] = $b >= $v;
  $ge["b"]["s"] = $b >= $s;

  $ge["v"]["n"] = $v >= $n;
  $ge["v"]["b"] = $v >= $b;
  $ge["v"]["v"] = $v >= $v;
  $ge["v"]["s"] = $v >= $s;

  $ge["s"]["n"] = $s >= $n;
  $ge["s"]["b"] = $s >= $b;
  $ge["s"]["v"] = $s >= $v;
  $ge["s"]["s"] = $s >= $s;
' <<EOF
n=1,b=true,v=,s=abc
EOF
