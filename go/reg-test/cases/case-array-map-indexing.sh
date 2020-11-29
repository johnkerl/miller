# Right-hand sides, non-indexed

run_mlr --oxtab --from reg-test/input/abixy head -n 1 then put '
  $_1 = ""; # just for visual output-spacing
  $srec_keyed_by_2        = $[2];
  $srec_name_2            = $[[2]];
  $srec_value_2           = $[[[2]]];
  $_2 = "";
  $dollar_star_keyed_by_2 = $*[2];
  $dollar_star_name_2     = $*[[2]];
  $dollar_star_value_2    = $*[[[2]]];
  $_3 = "";
  mymap                   = {"a":7, "b":8, "c":9};
  $mymap_keyed_by_2       = mymap[2];
  $mymap_name_2           = mymap[[2]];
  $mymap_value_2          = mymap[[[2]]];
  $_4 = "";
  myarray                 = [7, 8, 9];
  $myarray_keyed_by_2     = myarray[2];
  $myarray_name_2         = myarray[[2]];
  $myarray_value_2        = myarray[[[2]]];
'

run_mlr --oxtab --from reg-test/input/abixy head -n 1 then put '
  $_1 = ""; # just for visual output-spacing
  $srec_keyed_by_2        = $[900];
  $srec_name_2            = $[[900]];
  $srec_value_2           = $[[[900]]];
  $_2 = "";
  $dollar_star_keyed_by_2 = $*[900];
  $dollar_star_name_2     = $*[[900]];
  $dollar_star_value_2    = $*[[[900]]];
  $_3 = "";
  mymap                   = {"a":7, "b":8, "c":9};
  $mymap_keyed_by_2       = mymap[900];
  $mymap_name_2           = mymap[[900]];
  $mymap_value_2          = mymap[[[900]]];
  $_4 = "";
  myarray                 = [7, 8, 9];
  $myarray_keyed_by_2     = myarray[900];
  $myarray_name_2         = myarray[[900]];
  $myarray_value_2        = myarray[[[900]]];
'

run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = [$a,$b,$i,$x,$y][1]'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = [$a,$b,$i,$x,$y][-1]'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = [$a,$b,$i,$x,$y][NR]'

run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$up   = $[NF]'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$down = $[-NF]'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$up   = $*[NF]'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$down = $*[-NF]'

run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$up   = $[[NF]]'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$down = $[[-NF]]'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$up   = $*[[NF]]'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$down = $[*[-NF]]'

run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$up   = $[[[NF]]]'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$down = $[[[-NF]]]'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$up   = $*[[[NF]]]'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$down = $*[[[-NF]]]'

# Unsets

run_mlr --opprint put 'unset $c' $indir/abixy
run_mlr --opprint put 'unset $c; $c="new"' $indir/abixy
run_mlr --opprint put '$c=$a.$b; unset $c; $c="new"' $indir/abixy
run_mlr --opprint put '$c=$a.$b; unset $c' $indir/abixy

# Left-hand sides, non-indexed

run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[1] = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[2] = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[5] = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[-1] = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[-2] = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[-5] = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '@idx = NR % 5; @idx = @idx == 0 ? 5 : @idx; $[@idx] = "NEW"'

run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[1]       = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[2]       = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[5]       = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[-1]      = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[-2]      = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[-5]      = "new"'

run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[[1]]     = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[[2]]     = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[[5]]     = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[[-1]]    = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[[-2]]    = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[[-5]]    = "new"'

run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[[[1]]]   = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[[[2]]]   = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[[[5]]]   = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[[[-1]]]  = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[[[-2]]]  = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[[[-5]]]  = "new"'

run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$*[1]      = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$*[2]      = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$*[5]      = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$*[-1]     = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$*[-2]     = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$*[-5]     = "new"'

run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$*[[1]]    = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$*[[2]]    = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$*[[5]]    = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$*[[-1]]   = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$*[[-2]]   = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$*[[-5]]   = "new"'

run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$*[[[1]]]  = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$*[[[2]]]  = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$*[[[5]]]  = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$*[[[-1]]] = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$*[[[-2]]] = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$*[[[-5]]] = "new"'

# Mix of left-hand and right-hand sides

run_mlr --opprint put '$NEW = $[[3]]'     $indir/abixy
run_mlr --opprint put '$NEW = $[[[3]]]'   $indir/abixy

run_mlr --opprint put '$NEW = $[[11]]'    $indir/abixy
run_mlr --opprint put '$NEW = $[[[11]]]'  $indir/abixy

run_mlr --opprint put '$[[3]]   = "NEW"'  $indir/abixy
run_mlr --opprint put '$[[[3]]] = "NEW"'  $indir/abixy

run_mlr --opprint put '$[[11]]   = "NEW"' $indir/abixy
run_mlr --opprint put '$[[[11]]] = "NEW"' $indir/abixy

run_mlr --opprint put '$[[1]] = $[[2]]' $indir/abixy

run_mlr --opprint put '$a     = $[[2]]; unset $["a"]' $indir/abixy
run_mlr --opprint put '$[[1]] = $b;     unset $[[1]]' $indir/abixy
run_mlr --opprint put '$[[1]] = $[[2]]; unset $["a"]' $indir/abixy
