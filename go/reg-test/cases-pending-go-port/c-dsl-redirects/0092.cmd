mlr head -n 4 then put -q -o pprint '@a[NR]=$a; @b[NR]=$b; emit > "reg-test/cases-pending-go-port/c-dsl-redirects/0092.out.".$a.$b, @a, "NR"' reg-test/input/abixy
