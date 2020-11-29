run_mlr --from $indir/s.dkvp put '$si = 0; for (i = 0; i < NR; i += 1) { if (i == 2) { $si += 0   } $si += i }'
run_mlr --from $indir/s.dkvp put '$si = 0; for (i = 0; i < NR; i += 1) { if (i == 2) { $si += 100 } $si += i }'
run_mlr --from $indir/s.dkvp put '$si = 0; for (i = 0; i < NR; i += 1) { if (i == 2) { break }      $si += i }'
run_mlr --from $indir/s.dkvp put '$si = 0; for (i = 0; i < NR; i += 1) { if (i == 2) { continue }   $si += i }'

run_mlr --from $indir/s.dkvp --opprint put '
  $si = 0;
  for (i = 0; i < NR; i += 1) {
    if (true) {
      if (i == 2) {
        $si += 0
      }
    }
    $si += i
  }'

run_mlr --from $indir/s.dkvp --opprint put '
  $si = 0;
  for (i = 0; i < NR; i += 1) {
    if (true) {
      if (i == 2) {
        $si += 100
      }
    }
    $si += i
  }'

run_mlr --from $indir/s.dkvp --opprint put '
  $si = 0;
  for (i = 0; i < NR; i += 1) {
    if (true) {
      if (i == 2) {
        break
      }
    }
    $si += i
  }'

run_mlr --from $indir/s.dkvp --opprint put '
  $si = 0;
  for (i = 0; i < NR; i += 1) {
    if (true) {
      if (i == 2) {
        continue
      }
    }
    $si += i
  }'

run_mlr --from $indir/s.dkvp --opprint put '
  $si = 0;
  for (p = 1; p <= 3; p += 1) {
    for (i = 0; i < NR; i += 1) {
      if (i == 2) {
        $si += 0
      }
      $si += i * 10**p
    }
  }'

run_mlr --from $indir/s.dkvp --opprint put '
  $si = 0;
  for (p = 1; p <= 3; p += 1) {
    for (i = 0; i < NR; i += 1) {
      if (i == 2) {
        break
      }
      $si += i * 10**p
    }
  }'

run_mlr --from $indir/s.dkvp --opprint put '
  $si = 0;
  for (p = 1; p <= 3; p += 1) {
    for (i = 0; i < NR; i += 1) {
      if (i == 2) {
        continue
      }
      $si += i * 10**p
    }
  }'
