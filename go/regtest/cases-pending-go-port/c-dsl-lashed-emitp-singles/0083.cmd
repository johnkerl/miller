mlr         --from regtest/input/abixy-het --opprint put -q '@o = {"ab": $a . "_" . $b}; @p = {"ab": $x . "_" . $y}; emit (@o, @p), "ab"'
