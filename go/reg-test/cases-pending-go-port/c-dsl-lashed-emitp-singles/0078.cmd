mlr         --from reg-test/input/abixy-het --opprint put -q 'func f(a, b) { return {"ab": a . "_" . b} }  o = f($a, $b);  p = f($x, $y); emitp  (o, p)'
