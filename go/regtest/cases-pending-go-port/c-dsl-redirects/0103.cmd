mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > "regtest/cases-pending-go-port/c-dsl-redirects/0103.out.".$a.$b, (@a, @b), "NR"' regtest/input/abixy
