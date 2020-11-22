mention IF/ELIF WITH ELSE
run_mlr --from $indir/xy40.dkvp put -v '
  if (NR==1) {
    $x = 2;
    $y = 3
  }'

run_mlr --from $indir/xy40.dkvp put -v '
  if (NR == 4) {
    $x = 5;
    $y = 6
  } else {
    $x = 1007;
    $y = 1008
  }'

run_mlr --from $indir/xy40.dkvp put -v '
  if (NR == 9) {
    $x = 10;
    $y = 11
  } elif (NR == 12) {
    $x = 13;
    $y = 14
  } else {
    $x = 1015;
    $y = 1016
  }'

run_mlr --from $indir/xy40.dkvp put -v '
  if (NR == 17) {
    $x = 18;
    $y = 19
  } elif (NR == 20) {
    $x = 21;
    $y = 22
  } elif (NR == 23) {
    $x = 24;
    $y = 25
  } else {
    $x = 1026;
    $y = 1027
  }'

run_mlr --from $indir/xy40.dkvp put -v '
  if (NR == 28) {
    $x = 29;
    $y = 30
  } elif (NR == 31) {
    $x = 32;
    $y = 33
  } elif (NR == 34) {
    $x = 35;
    $y = 36
  } elif (NR == 37) {
    $x = 38;
    $y = 39
  } else {
    $x = 1040;
    $y = 1041
  }'

mention IF/ELIF WITHOUT ELSE
run_mlr --from $indir/xy40.dkvp put -v '
  if (NR == 1) {
    $x = 2;
    $y = 3
  }'

run_mlr --from $indir/xy40.dkvp put -v '
  if (NR == 4) {
    $x = 5;
    $y = 6
  } elif (NR == 7) {
    $x = 8;
    $y = 9
  }'

run_mlr --from $indir/xy40.dkvp put -v '
  if (NR == 10) {
    $x = 11;
    $y = 12
  } elif (NR == 13) {
    $x = 14;
    $y = 15
  } elif (NR == 16) {
    $x = 17;
    $y = 18
  }'

run_mlr --from $indir/xy40.dkvp put -v '
  if (NR == 19) {
    $x = 20;
    $y = 21
  } elif (NR == 22) {
    $x = 23;
    $y = 24
  } elif (NR == 25) {
    $x = 26;
    $y = 37
  } elif (NR == 28) {
    $x = 29;
    $y = 30
  }'
