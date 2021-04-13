mlr put -q '@v[NR] = $*; NR == 2 { dump > stderr }' regtest/input/abixy
