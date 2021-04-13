mlr head -n 4 then put -q -o json '@a[NR]=$a; @b[NR]=$b; emitp > "regtest/cases-pending-go-port/c-dsl-redirects/0041.out.".$a.$b, @a' regtest/input/abixy
