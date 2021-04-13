mlr         --from regtest/input/abixy-het --opprint put -q 'func f(a, b) { return {"ab": a . "_" . b} } emit (f($a, $b), f($x, $y)), "ab"'
