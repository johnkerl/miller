mlr --oxtab put -q '@sum[$a][$b] += $i; NR == 3 { @x = $x }; NR == 7 { @v = $* }; end {emitp all}' reg-test/input/abixy-het
