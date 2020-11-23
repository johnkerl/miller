# NUMERICS < BOOL < VOID < STRING

run_mlr --ojson put '
  $lt["n"]["n"] = $n < $n;
  $lt["n"]["b"] = $n < $b;
  $lt["n"]["v"] = $n < $v;
  $lt["n"]["s"] = $n < $s;

  $lt["b"]["n"] = $b < $n;
  $lt["b"]["b"] = $b < $b;
  $lt["b"]["v"] = $b < $v;
  $lt["b"]["s"] = $b < $s;

  $lt["v"]["n"] = $v < $n;
  $lt["v"]["b"] = $v < $b;
  $lt["v"]["v"] = $v < $v;
  $lt["v"]["s"] = $v < $s;

  $lt["s"]["n"] = $s < $n;
  $lt["s"]["b"] = $s < $b;
  $lt["s"]["s"] = $s < $s;
  $lt["s"]["v"] = $s < $v;

  $gt["n"]["n"] = $n > $n;
  $gt["n"]["b"] = $n > $b;
  $gt["n"]["v"] = $n > $v;
  $gt["n"]["s"] = $n > $s;

  $gt["b"]["n"] = $b > $n;
  $gt["b"]["b"] = $b > $b;
  $gt["b"]["v"] = $b > $v;
  $gt["b"]["s"] = $b > $s;

  $gt["v"]["n"] = $v > $n;
  $gt["v"]["b"] = $v > $b;
  $gt["v"]["v"] = $v > $v;
  $gt["v"]["s"] = $v > $s;

  $gt["s"]["n"] = $s > $n;
  $gt["s"]["b"] = $s > $b;
  $gt["s"]["v"] = $s > $v;
  $gt["s"]["s"] = $s > $s;
' <<EOF
n=1,b=true,v=,s=abc
EOF

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
  $min["s"]["s"] = min($s,$s);
  $min["s"]["v"] = min($s,$v);

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
