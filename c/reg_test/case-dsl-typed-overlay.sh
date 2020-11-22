# ----------------------------------------------------------------
announce DSL TYPED OVERLAY

run_mlr put '$y = string($x); $z=$y.$y' $indir/int-float.dkvp
run_mlr put '$z=string($x).string($x)' $indir/int-float.dkvp
run_mlr put '$y = string($x)' then put '$z=$y.$y' $indir/int-float.dkvp
run_mlr put '$a="hello"' then put '$b=$a." world";$z=$x+$y;$c=$b;$a=sub($b,"hello","farewell")' $indir/int-float.dkvp
