run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = true  && true'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = true  && false'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = false && true'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = false && false'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = true  && 4'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = false && 4'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = 4     && true'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = 4     && false'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = false && %%%panic%%%'

run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = true  || true'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = true  || false'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = false || true'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = false || false'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = true  || 4'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = false || 4'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = 4     || true'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = 4     || false'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = true  || %%%panic%%%'

run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = true  ? 4 : %%%panic%%%'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = false ? %%%panic%%% : 5'

run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = $x ?? %%%panic%%%'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = $x ??? %%%panic%%%'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = $x ?? 999'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = $x ??? 999'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = $nonesuch ?? 999'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$y ??= 999'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z ??= 999'

run_mlr --ojson put '
  $a = $a ??  "filla";
  $b = $b ??? "fillb"
' <<EOF
a=1,b=2,c=3
a=,b=,c=3
x=7,y=8,z=9
EOF
