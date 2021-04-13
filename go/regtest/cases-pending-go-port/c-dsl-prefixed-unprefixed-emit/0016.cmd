mlr --oxtab put -q '@sum[$a][$b][$a.$b] += $x; end{dump;emit @sum, "a"}'  regtest/input/abixy-wide
