
run_mlr put '$z = $s . $s'     $indir/null-vs-empty.dkvp
run_mlr put '$z = $s == ""'    $indir/null-vs-empty.dkvp
run_mlr put '$z = $s == $s'    $indir/null-vs-empty.dkvp
run_mlr put '$z = is_empty($s)' $indir/null-vs-empty.dkvp

run_mlr put '$z = $x + $y'      $indir/null-vs-empty.dkvp
run_mlr put '$z = $y + $y'      $indir/null-vs-empty.dkvp
run_mlr put '$z = $x + $nosuch' $indir/null-vs-empty.dkvp
run_mlr put '$t = sub($s,       "ell", "X")' $indir/null-vs-empty.dkvp
run_mlr put '$t = sub($s,       "ell", "")'  $indir/null-vs-empty.dkvp
run_mlr put '$t = sub($nosuch,  "ell", "X")' $indir/null-vs-empty.dkvp
run_mlr put '$t = gsub($s,      "l",   "X")' $indir/null-vs-empty.dkvp
run_mlr put '$t = gsub($s,      "l",   "")'  $indir/null-vs-empty.dkvp
run_mlr put '$t = gsub($nosuch, "l",   "X")' $indir/null-vs-empty.dkvp

mention EMIT
run_mlr put -q '@v=1; @nonesuch       {emit @v}' $indir/abixy
run_mlr put -q '@v=1; @nonesuch==true {emit @v}' $indir/abixy
run_mlr put -q '@v=1; $nonesuch       {emit @v}' $indir/abixy
run_mlr put -q '@v=1; $nonesuch==true {emit @v}' $indir/abixy

mention PLUS
run_mlr --ofs tab put 'begin{};          $xy = $x + $y; $sy = @s + $y; $xt = $x + @t; $st = @s + @t' $indir/typeof.dkvp
run_mlr --ofs tab put 'begin{@s=3};      $xy = $x + $y; $sy = @s + $y; $xt = $x + @t; $st = @s + @t' $indir/typeof.dkvp
run_mlr --ofs tab put 'begin{@t=4};      $xy = $x + $y; $sy = @s + $y; $xt = $x + @t; $st = @s + @t' $indir/typeof.dkvp
run_mlr --ofs tab put 'begin{@s=3;@t=4}; $xy = $x + $y; $sy = @s + $y; $xt = $x + @t; $st = @s + @t' $indir/typeof.dkvp

mention MINUS
run_mlr --ofs tab put 'begin{};          $xy = $x - $y; $sy = @s - $y; $xt = $x - @t; $st = @s - @t' $indir/typeof.dkvp
run_mlr --ofs tab put 'begin{@s=3};      $xy = $x - $y; $sy = @s - $y; $xt = $x - @t; $st = @s - @t' $indir/typeof.dkvp
run_mlr --ofs tab put 'begin{@t=4};      $xy = $x - $y; $sy = @s - $y; $xt = $x - @t; $st = @s - @t' $indir/typeof.dkvp
run_mlr --ofs tab put 'begin{@s=3;@t=4}; $xy = $x - $y; $sy = @s - $y; $xt = $x - @t; $st = @s - @t' $indir/typeof.dkvp

mention TIMES
run_mlr --ofs tab put 'begin{};          $xy = $x * $y; $sy = @s * $y; $xt = $x * @t; $st = @s * @t' $indir/typeof.dkvp
run_mlr --ofs tab put 'begin{@s=3};      $xy = $x * $y; $sy = @s * $y; $xt = $x * @t; $st = @s * @t' $indir/typeof.dkvp
run_mlr --ofs tab put 'begin{@t=4};      $xy = $x * $y; $sy = @s * $y; $xt = $x * @t; $st = @s * @t' $indir/typeof.dkvp
run_mlr --ofs tab put 'begin{@s=3;@t=4}; $xy = $x * $y; $sy = @s * $y; $xt = $x * @t; $st = @s * @t' $indir/typeof.dkvp

mention DIVIDE
run_mlr --ofs tab put 'begin{};          $xy = $x / $y; $sy = @s / $y; $xt = $x / @t; $st = @s / @t' $indir/typeof.dkvp
run_mlr --ofs tab put 'begin{@s=3};      $xy = $x / $y; $sy = @s / $y; $xt = $x / @t; $st = @s / @t' $indir/typeof.dkvp
run_mlr --ofs tab put 'begin{@t=4};      $xy = $x / $y; $sy = @s / $y; $xt = $x / @t; $st = @s / @t' $indir/typeof.dkvp
run_mlr --ofs tab put 'begin{@s=3;@t=4}; $xy = $x / $y; $sy = @s / $y; $xt = $x / @t; $st = @s / @t' $indir/typeof.dkvp

