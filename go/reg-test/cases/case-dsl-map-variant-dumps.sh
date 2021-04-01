run_mlr --from $indir/abixy-het put -q 'dump {"a"."b":$a.$b}'
run_mlr --from $indir/abixy-het put -q 'func f(a, b) { return {"a"."b":a.b} } dump f($a, $b)'
