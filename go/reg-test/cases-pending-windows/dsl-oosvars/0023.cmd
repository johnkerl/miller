mlr put -q '@s=$x; @t[$a]=$x; @u[$a][$b]=$x; end{dump; unset @u[1][2]; dump}' reg-test/input/unset1.dkvp
