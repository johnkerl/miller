mlr --opprint put -q '@sum[$a][$b][$a.$b] += $x; end{dump;emit @sum, "a", "b"}'  reg-test/input/abixy-wide
