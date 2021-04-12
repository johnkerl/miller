mlr head -n 4 then put -q -o json '@a[NR]=$a; @b[NR]=$b; emit > "reg-test/cases-pending-go-port/c-dsl-redirects/0085.out.".$a.$b, @a' reg-test/input/abixy
