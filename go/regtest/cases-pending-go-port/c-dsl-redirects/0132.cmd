mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > stdout, mapexcept($*, "a", "b"), "NR"' regtest/input/abixy
