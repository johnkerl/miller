run_mlr --opprint --from $indir/abixy put '
  num sum = 0;
  num a = 1;
  for (;;) {
    if (a > NR) {
        break;
    }
    sum += a;
    a += 1
  }
  $z = sum
'

run_mlr --opprint --from $indir/abixy put '
  num sum = 0;
  num a = 1;
  for (;;) {
    if (a > NR) {
        break;
    }
    sum += a;
    a += 1
  }
  $z = sum
'

run_mlr --opprint --from $indir/abixy put '
  num sum = 0;
  for (num a = 1; ;) {
    if (a > NR) {
        break;
    }
    sum += a;
    a += 1
  }
  $z = sum
'

run_mlr --opprint --from $indir/abixy put '
  num sum = 0;
  num a = 1;
  for (; ; a += 1) {
    if (a > NR) {
        break;
    }
    sum += a;
  }
  $z = sum
'

run_mlr --opprint --from $indir/abixy put '
  num sum = 0;
  for (num a=1; ; a += 1) {
    if (a > NR) {
        break;
    }
    sum += a;
  }
  $z = sum
'

run_mlr --opprint --from $indir/abixy put '
  num sum = 0;
  for (num a=1; a <= NR; a += 1) {
    sum += a;
  }
  $z = sum
'

echo x=1 | run_mlr put '
  num a = 100;
  num b = 100;
  for (num a = 200, b = 300; a <= 210; a += 1, b += 1) {
    print "a:".a.",b:".b
  }
  $oa = a;
  $ob = b;
'

run_mlr --from $indir/abixy put '
  for ( ; $x <= 10; $x += 1) {
  }
'

run_mlr --opprint --from $indir/abixy put '
  num a = 100;
  b  = 200;
  @c = 300;
  $d = 400;
  for (num i = 1; i < 1024; i *= 2) {
    a  += 1;
    b  += 1;
    @c += 1;
    $d += 1;
    $oa = a;
    $ob = b;
    $oc = @c;
    $od = $d;
  }
'

run_mlr --opprint --from $indir/abixy put '
  num sum = 0;
  for (num a = 1; a <= 10; a += 1) {
    continue;
    sum += a;
  }
  $z = sum
'

run_mlr --opprint --from $indir/abixy put '
  num sum = 0;
  for (num a = 1; a <= 10; a += 1) {
    sum += a;
    break;
  }
  $z = sum
'

run_mlr --opprint --from $indir/abixy put '
  num sum = 0;
  for (num a = 1; a <= NR; a += 1) {
    if (a == 4 || a == 5) {
      continue;
    }
    if (a == 8) {
        break;
    }
    sum += a;
  }
  $z = sum
'

# Multi-continutation cases

run_mlr --opprint --from $indir/abixy put '
    for ($o1 = 1; ; $o3 = 3) {
        break;
    }
'

run_mlr --opprint --from $indir/abixy put '
    for ($o1 = 1; $o1 < NR; $o1 += 1) {
    }
'

mlr_expect_fail --opprint --from $indir/abixy put '
    for ($o1 = 1, $o2 = 2; $o3 = 3; $o4 = 4) {
    }
'

mlr_expect_fail --opprint --from $indir/abixy put '
    for ($o1 = 1, $o2 = 2; $o3 < 3, $o4 = 4; $o5 = 5) {
        break;
    }
'

run_mlr --opprint --from $indir/abixy put '
    $o4 = 0;
    for ($o1 = 1, $o2 = 2; $o3 = 3, $o4 < 4; $o5 = 5) {
        break;
    }
'
