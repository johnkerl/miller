
run_mlr --xtab put       '$y     = $pi1 + $pi2' $indir/mixed-types.xtab
run_mlr --xtab put    -F '$y     = $pi1 + $pi2' $indir/mixed-types.xtab
run_mlr --xtab put    -S '$y     = $pi1 . $pi2' $indir/mixed-types.xtab
run_mlr --xtab filter    '999   != $pi1 + $pi2' $indir/mixed-types.xtab
run_mlr --xtab filter -F '999   != $pi1 + $pi2' $indir/mixed-types.xtab
run_mlr --xtab filter -S '"999" != $pi1 . $pi2' $indir/mixed-types.xtab

run_mlr --oxtab put    '$s = $a; $t = $b; $u = 3; $v = 4.0; $ts=typeof($s); $tt=typeof($t); $tu=typeof($u); $tv=typeof($v);' <<EOF
a=1,b=2.0
EOF
run_mlr --oxtab put -F '$s = $a; $t = $b; $u = 3; $v = 4.0; $ts=typeof($s); $tt=typeof($t); $tu=typeof($u); $tv=typeof($v);' <<EOF
a=1,b=2.0
EOF
run_mlr --oxtab put -S '$s = $a; $t = $b; $u = 3; $v = 4.0; $ts=typeof($s); $tt=typeof($t); $tu=typeof($u); $tv=typeof($v);' <<EOF
a=1,b=2.0
EOF

run_mlr --xtab put    '$y=abs($pf1)' $indir/mixed-types.xtab
run_mlr --xtab put    '$y=abs($nf1)' $indir/mixed-types.xtab
run_mlr --xtab put    '$y=abs($zf)'  $indir/mixed-types.xtab
run_mlr --xtab put    '$y=abs($pi1)' $indir/mixed-types.xtab
run_mlr --xtab put    '$y=abs($ni1)' $indir/mixed-types.xtab
run_mlr --xtab put    '$y=abs($zi)'  $indir/mixed-types.xtab

run_mlr --xtab put -F '$y=abs($pf1)' $indir/mixed-types.xtab
run_mlr --xtab put -F '$y=abs($nf1)' $indir/mixed-types.xtab
run_mlr --xtab put -F '$y=abs($zf)'  $indir/mixed-types.xtab
run_mlr --xtab put -F '$y=abs($pi1)' $indir/mixed-types.xtab
run_mlr --xtab put -F '$y=abs($ni1)' $indir/mixed-types.xtab
run_mlr --xtab put -F '$y=abs($zi)'  $indir/mixed-types.xtab

run_mlr --xtab put    '$y=ceil($pf1)' $indir/mixed-types.xtab
run_mlr --xtab put    '$y=ceil($nf1)' $indir/mixed-types.xtab
run_mlr --xtab put    '$y=ceil($zf)'  $indir/mixed-types.xtab
run_mlr --xtab put    '$y=ceil($pi1)' $indir/mixed-types.xtab
run_mlr --xtab put    '$y=ceil($ni1)' $indir/mixed-types.xtab
run_mlr --xtab put    '$y=ceil($zi)'  $indir/mixed-types.xtab

run_mlr --xtab put -F '$y=floor($pf1)' $indir/mixed-types.xtab
run_mlr --xtab put -F '$y=floor($nf1)' $indir/mixed-types.xtab
run_mlr --xtab put -F '$y=floor($zf)'  $indir/mixed-types.xtab
run_mlr --xtab put -F '$y=floor($pi1)' $indir/mixed-types.xtab
run_mlr --xtab put -F '$y=floor($ni1)' $indir/mixed-types.xtab
run_mlr --xtab put -F '$y=floor($zi)'  $indir/mixed-types.xtab

run_mlr --xtab put    '$y=round($pf1)' $indir/mixed-types.xtab
run_mlr --xtab put    '$y=round($nf1)' $indir/mixed-types.xtab
run_mlr --xtab put    '$y=round($zf)'  $indir/mixed-types.xtab
run_mlr --xtab put    '$y=round($pi1)' $indir/mixed-types.xtab
run_mlr --xtab put    '$y=round($ni1)' $indir/mixed-types.xtab
run_mlr --xtab put    '$y=round($zi)'  $indir/mixed-types.xtab

