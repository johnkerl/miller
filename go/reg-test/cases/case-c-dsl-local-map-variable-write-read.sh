run_mlr --from $indir/xyz345 put '
    $s = a;
    $t = b[$x];
    $u = c[$x][$y];
    $v = d[$x][$y][$z];
'

run_mlr --from $indir/xyz345 put '
    d[$x][$y][$z] = 9;
    $d = d[$x][$y][$z];
'

run_mlr --from $indir/xyz345 put '
    a = 6;
    b[$x] = 7;
    c[$x][$y] = 8;
    d[$x][$y][$z] = 9;

    $a = a;
    $b = b[$x];
    $c = c[$x][$y];
    $d = d[$x][$y][$z];
'

run_mlr --from $indir/xyz345 put '
    a = 6;
    b[$x] = 7;
    c[$x][$y] = 8;
    d[$x][$y][$z] = 9;

    $a0 = a;
    $a1 = a[$x];
    $a2 = a[$x][$y];
    $a3 = a[$x][$y][$z];

    $b0 = b;
    $b1 = b[$x];
    $b2 = b[$x][$y];
    $b3 = b[$x][$y][$z];

    $c0 = c;
    $c1 = c[$x];
    $c2 = c[$x][$y];
    $c3 = c[$x][$y][$z];

    $d0 = d;
    $d1 = d[$x];
    $d2 = d[$x][$y];
    $d3 = d[$x][$y][$z];
'
