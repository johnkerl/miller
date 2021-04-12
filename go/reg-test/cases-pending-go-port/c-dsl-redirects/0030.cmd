mlr put -q '@v[NR] = $*; NR == 2 { dump > stderr }' reg-test/input/abixy