run_mlr --xtab put -F '$y=round($pf1)' $indir/mixed-types.xtab
run_mlr --xtab put -F '$y=round($nf1)' $indir/mixed-types.xtab
run_mlr --xtab put -F '$y=round($zf)'  $indir/mixed-types.xtab
run_mlr --xtab put -F '$y=round($pi1)' $indir/mixed-types.xtab
run_mlr --xtab put -F '$y=round($ni1)' $indir/mixed-types.xtab
run_mlr --xtab put -F '$y=round($zi)'  $indir/mixed-types.xtab

run_mlr --xtab put    '$y=sgn($pf1)' $indir/mixed-types.xtab
run_mlr --xtab put    '$y=sgn($nf1)' $indir/mixed-types.xtab
run_mlr --xtab put    '$y=sgn($zf)'  $indir/mixed-types.xtab
run_mlr --xtab put    '$y=sgn($pi1)' $indir/mixed-types.xtab
run_mlr --xtab put    '$y=sgn($ni1)' $indir/mixed-types.xtab
run_mlr --xtab put    '$y=sgn($zi)'  $indir/mixed-types.xtab

run_mlr --xtab put -F '$y=sgn($pf1)' $indir/mixed-types.xtab
run_mlr --xtab put -F '$y=sgn($nf1)' $indir/mixed-types.xtab
run_mlr --xtab put -F '$y=sgn($zf)'  $indir/mixed-types.xtab
run_mlr --xtab put -F '$y=sgn($pi1)' $indir/mixed-types.xtab
run_mlr --xtab put -F '$y=sgn($ni1)' $indir/mixed-types.xtab
run_mlr --xtab put -F '$y=sgn($zi)'  $indir/mixed-types.xtab

run_mlr --xtab put    '$min=min($pf1,$pf2);$max=max($pf1,$pf2)' $indir/mixed-types.xtab
run_mlr --xtab put    '$min=min($pf1,$pi2);$max=max($pf1,$pi2)' $indir/mixed-types.xtab
run_mlr --xtab put    '$min=min($pi1,$pf2);$max=max($pi1,$pf2)' $indir/mixed-types.xtab
run_mlr --xtab put    '$min=min($pi1,$pi2);$max=max($pi1,$pi2)' $indir/mixed-types.xtab

run_mlr --xtab put -F '$min=min($pf1,$pf2);$max=max($pf1,$pf2)' $indir/mixed-types.xtab
run_mlr --xtab put -F '$min=min($pf1,$pi2);$max=max($pf1,$pi2)' $indir/mixed-types.xtab
run_mlr --xtab put -F '$min=min($pi1,$pf2);$max=max($pi1,$pf2)' $indir/mixed-types.xtab
run_mlr --xtab put -F '$min=min($pi1,$pi2);$max=max($pi1,$pi2)' $indir/mixed-types.xtab

run_mlr --xtab put    '$min=min($pf1,$pf2);$max=max($pf1,$pf2)' $indir/mixed-types.xtab
run_mlr --xtab put    '$min=min($pf1,$pi2);$max=max($pf1,$pi2)' $indir/mixed-types.xtab
run_mlr --xtab put    '$min=min($pi1,$pf2);$max=max($pi1,$pf2)' $indir/mixed-types.xtab
run_mlr --xtab put    '$min=min($pi1,$pi2);$max=max($pi1,$pi2)' $indir/mixed-types.xtab

run_mlr --xtab put -F '$min=min($pf1,$pf2);$max=max($pf1,$pf2)' $indir/mixed-types.xtab
run_mlr --xtab put -F '$min=min($pf1,$pi2);$max=max($pf1,$pi2)' $indir/mixed-types.xtab
run_mlr --xtab put -F '$min=min($pi1,$pf2);$max=max($pi1,$pf2)' $indir/mixed-types.xtab
run_mlr --xtab put -F '$min=min($pi1,$pi2);$max=max($pi1,$pi2)' $indir/mixed-types.xtab

