run_mlr put -q 's=$x; t[$a]=$x; u[$a][$b]=$x; dump s; dump t; dump u; unset s      ; dump s; dump t; dump u' $indir/unset1.dkvp

run_mlr put -q 's=$x; t[$a]=$x; u[$a][$b]=$x; dump s; dump t; dump u; unset t      ; dump s; dump t; dump u' $indir/unset1.dkvp
run_mlr put -q 's=$x; t[$a]=$x; u[$a][$b]=$x; dump s; dump t; dump u; unset t[1]   ; dump s; dump t; dump u' $indir/unset1.dkvp

run_mlr put -q 's=$x; t[$a]=$x; u[$a][$b]=$x; dump s; dump t; dump u; unset u      ; dump s; dump t; dump u' $indir/unset1.dkvp
run_mlr put -q 's=$x; t[$a]=$x; u[$a][$b]=$x; dump s; dump t; dump u; unset u[1]   ; dump s; dump t; dump u' $indir/unset1.dkvp
run_mlr put -q 's=$x; t[$a]=$x; u[$a][$b]=$x; dump s; dump t; dump u; unset u[1][2]; dump s; dump t; dump u' $indir/unset1.dkvp

run_mlr put -q 's=$x; t[$a]=$x; u[$a][$b]=$x; dump s; dump t; dump u; unset s      ; dump s; dump t; dump u' $indir/unset4.dkvp

run_mlr put -q 's=$x; t[$a]=$x; u[$a][$b]=$x; dump s; dump t; dump u; unset t      ; dump s; dump t; dump u' $indir/unset4.dkvp
run_mlr put -q 's=$x; t[$a]=$x; u[$a][$b]=$x; dump s; dump t; dump u; unset t[1]   ; dump s; dump t; dump u' $indir/unset4.dkvp

run_mlr put -q 's=$x; t[$a]=$x; u[$a][$b]=$x; dump s; dump t; dump u; unset u      ; dump s; dump t; dump u' $indir/unset4.dkvp
run_mlr put -q 's=$x; t[$a]=$x; u[$a][$b]=$x; dump s; dump t; dump u; unset u[1]   ; dump s; dump t; dump u' $indir/unset4.dkvp
run_mlr put -q 's=$x; t[$a]=$x; u[$a][$b]=$x; dump s; dump t; dump u; unset u[1][2]; dump s; dump t; dump u' $indir/unset4.dkvp
