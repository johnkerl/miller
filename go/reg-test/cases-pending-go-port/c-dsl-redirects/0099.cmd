mlr head -n 4 then put -q --oxtab '@a[NR]=$a; @b[NR]=$b; emit > "reg-test/cases-pending-go-port/c-dsl-redirects/0099.out.".$a.$b, (@a, @b)' reg-test/input/abixy