mention INTEGER DIVIDE
run_mlr --ofs tab put 'begin{};          $xy = $x // $y; $sy = @s // $y; $xt = $x // @t; $st = @s // @t' $indir/typeof.dkvp
run_mlr --ofs tab put 'begin{@s=3};      $xy = $x // $y; $sy = @s // $y; $xt = $x // @t; $st = @s // @t' $indir/typeof.dkvp
run_mlr --ofs tab put 'begin{@t=4};      $xy = $x // $y; $sy = @s // $y; $xt = $x // @t; $st = @s // @t' $indir/typeof.dkvp
run_mlr --ofs tab put 'begin{@s=3;@t=4}; $xy = $x // $y; $sy = @s // $y; $xt = $x // @t; $st = @s // @t' $indir/typeof.dkvp

mention REMAINDER
run_mlr --ofs tab put 'begin{};          $xy = $x % $y; $sy = @s % $y; $xt = $x % @t; $st = @s % @t' $indir/typeof.dkvp
run_mlr --ofs tab put 'begin{@s=3};      $xy = $x % $y; $sy = @s % $y; $xt = $x % @t; $st = @s % @t' $indir/typeof.dkvp
run_mlr --ofs tab put 'begin{@t=4};      $xy = $x % $y; $sy = @s % $y; $xt = $x % @t; $st = @s % @t' $indir/typeof.dkvp
run_mlr --ofs tab put 'begin{@s=3;@t=4}; $xy = $x % $y; $sy = @s % $y; $xt = $x % @t; $st = @s % @t' $indir/typeof.dkvp

mention BITWISE AND
run_mlr --ofs tab put 'begin{};          $xy = $x & $y; $sy = @s & $y; $xt = $x & @t; $st = @s & @t' $indir/typeof.dkvp
run_mlr --ofs tab put 'begin{@s=3};      $xy = $x & $y; $sy = @s & $y; $xt = $x & @t; $st = @s & @t' $indir/typeof.dkvp
run_mlr --ofs tab put 'begin{@t=4};      $xy = $x & $y; $sy = @s & $y; $xt = $x & @t; $st = @s & @t' $indir/typeof.dkvp
run_mlr --ofs tab put 'begin{@s=3;@t=4}; $xy = $x & $y; $sy = @s & $y; $xt = $x & @t; $st = @s & @t' $indir/typeof.dkvp

mention BITWISE OR
run_mlr --ofs tab put 'begin{};          $xy = $x | $y; $sy = @s | $y; $xt = $x | @t; $st = @s | @t' $indir/typeof.dkvp
run_mlr --ofs tab put 'begin{@s=3};      $xy = $x | $y; $sy = @s | $y; $xt = $x | @t; $st = @s | @t' $indir/typeof.dkvp
run_mlr --ofs tab put 'begin{@t=4};      $xy = $x | $y; $sy = @s | $y; $xt = $x | @t; $st = @s | @t' $indir/typeof.dkvp
run_mlr --ofs tab put 'begin{@s=3;@t=4}; $xy = $x | $y; $sy = @s | $y; $xt = $x | @t; $st = @s | @t' $indir/typeof.dkvp

mention BITWISE XOR
run_mlr --ofs tab put 'begin{};          $xy = $x ^ $y; $sy = @s ^ $y; $xt = $x ^ @t; $st = @s ^ @t' $indir/typeof.dkvp
run_mlr --ofs tab put 'begin{@s=3};      $xy = $x ^ $y; $sy = @s ^ $y; $xt = $x ^ @t; $st = @s ^ @t' $indir/typeof.dkvp
run_mlr --ofs tab put 'begin{@t=4};      $xy = $x ^ $y; $sy = @s ^ $y; $xt = $x ^ @t; $st = @s ^ @t' $indir/typeof.dkvp
run_mlr --ofs tab put 'begin{@s=3;@t=4}; $xy = $x ^ $y; $sy = @s ^ $y; $xt = $x ^ @t; $st = @s ^ @t' $indir/typeof.dkvp
