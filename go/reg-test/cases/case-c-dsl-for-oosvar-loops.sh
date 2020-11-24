run_mlr --opprint -n put -v '
  @o[1][1]["text1"][NR] = $a;
  @o[1][2]["text2"][NR] = $b;
  @o[1][2][$a][$i*100] = $x;
  for((k1,k2),v in @o[1][2]) {
    @n[3][4][k2][k1] = v;
  }
  end {
    emit @n, "a", "b", "c", "d"
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  @o[1][1]["text1"][NR] = $a;
  @o[1][2]["text2"][NR] = $b;
  @o[1][2][$a][$i*100] = $x;
  for((k1,k2),v in @o[1][2]) {
    @n[3][4][k2][k1] = v;
  }
  end {
    emit @n, "a", "b", "c", "d"
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  @sum[$a][$b] += $x;
  @count[$a][$b] += 1;
  end {
    for ((k1, k2), v in @sum) {
      @mean[k1][k2] = @sum[k1][k2] / @count[k1][k2]
    }
    emitp @sum, "a", "b";
    emitp @count, "a", "b";
    emitp @mean, "a", "b"
  }
'

run_mlr --opprint --from $indir/abixy-wide put -q '
  @value["sum"][$a][$b] += $x;
  @value["count"][$a][$b] += 1;
  end {
    for ((k1, k2), v in @value["sum"]) {
      @value["mean"][k1][k2] = @value["sum"][k1][k2] / @value["count"][k1][k2]
    }
    emitp @value, "stat", "a", "b";
  }
'

mlr_expect_fail -n put -v 'for (k, k in $*) {}'

mlr_expect_fail -n put -v 'for (k, k in @*) {}'

mlr_expect_fail -n put -v 'for ((a,a), c in @*) {}'
mlr_expect_fail -n put -v 'for ((a,b), a in @*) {}'
mlr_expect_fail -n put -v 'for ((a,b), b in @*) {}'

mlr_expect_fail -n put -v 'for ((a,a,c), d in @*) {}'
mlr_expect_fail -n put -v 'for ((a,b,a), d in @*) {}'
mlr_expect_fail -n put -v 'for ((a,b,c), a in @*) {}'
mlr_expect_fail -n put -v 'for ((a,b,b), d in @*) {}'
mlr_expect_fail -n put -v 'for ((a,b,c), b in @*) {}'
mlr_expect_fail -n put -v 'for ((a,b,c), c in @*) {}'

run_mlr --from $indir/xyz2 put -q 'func f() { return {"a"."b":"c"."d",3:4}}; for (k,v in f()){print "k=".k.",v=".v}'
run_mlr --from $indir/xyz2 put -q 'for (k,v in {"a"."b":"c"."d",3:"c"}) {print "k=".k.",v=".v}'
run_mlr --from $indir/xyz2 put -q 'o["a"."b"]="c"."d"; for (k,v in o) {print "k=".k.",v=".v}'
run_mlr --from $indir/xyz2 put -q '@o["a"."b"]="c"."d"; for (k,v in @o) {print "k=".k.",v=".v}'
run_mlr --from $indir/xyz2 put 'for (k in $*) { print k}'
run_mlr --from $indir/xyz2 put 'm=$*; for (k in m) { print k}'
