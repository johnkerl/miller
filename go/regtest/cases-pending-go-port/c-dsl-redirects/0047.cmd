mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp >> "regtest/cases-pending-go-port/c-dsl-redirects/0047.out.".$a.$b, @a, "NR"' regtest/input/abixy