run_mlr --xtab put    '$sum=$pf1+$pf2;$diff=$pf1-$pf2' $indir/mixed-types.xtab
run_mlr --xtab put    '$sum=$pf1+$pi2;$diff=$pf1-$pi2' $indir/mixed-types.xtab
run_mlr --xtab put    '$sum=$pi1+$pf2;$diff=$pi1-$pf2' $indir/mixed-types.xtab
run_mlr --xtab put    '$sum=$pi1+$pi2;$diff=$pi1-$pi2' $indir/mixed-types.xtab

run_mlr --xtab put -F '$sum=$pf1+$pf2;$diff=$pf1-$pf2' $indir/mixed-types.xtab
run_mlr --xtab put -F '$sum=$pf1+$pi2;$diff=$pf1-$pi2' $indir/mixed-types.xtab
run_mlr --xtab put -F '$sum=$pi1+$pf2;$diff=$pi1-$pf2' $indir/mixed-types.xtab
run_mlr --xtab put -F '$sum=$pi1+$pi2;$diff=$pi1-$pi2' $indir/mixed-types.xtab

run_mlr --xtab put    '$prod=$pf1*$pf2;$quot=$pf1/$pf2' $indir/mixed-types.xtab
run_mlr --xtab put    '$prod=$pf1*$pi2;$quot=$pf1/$pi2' $indir/mixed-types.xtab
run_mlr --xtab put    '$prod=$pi1*$pf2;$quot=$pi1/$pf2' $indir/mixed-types.xtab
run_mlr --xtab put    '$prod=$pi1*$pi2;$quot=$pi1/$pi2' $indir/mixed-types.xtab

run_mlr --xtab put -F '$prod=$pf1*$pf2;$quot=$pf1/$pf2' $indir/mixed-types.xtab
run_mlr --xtab put -F '$prod=$pf1*$pi2;$quot=$pf1/$pi2' $indir/mixed-types.xtab
run_mlr --xtab put -F '$prod=$pi1*$pf2;$quot=$pi1/$pf2' $indir/mixed-types.xtab
run_mlr --xtab put -F '$prod=$pi1*$pi2;$quot=$pi1/$pi2' $indir/mixed-types.xtab

run_mlr --xtab put    '$iquot=$pf1//$pf2;$mod=$pf1%$pf2' $indir/mixed-types.xtab
run_mlr --xtab put    '$iquot=$pf1//$pi2;$mod=$pf1%$pi2' $indir/mixed-types.xtab
run_mlr --xtab put    '$iquot=$pi1//$pf2;$mod=$pi1%$pf2' $indir/mixed-types.xtab
run_mlr --xtab put    '$iquot=$pi1//$pi2;$mod=$pi1%$pi2' $indir/mixed-types.xtab

run_mlr --xtab put -F '$iquot=$pf1//$pf2;$mod=$pf1%$pf2' $indir/mixed-types.xtab
run_mlr --xtab put -F '$iquot=$pf1//$pi2;$mod=$pf1%$pi2' $indir/mixed-types.xtab
run_mlr --xtab put -F '$iquot=$pi1//$pf2;$mod=$pi1%$pf2' $indir/mixed-types.xtab
run_mlr --xtab put -F '$iquot=$pi1//$pi2;$mod=$pi1%$pi2' $indir/mixed-types.xtab

run_mlr --xtab put    '$a=roundm($pf1,10.0);$b=roundm($pf1,-10.0)' $indir/mixed-types.xtab
run_mlr --xtab put    '$a=roundm($pf1,10)  ;$b=roundm($pf1,-10)  ' $indir/mixed-types.xtab
run_mlr --xtab put    '$a=roundm($pi1,10.0);$b=roundm($pi1,-10.0)' $indir/mixed-types.xtab
run_mlr --xtab put    '$a=roundm($pi1,10)  ;$b=roundm($pi1,-10)  ' $indir/mixed-types.xtab

run_mlr --oxtab put '$z=$x+$y; $a=3+4; $b="3"."4"; $c="3"+4' <<EOF
x=3,y=4
EOF
