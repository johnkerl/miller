mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > "regtest/cases-pending-go-port/c-dsl-redirects/0082.out.".$a.$b, @a' regtest/input/abixy
