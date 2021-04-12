mlr head -n 4 then put -q --ojson '@a[NR]=$a; @b[NR]=$b; emit > "reg-test/cases-pending-go-port/c-dsl-redirects/0106.out.".$a.$b, (@a, @b), "NR"' reg-test/input/abixy
