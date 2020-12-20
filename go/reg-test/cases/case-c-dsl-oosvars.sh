
run_mlr --opprint put -v 'begin{@ox=0}; $d=$x-@ox; @ox=$x' $indir/abixy
run_mlr --opprint put -v 'begin{@ox="no"}; $d=@ox == "no" ? 1.0 : $x/@ox; @ox=$x' then step -a ratio -f x $indir/abixy
run_mlr --opprint put -v '$d=$x/@ox; @ox=$x' then step -a ratio -f x $indir/abixy
run_mlr --opprint put -v 'begin{@ox="no"}; $d=@ox == "no" ? 1.0 : $x/@ox; @ox=$x' then step -a ratio -f x $indir/abixy
run_mlr --opprint put -v 'begin{@rsum = 0}; @rsum = @rsum + $x; $rsum = @rsum' $indir/abixy
run_mlr --opprint put -v 'begin{@a=0; @b=0; @c=0}; $za=@a; $zb=@b; $zc=@c; $d=@a+@b+@c; @a=@b; @b=@c; @c=$i' $indir/abixy
run_mlr --opprint put -v 'begin {@a=0; @b=0; @c=0}; $za=@a; $zb=@b; $zc=@c; $d=@a+@b+@c; @a=@b; @b=@c; @c=$i' $indir/abixy
run_mlr --opprint put -v 'begin{@ox=0}; $d=$x-@ox; @ox=$x' $indir/abixy

run_mlr put -q '@a=$a; @b=$b; @c=$x; end {emitf @a; emitf @b; emitf @c}' $indir/abixy
run_mlr put -q '@a=$a; @b=$b; @c=$x; end{emitf @a, @b, @c}' $indir/abixy
run_mlr --from $indir/abixy put -q '@a=1;b=2;$c=3;emitf @a,b,$c'

run_mlr --opprint put -v 'begin {@count=0; @sum=0.0}; @count=@count+1; @sum=@sum+$x; end{@mean=@sum/@count; emitf @mean}' $indir/abixy
run_mlr --opprint put -v 'end{@mean=@sum/@count; emitf @mean}; begin {@count=0; @sum=0.0}; @count=@count+1; @sum=@sum+$x' $indir/abixy

run_mlr put -v 'begin{ @a = @b[1] }; $c = @d; @e[$i][2+$j][3] = $4; end{@f[@g[5][@h]] = 6}' /dev/null

run_mlr put '@y[$a]=$y; end{dump}' $indir/abixy

run_mlr stats1 -a sum -f y -g a $indir/abixy
run_mlr put '@y_sum[$a] = $y; end{dump}' $indir/abixy


run_mlr put -q '@s=$x; @t[$a]=$x; @u[$a][$b]=$x; end{dump; unset @s      ; dump}' $indir/unset1.dkvp

run_mlr put -q '@s=$x; @t[$a]=$x; @u[$a][$b]=$x; end{dump; unset @t      ; dump}' $indir/unset1.dkvp
run_mlr put -q '@s=$x; @t[$a]=$x; @u[$a][$b]=$x; end{dump; unset @t[1]   ; dump}' $indir/unset1.dkvp

run_mlr put -q '@s=$x; @t[$a]=$x; @u[$a][$b]=$x; end{dump; unset @u      ; dump}' $indir/unset1.dkvp
run_mlr put -q '@s=$x; @t[$a]=$x; @u[$a][$b]=$x; end{dump; unset @u[1]   ; dump}' $indir/unset1.dkvp
run_mlr put -q '@s=$x; @t[$a]=$x; @u[$a][$b]=$x; end{dump; unset @u[1][2]; dump}' $indir/unset1.dkvp


run_mlr put -q '@s=$x; @t[$a]=$x; @u[$a][$b]=$x; end{dump; unset @s      ; dump}' $indir/unset4.dkvp

run_mlr put -q '@s=$x; @t[$a]=$x; @u[$a][$b]=$x; end{dump; unset @t      ; dump}' $indir/unset4.dkvp
run_mlr put -q '@s=$x; @t[$a]=$x; @u[$a][$b]=$x; end{dump; unset @t[1]   ; dump}' $indir/unset4.dkvp

run_mlr put -q '@s=$x; @t[$a]=$x; @u[$a][$b]=$x; end{dump; unset @u      ; dump}' $indir/unset4.dkvp
run_mlr put -q '@s=$x; @t[$a]=$x; @u[$a][$b]=$x; end{dump; unset @u[1]   ; dump}' $indir/unset4.dkvp
run_mlr put -q '@s=$x; @t[$a]=$x; @u[$a][$b]=$x; end{dump; unset @u[1][2]; dump}' $indir/unset4.dkvp

run_mlr put -q '@s=$x; @t[$a]=$x; @u[$a][$b]=$x; end{dump; unset all;      dump}' $indir/unset4.dkvp
run_mlr put -q '@s=$x; @t[$a]=$x; @u[$a][$b]=$x; end{dump; unset @*;       dump}' $indir/unset4.dkvp

run_mlr put 'unset $x' $indir/unset4.dkvp
run_mlr put 'unset $*; $aaa = 999' $indir/unset4.dkvp
run_mlr --from $indir/abixy put 'x = 1; print "OX=".x; unset x; print "NX=".x'

run_mlr put -q '@{variable.name} += $x; end{emit @{variable.name}}' $indir/abixy
run_mlr put -q '@{variable.name}[$a] += $x; end{emit @{variable.name},"a"}' $indir/abixy

run_mlr put 'for (k,v in $*) { if (k == "i") {unset $[k]}}' $indir/abixy

run_mlr --opprint --from $indir/abixy put -q '
  @output[NR] = $*;
  end {
    for ((nr, k), v in @output) {
      if (nr == 4 || k == "i") {
        unset @output[nr][k]
      }
    }
    emitp @output, "NR", "k"
  }
'
