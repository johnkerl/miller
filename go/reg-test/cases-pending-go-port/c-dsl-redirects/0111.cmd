mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit | "tr \[a-z\] \[A-Z\]", @*' reg-test/input/abixy
