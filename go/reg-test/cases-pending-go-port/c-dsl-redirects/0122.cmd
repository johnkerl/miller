mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > "reg-test/cases-pending-go-port/c-dsl-redirects/0122.out.".$a.$b, all, "NR"' reg-test/input/abixy
