
run_mlr --from $indir/xyz345 put '
    map a = {};
'

run_mlr --from $indir/xyz345 put '
    map a = {};
    a[1]=2;
    a[3][4]=5;
'

mlr_expect_fail --from $indir/xyz345 put '
    map a = {};
    a=2;
    a[3][4]=5;
'

mlr_expect_fail --from $indir/xyz345 put '
    map a = {};
    a[3][4]=5;
    a=2;
'
