



run_mlr -n put 'end { @eq = [1,2,3]       == [1,2,3]       ; print @eq}'
run_mlr -n put 'end { @eq = [1,2,3]       == [1,2,3,4]     ; print @eq}'
run_mlr -n put 'end { @eq = [1,2,3]       == [1,3,3]       ; print @eq}'
run_mlr -n put 'end { @eq = ["a",2,3]     == [1,2,3]       ; print @eq}'
run_mlr -n put 'end { @eq = []            == {}            ; print @eq}'
run_mlr -n put 'end { @eq = {}            == {}            ; print @eq}'
run_mlr -n put 'end { @eq = {"a":1}       == {"a":1}       ; print @eq}'
run_mlr -n put 'end { @eq = {"a":1}       == {"a":2}       ; print @eq}'
run_mlr -n put 'end { @eq = {"a":1}       == {"b":1}       ; print @eq}'
run_mlr -n put 'end { @eq = {"a":1,"b":2} == {"b":2}       ; print @eq}'
run_mlr -n put 'end { @eq = {"a":1,"b":2} == {"a":1,"b":2} ; print @eq}'
run_mlr -n put 'end { @eq = {"b":2,"a":1} == {"a":1,"b":2} ; print @eq}'

run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = true  && true'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = true  && false'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = false && true'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = false && false'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = true  && 4'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = false && 4'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = 4     && true'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = 4     && false'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = false && %%%panic%%%'

run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = true  || true'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = true  || false'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = false || true'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = false || false'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = true  || 4'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = false || 4'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = 4     || true'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = 4     || false'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = true  || %%%panic%%%'

run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = true  ? 4 : %%%panic%%%'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = false ? %%%panic%%% : 5'

run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = $x ?? %%%panic%%%'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = $x ?? 999'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = $nonesuch ?? 999'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$y ??= 999'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z ??= 999'

run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = "abc\"def\"ghi"'

run_mlr --from $indir/s.dkvp --idkvp --opprint put -v '$i += 2'
run_mlr --from $indir/s.dkvp --idkvp --opprint put -v '$i *= 2'
run_mlr --from $indir/s.dkvp --idkvp --opprint put -v '$i /= 2'
run_mlr --from $indir/s.dkvp --idkvp --opprint put -v '$i |= 2'
run_mlr --from $indir/s.dkvp --idkvp --opprint put -v '$j = true; $j &&= $i < 2'

run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = [$a,$b,$i,$x,$y][1]'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = [$a,$b,$i,$x,$y][-1]'
run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = [$a,$b,$i,$x,$y][NR]'

run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = {"a":$a,"b":$b,"i":$i,"x":$x,"y":$y}["b"]'

run_mlr --from $indir/s.dkvp --from $indir/t.dkvp --ojson put '$z=[1,2,[NR,[FILENAME,5],$x*$y]]'

echo '{"x":1}'                                           | run_mlr --json cat
echo '{"x":[1,2,3]}'                                     | run_mlr --json cat
echo '{"x":[1,[2,3,4],5]}'                               | run_mlr --json cat
echo '{"x":[1,[2,[3,4,5],6],7]}'                         | run_mlr --json cat

echo '{"x":{}}'                                          | run_mlr --json cat
echo '{"x":{"a":1,"b":2,"c":3}}'                         | run_mlr --json cat
echo '{"x":{"a":1,"b":{"c":3,"d":4,"e":5},"f":6}}'       | run_mlr --json cat

echo '{"x":{},"y":1}'                                    | run_mlr --json cat
echo '{"x":{"a":1,"b":2,"c":3},"y":4}'                   | run_mlr --json cat
echo '{"x":{"a":1,"b":{"c":3,"d":4,"e":5},"f":6},"y":7}' | run_mlr --json cat

echo '{"x":1}'                                           | run_mlr --json cat | run_mlr --json cat
echo '{"x":[1,2,3]}'                                     | run_mlr --json cat | run_mlr --json cat
echo '{"x":[1,[2,3,4],5]}'                               | run_mlr --json cat | run_mlr --json cat
echo '{"x":[1,[2,[3,4,5],6],7]}'                         | run_mlr --json cat | run_mlr --json cat

echo '{"x":{}}'                                          | run_mlr --json cat | run_mlr --json cat
echo '{"x":{"a":1,"b":2,"c":3}}'                         | run_mlr --json cat | run_mlr --json cat
echo '{"x":{"a":1,"b":{"c":3,"d":4,"e":5},"f":6}}'       | run_mlr --json cat | run_mlr --json cat

echo '{"x":{},"y":1}'                                    | run_mlr --json cat | run_mlr --json cat
echo '{"x":{"a":1,"b":2,"c":3},"y":4}'                   | run_mlr --json cat | run_mlr --json cat
echo '{"x":{"a":1,"b":{"c":3,"d":4,"e":5},"f":6},"y":7}' | run_mlr --json cat | run_mlr --json cat

run_mlr --from $indir/s.dkvp --idkvp --ojson put '$z = $*["a"]'
run_mlr --from $indir/s.dkvp --idkvp --ojson put '$z = $*'

run_mlr --from $indir/s.dkvp --idkvp --ojson put '$* = {"s": 7, "t": 8}'
run_mlr --from $indir/s.dkvp --idkvp --ojson put '$*["st"] = 78'
run_mlr --from $indir/s.dkvp --idkvp --ojson put '$*["a"] = 78'
run_mlr --from $indir/s.dkvp --idkvp --ojson put '$*["a"] = {}'
run_mlr --from $indir/s.dkvp --idkvp --opprint put -v '$new = $["a"]'
run_mlr --from $indir/s.dkvp --idkvp --opprint put -v '$["new"] = $a'
run_mlr --from $indir/s.dkvp --idkvp --opprint put -v '${new} = $a . $b'
run_mlr --from $indir/s.dkvp --idkvp --opprint put -v '$new = ${a} . ${b}'

run_mlr --from $indir/s.dkvp --idkvp --opprint put '@tmp = $a . $b; $ab = @tmp'
run_mlr --ojson --from $indir/s.dkvp put '@curi=$i; $curi = @curi; $lagi=@lagi; @lagi=$i'
run_mlr --from $indir/s.dkvp --ojson put '$z["abc"]["def"]["ghi"]=NR'

run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$up   = $[NR]'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$down = $[-NR]'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$up   = $*[NR]'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$down = $*[-NR]'

run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[1] = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[2] = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[5] = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[-1] = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[-2] = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[-5] = "new"'
# expect fail run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[1] = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '@idx = NR % 5; @idx = @idx == 0 ? 5 : @idx; $[@idx] = "NEW"'

run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$*[1] = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$*[2] = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$*[5] = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$*[-1] = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$*[-2] = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$*[-5] = "new"'
# expect fail run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '$[1] = "new"'
run_mlr --opprint --from $indir/s.dkvp --from $indir/t.dkvp put '@idx = NR % 5; @idx = @idx == 0 ? 5 : @idx; $*[@idx] = "NEW"'

run_mlr --json put '$a=$a[2]["b"][1]' $indir/nested.json

run_mlr --ojson --from $indir/2.dkvp put '$abc[FILENAME] = "def"'
run_mlr --ojson --from $indir/2.dkvp put '$abc[NR] = "def"'
run_mlr --ojson --from $indir/2.dkvp put '$abc[FILENAME][NR] = "def"'
run_mlr --ojson --from $indir/2.dkvp put '$abc[NR][FILENAME] = "def"'

run_mlr --ojson --from $indir/2.dkvp put '@abc[FILENAME] = "def"; $ghi = @abc'
run_mlr --ojson --from $indir/2.dkvp put '@abc[NR] = "def"; $ghi = @abc'
run_mlr --ojson --from $indir/2.dkvp put '@abc[FILENAME][NR] = "def"; $ghi = @abc'
run_mlr --ojson --from $indir/2.dkvp put '@abc[NR][FILENAME] = "def"; $ghi = @abc'

run_mlr --from $indir/2.dkvp --ojson put '@a = 3; $new=@a'
run_mlr --from $indir/2.dkvp --ojson put '@a = 3; @a[1]=4; $new=@a'
run_mlr --from $indir/2.dkvp --ojson put '@a = 3; @a[1]=4;@a[1][1]=5; $new=@a'

run_mlr --from $indir/2.dkvp --ojson put '@a = 3; @a["x"]=4; $new=@a'
run_mlr --from $indir/2.dkvp --ojson put '@a = 3; @a["x"]=4;@a["x"]["x"]=5; $new=@a'

run_mlr -n put -v '$z=max()'
run_mlr -n put -v '$z=max(1)'
run_mlr -n put -v '$z=max(1,)'
run_mlr -n put -v '$z=max(1,2)'
run_mlr -n put -v '$z=max(1,2,)'
run_mlr -n put -v '$z=max(1,2,3)'
run_mlr -n put -v '$z=max(1,2,3,)'

run_mlr --from $indir/s.dkvp --opprint put '$z = max($x,$y)'
run_mlr --from $indir/s.dkvp --opprint put '$z = min($x,$y,$i)'
run_mlr --from $indir/s.dkvp --opprint put '$z = abs($x)'
run_mlr --from $indir/s.dkvp --opprint put '$c = cos(2*M_PI*NR/32); $s = sin(2*M_PI*NR/32)'

run_mlr --from $indir/ten.dkvp --opprint put '$si = sgn($i-5); $sy = sgn($y); $t = atan2($y, $x); $p = $x ** $y; $q = pow($x, $y)'
run_mlr --opprint --from $indir/ten.dkvp put '$q = qnorm(-5 + $i); $r = 5 + invqnorm($q)'
run_mlr --opprint --from $indir/ten.dkvp put '
  $r2 = roundm($i + $x, 2.0);
  $r4 = roundm($i + $x, 4.0);
'
run_mlr --opprint --from $indir/ten.dkvp put '$z=logifit($i,$x,$y)'

run_mlr --from $indir/ten.dkvp --opprint put '$nx = bitcount($x); $ni = bitcount($i)'

run_mlr --from $indir/s.dkvp --opprint put 'filter NR > 2'
run_mlr --from $indir/s.dkvp --opprint filter 'NR > 2'
run_mlr --from $indir/s.dkvp --opprint filter -x 'NR > 2'

# run_mlr --from $indir/s.dkvp --opprint put -q '@sum += $i; emit {"sum": @sum}'
run_mlr --from $indir/s.dkvp --opprint put -q '@sum[$a] += $i; emit {"sum": @sum}'

# ----------------------------------------------------------------
run_mlr --opprint --from $indir/s.dkvp put          '@sum += 1; $z=@sum'
run_mlr --opprint --from $indir/s.dkvp put -s sum=0 '@sum += 1; $z=@sum'
run_mlr --opprint --from $indir/s.dkvp put -s sum=8 '@sum += 1; $z=@sum'

run_mlr --opprint --from $indir/ten.dkvp put -s a=0 -s b=1 '
  @c = @a + @b;
  $fa = @a;
  $fb = @b;
  $fc = @c;
  @a = @b;
  @b = @c;
'

run_mlr --opprint --from $indir/ten.dkvp put -e 'begin {@a=0}' -e 'begin {@b=1}' -e '
  @c = @a + @b;
  $fa = @a;
  $fb = @b;
  $fc = @c;
  @a = @b;
  @b = @c;
'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '@sum += $x; @count += 1; dump'
run_mlr --from $indir/s.dkvp put -q '@sum += $x; @count += 1; dump @sum'
run_mlr --from $indir/s.dkvp put -q '@sum += $x; @count += 1; dump @sum, @count'

run_mlr --from $indir/s.dkvp put -q '@sum += $x; @count += 1; print'
run_mlr --from $indir/s.dkvp put -q '@sum += $x; @count += 1; print @sum'
run_mlr --from $indir/s.dkvp put -q '@sum += $x; @count += 1; print @sum, @count'

run_mlr --from $indir/s.dkvp put -q 'print'
run_mlr --from $indir/s.dkvp put -q 'print $x'
run_mlr --from $indir/s.dkvp put -q 'print $x,$y'
run_mlr --from $indir/s.dkvp put -q 'print $x,$y,$i'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '@sum += $x; dump'
run_mlr --from $indir/s.dkvp put -q '@sum[$a] += $x; dump'
run_mlr --from $indir/s.dkvp put -q 'begin{@sum=0} @sum += $x; end{dump}'
run_mlr --from $indir/s.dkvp put -q 'begin{@sum={}} @sum[$a] += $x; end{dump}'

run_mlr --from $indir/s.dkvp put -q 'begin{@sum=[3,4]} @sum[1+NR%2] += $x; end{dump}'
run_mlr --from $indir/s.dkvp put -q 'begin{@sum=[0,0]} @sum[1+NR%2] += $x; end{dump}'
# TODO: fix this
run_mlr --from $indir/s.dkvp put -q 'begin{@sum=[]} @sum[1+NR%2] += $x; end{dump}'
# TODO: fix this
run_mlr --from $indir/s.dkvp put -q 'begin{} @sum[1+(NR%2)] += $x; end{dump}'

run_mlr --from $indir/s.dkvp put 'if (NR == 1) { $z = 100 }'
run_mlr --from $indir/s.dkvp put 'if (NR == 1) { $z = 100 } else { $z = 900 }'
run_mlr --from $indir/s.dkvp put 'if (NR == 1) { $z = 100 } elif (NR == 2) { $z = 200 }'
run_mlr --from $indir/s.dkvp put 'if (NR == 1) { $z = 100 } elif (NR == 2) { $z = 200 } else { $z = 900 }'
run_mlr --from $indir/s.dkvp put 'if (NR == 1) { $z = 100 } elif (NR == 2) { $z = 200 } elif (NR == 3) { $z = 300 } else { $z = 900 }'

run_mlr --from $indir/s.dkvp put 'NR == 2 { $z = 100 }'
run_mlr --from $indir/s.dkvp put 'NR != 2 { $z = 100 }'

echo x=eeee | run_mlr put '$y=ssub($x, "e", "X")'
echo x=eeee | run_mlr put '$y=sub($x, "e", "X")'
echo x=eeee | run_mlr put '$y=gsub($x, "e", "X")'

run_mlr --opprint --from $indir/s.dkvp put '$z = truncate($a, -1)'
run_mlr --opprint --from $indir/s.dkvp put '$z = truncate($a, 0)'
run_mlr --opprint --from $indir/s.dkvp put '$z = truncate($a, 1)'
run_mlr --opprint --from $indir/s.dkvp put '$z = truncate($a, 2)'
run_mlr --opprint --from $indir/s.dkvp put '$z = truncate($a, 3)'
run_mlr --opprint --from $indir/s.dkvp put '$z = truncate($a, 4)'

run_mlr --from $indir/s.dkvp head -n 2 then put -q 'for (k in $*) { emit { "foo" : "bar" } }'
run_mlr --from $indir/s.dkvp head -n 2 then put -q 'for (k in $*) { emit { "foo" : k } }'
run_mlr --from $indir/s.dkvp head -n 2 then put -q 'for (k in $*) { emit { k: "bar" } }'
run_mlr --from $indir/s.dkvp head -n 2 then put -q 'for (k in $*) { emit { k : k } }'

run_mlr --from $indir/s.dkvp head -n 2 then put -q 'for (k,v in $*) { emit { "foo" : "bar" } }'
run_mlr --from $indir/s.dkvp head -n 2 then put -q 'for (k,v in $*) { emit { "foo" : v } }'
run_mlr --from $indir/s.dkvp head -n 2 then put -q 'for (k,v in $*) { emit { k: "bar" } }'
run_mlr --from $indir/s.dkvp head -n 2 then put -q 'for (k,v in $*) { emit { k : v } }'

run_mlr -n put -v 'for (k             in @*) {}'
run_mlr -n put -v 'for (k, v          in @*) {}'
run_mlr -n put -v 'for ((k1,k2),    v in @*) {}'
run_mlr -n put -v 'for ((k1,k2,k3), v in @*) {}'

run_mlr --from $indir/s.dkvp put '$z = 0; while ($z < $i) {$z += 1}'
run_mlr --from $indir/s.dkvp put '$z = 0; do {$z += 1} while ($z < $i)'
run_mlr --from $indir/s.dkvp put '$z = 10; while ($z < $i) {$z += 1}'
run_mlr --from $indir/s.dkvp put '$z = 10; do {$z += 1} while ($z < $i)'

run_mlr --from $indir/s.dkvp head -n 1 then put -q 'for (e in [3,4,5]) { emit { "foo" : "bar" } }'
run_mlr --from $indir/s.dkvp head -n 1 then put -q 'for (e in [3,4,5]) { emit { "foo" : e } }'

run_mlr --from $indir/s.dkvp head -n 1 then put -q 'for (i,e in [3,4,5]) { emit { "foo" : "bar" } }'
run_mlr --from $indir/s.dkvp head -n 1 then put -q 'for (i,e in [3,4,5]) { emit { "foo" : i } }'
run_mlr --from $indir/s.dkvp head -n 1 then put -q 'for (i,e in [3,4,5]) { emit { "foo" : e } }'

run_mlr --from $indir/s.dkvp put 'nr=NR; $nr=nr'

run_mlr --from $indir/s.dkvp put '
  z = 1;
  if (NR <= 2) {
    z = 2;
  } else {
    z = 3;
  }
  $z = z
'

run_mlr --from $indir/s.dkvp put 'for (@i = 0; @i < NR; @i += 1) { $i += @i }'
run_mlr --from $indir/s.dkvp put 'i=999; for (i = 0; i < NR; i += 1) { $i += i }'
run_mlr --from $indir/s.dkvp put -v 'j = 2; for (i = 0; i < NR; i += 1) { $i += i }'

run_mlr --from $indir/ten.dkvp --opprint put '$z=sec2gmt($i)'
run_mlr --from $indir/ten.dkvp --opprint put '$z=sec2gmt($i, $i-1)'
run_mlr --from $indir/ten.dkvp --opprint put '$z=sec2gmt($i+0.123456789)'
run_mlr --from $indir/ten.dkvp --opprint put '$z=sec2gmt($i+0.123456789,$i-1)'

echo 'x= a     b '| run_mlr --ojson put '$y = strip($x)'
echo 'x= a     b '| run_mlr --ojson put '$y = lstrip($x)'
echo 'x= a     b '| run_mlr --ojson put '$y = rstrip($x)'
echo 'x= a     b '| run_mlr --ojson put '$y = collapse_whitespace($x)'
echo 'x= a     b '| run_mlr --ojson put '$y = clean_whitespace($x)'

run_mlr --from $indir/s.dkvp put '$z = strlen($a)'

echo "x=abcdefg" | run_mlr put '$y = substr($x, 0, 0)'
echo "x=abcdefg" | run_mlr put '$y = substr($x, 0, 7)'
echo "x=abcdefg" | run_mlr put '$y = substr($x, 1, 7)'
echo "x=abcdefg" | run_mlr put '$y = substr($x, 1, 6)'
echo "x=abcdefg" | run_mlr put '$y = substr($x, 2, 5)'
echo "x=abcdefg" | run_mlr put '$y = substr($x, 2, 3)'
echo "x=abcdefg" | run_mlr put '$y = substr($x, 3, 3)'
echo "x=abcdefg" | run_mlr put '$y = substr($x, 4, 3)'
echo "x=1234567" | run_mlr put '$y = substr($x, 2, 5)'

echo "x=1,y=abcdefg,z=3" | run_mlr put '$n = length($x)'
echo "x=1,y=abcdefg,z=3" | run_mlr put '$n = length($y)'
echo "x=1,y=abcdefg,z=3" | run_mlr put '$n = length($nonesuch)'
echo "x=1,y=abcdefg,z=3" | run_mlr put '$n = length($*)'
echo "x=1,y=abcdefg,z=3" | run_mlr put '$n = length([])'
echo "x=1,y=abcdefg,z=3" | run_mlr put '$n = length([5,6,7])'
echo "x=1,y=abcdefg,z=3" | run_mlr put '$n = length({})'
echo "x=1,y=abcdefg,z=3" | run_mlr put '$n = length({"a":5,"b":6,"c":7})'

run_mlr --from $indir/s.dkvp put '$si = 0; for (i = 0; i < NR; i += 1) { if (i == 2) { $si += 0   } $si += i }'
run_mlr --from $indir/s.dkvp put '$si = 0; for (i = 0; i < NR; i += 1) { if (i == 2) { $si += 100 } $si += i }'
run_mlr --from $indir/s.dkvp put '$si = 0; for (i = 0; i < NR; i += 1) { if (i == 2) { break }      $si += i }'
run_mlr --from $indir/s.dkvp put '$si = 0; for (i = 0; i < NR; i += 1) { if (i == 2) { continue }   $si += i }'

run_mlr --from $indir/s.dkvp --opprint put '
  $si = 0;
  for (i = 0; i < NR; i += 1) {
    if (true) {
      if (i == 2) {
        $si += 0
      }
    }
    $si += i
  }'

run_mlr --from $indir/s.dkvp --opprint put '
  $si = 0;
  for (i = 0; i < NR; i += 1) {
    if (true) {
      if (i == 2) {
        $si += 100
      }
    }
    $si += i
  }'

run_mlr --from $indir/s.dkvp --opprint put '
  $si = 0;
  for (i = 0; i < NR; i += 1) {
    if (true) {
      if (i == 2) {
        break
      }
    }
    $si += i
  }'

run_mlr --from $indir/s.dkvp --opprint put '
  $si = 0;
  for (i = 0; i < NR; i += 1) {
    if (true) {
      if (i == 2) {
        continue
      }
    }
    $si += i
  }'

run_mlr --from $indir/s.dkvp --opprint put '
  $si = 0;
  for (p = 1; p <= 3; p += 1) {
    for (i = 0; i < NR; i += 1) {
      if (i == 2) {
        $si += 0
      }
      $si += i * 10**p
    }
  }'

run_mlr --from $indir/s.dkvp --opprint put '
  $si = 0;
  for (p = 1; p <= 3; p += 1) {
    for (i = 0; i < NR; i += 1) {
      if (i == 2) {
        break
      }
      $si += i * 10**p
    }
  }'

run_mlr --from $indir/s.dkvp --opprint put '
  $si = 0;
  for (p = 1; p <= 3; p += 1) {
    for (i = 0; i < NR; i += 1) {
      if (i == 2) {
        continue
      }
      $si += i * 10**p
    }
  }'

run_mlr --opprint --from $indir/ten.dkvp put -f $indir/f.mlr
run_mlr --opprint --from $indir/ten.dkvp put -f $indir/ff.mlr
run_mlr --opprint --from $indir/ten.dkvp put -f $indir/fg.mlr

run_mlr         --from $indir/s.dkvp put 'var x = 3'
run_mlr         --from $indir/s.dkvp put 'int x = 3'
run_mlr         --from $indir/s.dkvp put 'num x = 3'
mlr_expect_fail --from $indir/s.dkvp put 'str x = 3'
mlr_expect_fail --from $indir/s.dkvp put 'arr x = 3'

run_mlr         --from $indir/s.dkvp put 'func f(var x) { return 2*x} $y=f(3)'
run_mlr         --from $indir/s.dkvp put 'func f(int x) { return 2*x} $y=f(3)'
run_mlr         --from $indir/s.dkvp put 'func f(num x) { return 2*x} $y=f(3)'
mlr_expect_fail --from $indir/s.dkvp put 'func f(str x) { return 2*x} $y=f(3)'
mlr_expect_fail --from $indir/s.dkvp put 'func f(arr x) { return 2*x} $y=f(3)'

run_mlr         --from $indir/s.dkvp put 'func f(x): var { return 2*x} $y=f(3)'
run_mlr         --from $indir/s.dkvp put 'func f(x): int { return 2*x} $y=f(3)'
run_mlr         --from $indir/s.dkvp put 'func f(x): num { return 2*x} $y=f(3)'
mlr_expect_fail --from $indir/s.dkvp put 'func f(x): str { return 2*x} $y=f(3)'
mlr_expect_fail --from $indir/s.dkvp put 'func f(x): arr { return 2*x} $y=f(3)'

run_mlr --idkvp --opprint --from $indir/s.dkvp put '
  for (k, v in $*) {
    $["t".k] = typeof(v)
  }
  $tnonesuch = typeof($nonesuch)
'

run_mlr --idkvp --opprint --from $indir/s.dkvp put '
  for (k, v in $*) {
    $["s".k] = string(v)
  }
  $snonesuch = string($nonesuch)
'

run_mlr --j2p --from $indir/typecast.json put '
  $t       = typeof($a);
  $string  = string($a);
  $int     = int($a);
  $float   = float($a);
  $boolean = boolean($a);
' then reorder -f t,a

run_mlr -n put -f $indir/sieve.mlr
run_mlr --from /dev/null put -f $indir/sieve.mlr
run_mlr --from $indir/s.dkvp put -q -f $indir/sieve.mlr

mlr_expect_fail -n put 'begin{begin{}}'
mlr_expect_fail -n put 'begin{end{}}'
mlr_expect_fail -n put 'end{begin{}}'
mlr_expect_fail -n put 'end{end{}}'
mlr_expect_fail -n put 'begin { func f(x) { return 2*x} }'
# TODO: once subr exists
# mlr_expect_fail -n put 'begin { subr f(x) { return 2*x} }'
mlr_expect_fail -n put 'begin { emit $x }'
mlr_expect_fail -n put 'return 3'
mlr_expect_fail -n put 'break'
mlr_expect_fail -n put 'continue'
mlr_expect_fail -n put 'func f() { break }'
mlr_expect_fail -n put 'func f() { continue }'
run_mlr -n put -v 'true'
mlr_expect_fail -n put -v 'begin{true}'

run_mlr --from $indir/ten.dkvp head -n 1 then put -q '$v=[1,2,3,4,5]; unset $v[0]; dump $v'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q '$v=[1,2,3,4,5]; unset $v[1]; dump $v'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q '$v=[1,2,3,4,5]; unset $v[2]; dump $v'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q '$v=[1,2,3,4,5]; unset $v[3]; dump $v'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q '$v=[1,2,3,4,5]; unset $v[4]; dump $v'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q '$v=[1,2,3,4,5]; unset $v[5]; dump $v'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q '$v=[1,2,3,4,5]; unset $v[6]; dump $v'

run_mlr --from $indir/ten.dkvp head -n 1 then put -q 'end { @v=[1,2,3,4,5]; unset @v[0]; dump @v }'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q 'end { @v=[1,2,3,4,5]; unset @v[1]; dump @v }'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q 'end { @v=[1,2,3,4,5]; unset @v[2]; dump @v }'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q 'end { @v=[1,2,3,4,5]; unset @v[3]; dump @v }'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q 'end { @v=[1,2,3,4,5]; unset @v[4]; dump @v }'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q 'end { @v=[1,2,3,4,5]; unset @v[5]; dump @v }'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q 'end { @v=[1,2,3,4,5]; unset @v[6]; dump @v }'

run_mlr --from $indir/ten.dkvp head -n 1 then put -q 'end { v=[1,2,3,4,5]; unset v[0]; dump v }'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q 'end { v=[1,2,3,4,5]; unset v[1]; dump v }'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q 'end { v=[1,2,3,4,5]; unset v[2]; dump v }'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q 'end { v=[1,2,3,4,5]; unset v[3]; dump v }'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q 'end { v=[1,2,3,4,5]; unset v[4]; dump v }'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q 'end { v=[1,2,3,4,5]; unset v[5]; dump v }'
run_mlr --from $indir/ten.dkvp head -n 1 then put -q 'end { v=[1,2,3,4,5]; unset v[6]; dump v }'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '
  @v = {
    "a": {"x": 1},
    "b": {"y": 1},
  };
  if (NR == 1) {
    dump @v;
  } elif (NR == 2) {
    unset @v;
    dump @v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  @v = {
    "a": {"x": 1},
    "b": {"y": 1},
  };
  if (NR == 1) {
    dump @v;
  } elif (NR == 2) {
    unset @v["a"];
    dump @v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  @v = {
    "a": {"x": 1},
    "b": {"y": 1},
  };
  if (NR == 1) {
    dump @v;
  } elif (NR == 2) {
    unset @v["a"]["x"];
    dump @v;
  }
'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '
  @v = [
    {"x": 1},
    {"y": 1},
  ];
  if (NR == 1) {
    dump @v;
  } elif (NR == 2) {
    unset @v;
    dump @v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  @v = [
    {"x": 1},
    {"y": 1},
  ];
  if (NR == 1) {
    dump @v;
  } elif (NR == 2) {
    unset @v[1];
    dump @v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  @v = [
    {"x": 1},
    {"y": 1},
  ];
  if (NR == 1) {
    dump @v;
  } elif (NR == 2) {
    unset @v[1]["x"];
    dump @v;
  }
'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '
  @v = {
    "a": ["u", "v"],
    "b": ["w", "x"],
  };
  if (NR == 1) {
    dump @v;
  } elif (NR == 2) {
    unset @v;
    dump @v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  @v = {
    "a": ["u", "v"],
    "b": ["w", "x"],
  };
  if (NR == 1) {
    dump @v;
  } elif (NR == 2) {
    unset @v["a"];
    dump @v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  @v = {
    "a": ["u", "v"],
    "b": ["w", "x"],
  };
  if (NR == 1) {
    dump @v;
  } elif (NR == 2) {
    unset @v["a"][2];
    dump @v;
  }
'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '
  @v = [
    ["u", "v"],
    ["w", "x"],
  ];
  if (NR == 1) {
    dump @v;
  } elif (NR == 2) {
    unset @v;
    dump @v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  @v = [
    ["u", "v"],
    ["w", "x"],
  ];
  if (NR == 1) {
    dump @v;
  } elif (NR == 2) {
    unset @v[1];
    dump @v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  @v = [
    ["u", "v"],
    ["w", "x"],
  ];
  if (NR == 1) {
    dump @v;
  } elif (NR == 2) {
    unset @v[1][2];
    dump @v;
  }
'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '
  $v = {
    "a": {"x": 1},
    "b": {"y": 1},
  };
  if (NR == 1) {
    dump $v;
  } elif (NR == 2) {
    unset $v;
    dump $v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  $v = {
    "a": {"x": 1},
    "b": {"y": 1},
  };
  if (NR == 1) {
    dump $v;
  } elif (NR == 2) {
    unset $v["a"];
    dump $v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  $v = {
    "a": {"x": 1},
    "b": {"y": 1},
  };
  if (NR == 1) {
    dump $v;
  } elif (NR == 2) {
    unset $v["a"]["x"];
    dump $v;
  }
'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '
  $v = [
    {"x": 1},
    {"y": 1},
  ];
  if (NR == 1) {
    dump $v;
  } elif (NR == 2) {
    unset $v;
    dump $v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  $v = [
    {"x": 1},
    {"y": 1},
  ];
  if (NR == 1) {
    dump $v;
  } elif (NR == 2) {
    unset $v[1];
    dump $v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  $v = [
    {"x": 1},
    {"y": 1},
  ];
  if (NR == 1) {
    dump $v;
  } elif (NR == 2) {
    unset $v[1]["x"];
    dump $v;
  }
'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '
  $v = {
    "a": ["u", "v"],
    "b": ["w", "x"],
  };
  if (NR == 1) {
    dump $v;
  } elif (NR == 2) {
    unset $v;
    dump $v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  $v = {
    "a": ["u", "v"],
    "b": ["w", "x"],
  };
  if (NR == 1) {
    dump $v;
  } elif (NR == 2) {
    unset $v["a"];
    dump $v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  $v = {
    "a": ["u", "v"],
    "b": ["w", "x"],
  };
  if (NR == 1) {
    dump $v;
  } elif (NR == 2) {
    unset $v["a"][2];
    dump $v;
  }
'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '
  $v = [
    ["u", "v"],
    ["w", "x"],
  ];
  if (NR == 1) {
    dump $v;
  } elif (NR == 2) {
    unset $v;
    dump $v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  $v = [
    ["u", "v"],
    ["w", "x"],
  ];
  if (NR == 1) {
    dump $v;
  } elif (NR == 2) {
    unset $v[1];
    dump $v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  $v = [
    ["u", "v"],
    ["w", "x"],
  ];
  if (NR == 1) {
    dump $v;
  } elif (NR == 2) {
    unset $v[1][2];
    dump $v;
  }
'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '
  $* = {
    "a": {"x": 1},
    "b": {"y": 1},
  };
  if (NR == 1) {
    dump $*;
  } elif (NR == 2) {
    unset $*;
    dump $*;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  $* = {
    "a": {"x": 1},
    "b": {"y": 1},
  };
  if (NR == 1) {
    dump $*;
  } elif (NR == 2) {
    unset $*["a"];
    dump $*;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  $* = {
    "a": {"x": 1},
    "b": {"y": 1},
  };
  if (NR == 1) {
    dump $*;
  } elif (NR == 2) {
    unset $*["a"]["x"];
    dump $*;
  }
'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '
  $* = {
    "a": ["u", "v"],
    "b": ["w", "x"],
  };
  if (NR == 1) {
    dump $*;
  } elif (NR == 2) {
    unset $*;
    dump $*;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  $* = {
    "a": ["u", "v"],
    "b": ["w", "x"],
  };
  if (NR == 1) {
    dump $*;
  } elif (NR == 2) {
    unset $*["a"];
    dump $*;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  $* = {
    "a": ["u", "v"],
    "b": ["w", "x"],
  };
  if (NR == 1) {
    dump $*;
  } elif (NR == 2) {
    unset $*["a"][2];
    dump $*;
  }
'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '
  v = {
    "a": {"x": 1},
    "b": {"y": 1},
  };
  if (NR == 1) {
    dump v;
  } elif (NR == 2) {
    unset v;
    dump v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  v = {
    "a": {"x": 1},
    "b": {"y": 1},
  };
  if (NR == 1) {
    dump v;
  } elif (NR == 2) {
    unset v["a"];
    dump v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  v = {
    "a": {"x": 1},
    "b": {"y": 1},
  };
  if (NR == 1) {
    dump v;
  } elif (NR == 2) {
    unset v["a"]["x"];
    dump v;
  }
'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '
  v = [
    {"x": 1},
    {"y": 1},
  ];
  if (NR == 1) {
    dump v;
  } elif (NR == 2) {
    unset v;
    dump v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  v = [
    {"x": 1},
    {"y": 1},
  ];
  if (NR == 1) {
    dump v;
  } elif (NR == 2) {
    unset v[1];
    dump v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  v = [
    {"x": 1},
    {"y": 1},
  ];
  if (NR == 1) {
    dump v;
  } elif (NR == 2) {
    unset v[1]["x"];
    dump v;
  }
'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '
  v = {
    "a": ["u", "v"],
    "b": ["w", "x"],
  };
  if (NR == 1) {
    dump v;
  } elif (NR == 2) {
    unset v;
    dump v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  v = {
    "a": ["u", "v"],
    "b": ["w", "x"],
  };
  if (NR == 1) {
    dump v;
  } elif (NR == 2) {
    unset v["a"];
    dump v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  v = {
    "a": ["u", "v"],
    "b": ["w", "x"],
  };
  if (NR == 1) {
    dump v;
  } elif (NR == 2) {
    unset v["a"][2];
    dump v;
  }
'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp put -q '
  v = [
    ["u", "v"],
    ["w", "x"],
  ];
  if (NR == 1) {
    dump v;
  } elif (NR == 2) {
    unset v;
    dump v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  v = [
    ["u", "v"],
    ["w", "x"],
  ];
  if (NR == 1) {
    dump v;
  } elif (NR == 2) {
    unset v[1];
    dump v;
  }
'

run_mlr --from $indir/s.dkvp put -q '
  v = [
    ["u", "v"],
    ["w", "x"],
  ];
  if (NR == 1) {
    dump v;
  } elif (NR == 2) {
    unset v[1][2];
    dump v;
  }
'

# ----------------------------------------------------------------
run_mlr --from $indir/s.dkvp --opprint put '
  $u = toupper($a);
  $l = tolower($u);
  $c = capitalize($l);
'

run_mlr --from $indir/abixy-het put '$z = $a =~ ".[ak]."'
run_mlr --from $indir/abixy-het put '$z = $a !=~ ".[ak]."'

run_mlr --opprint --from $indir/s.dkvp put '
  $dx = depth($x);
  $dn = depth($nonesuch);
  $da1 = depth([1,2,3]);
  $da2 = depth([1,[4,5,6],3]);
  $da3 = depth([1,{"s":4,"t":[7,8,9],"u":6},3]);
  $dm1 = depth({"s":1,"t":2,"u":3});
  $dm2 = depth({"s":1,"t":[4,5,6],"u":3});
  $dm3 = depth({"s":1,"t":[4,$*,6],"u":3});
'

run_mlr --opprint --from $indir/s.dkvp put '
  $lcx = leafcount($x);
  $lcn = leafcount($nonesuch);
  $lca1 = leafcount([1,2,3]);
  $lca2 = leafcount([1,[4,5,6],3]);
  $lca3 = leafcount([1,{"s":4,"t":[7,8,9],"u":6},3]);
  $lcm1 = leafcount({"s":1,"t":2,"u":3});
  $lcm2 = leafcount({"s":1,"t":[4,5,6],"u":3});
  $lcm3 = leafcount({"s":1,"t":[4,{"x":8, "y": 9},6],"u":3});
'

run_mlr --oxtab --from $indir/s.dkvp head -n 1 then put '
  $hk01 = haskey($x, $a);
  $hk02 = haskey($nonesuch, $a);
  $hk03 = haskey($*, 7);
  $hk04 = haskey($*, "a");
  $hk05 = haskey($*, "nonesuch");
  $hk06 = haskey([10,20,30], 0);
  $hk07 = haskey([10,20,30], 1);
  $hk08 = haskey([10,20,30], 2);
  $hk09 = haskey([10,20,30], 3);
  $hk10 = haskey([10,20,30], 4);
  $hk11 = haskey([10,20,30], -4);
  $hk12 = haskey([10,20,30], -3);
  $hk13 = haskey([10,20,30], -2);
  $hk14 = haskey([10,20,30], -1);
  $hk15 = haskey([10,20,30], "a");
'

# ----------------------------------------------------------------
run_mlr --from $indir/abixy-het put -q 'o = mapsum(); dump o'
run_mlr --from $indir/abixy-het put -q 'o = mapdiff(); dump o'
mlr_expect_fail --from $indir/abixy-het put -q 'o = mapexcept(); dump o'
mlr_expect_fail --from $indir/abixy-het put -q 'o = mapselect(); dump o'

run_mlr --from $indir/abixy-het put -q 'o = mapsum($*); dump o'
run_mlr --from $indir/abixy-het put -q 'o = mapsum({"a":999}); dump o'
run_mlr --from $indir/abixy-het put -q 'o = mapdiff($*); dump o'
run_mlr --from $indir/abixy-het put -q 'o = mapdiff({"a":999}); dump o'
run_mlr --from $indir/abixy-het put -q 'o = mapexcept($*); dump o'
run_mlr --from $indir/abixy-het put -q 'o = mapexcept({"a":999}); dump o'
run_mlr --from $indir/abixy-het put -q 'o = mapselect($*); dump o'
run_mlr --from $indir/abixy-het put -q 'o = mapselect({"a":999}); dump o'

run_mlr --from $indir/abixy-het put -q 'o = mapsum($*, {"a":999}); dump o'
run_mlr --from $indir/abixy-het put -q 'o = mapsum({"a":999}, $*); dump o'
run_mlr --from $indir/abixy-het put -q 'o = mapdiff($*, {"a":999}); dump o'
run_mlr --from $indir/abixy-het put -q 'o = mapdiff({"a":999}, $*); dump o'
run_mlr --from $indir/abixy-het put -q 'o = mapexcept($*, "a"); dump o'
run_mlr --from $indir/abixy-het put -q 'o = mapexcept({"a":999}, "a"); dump o'
run_mlr --from $indir/abixy-het put -q 'o = mapexcept({"a":999}, "nonesuch"); dump o'
run_mlr --from $indir/abixy-het put -q 'o = mapselect($*, "a"); dump o'
run_mlr --from $indir/abixy-het put -q 'o = mapselect({"a":999}, "a"); dump o'
run_mlr --from $indir/abixy-het put -q 'o = mapselect({"a":999}, "nonesuch"); dump o'

run_mlr --from $indir/abixy-het put -q 'o = mapsum($*, {"a":999}, {"b": 444}); dump o'
run_mlr --from $indir/abixy-het put -q 'o = mapsum({"a":999}, $*, {"b": 444}); dump o'
run_mlr --from $indir/abixy-het put -q 'o = mapdiff($*, {"a":999}, {"b": 444}); dump o'
run_mlr --from $indir/abixy-het put -q 'o = mapdiff({"a":999}, $*, {"b": 444}); dump o'
run_mlr --from $indir/abixy-het put -q 'o = mapexcept($*, "a", "b"); dump o'
run_mlr --from $indir/abixy-het put -q 'o = mapselect({"a":999}, "b", "nonesuch"); dump o'
run_mlr --from $indir/abixy-het put -q 'o = mapselect($*, "a", "b"); dump o'
run_mlr --from $indir/abixy-het put -q 'o = mapselect({"a":999}, "b", "nonesuch"); dump o'

run_mlr --from $indir/abixy-het put -q 'o = mapselect($*, ["b", "nonesuch"]); dump o'
run_mlr --from $indir/abixy-het put -q 'o = mapexcept($*, ["b", "nonesuch"]); dump o'

run_mlr --c2p --from $indir/mod.csv put '
  $add = madd($a, $b, $m);
  $sub = msub($a, $b, $m);
  $mul = mmul($a, $b, $m);
  $exp = mexp($a, $b, $m);
'

run_mlr --opprint --from $indir/ten.dkvp put '
  $ha = hexfmt($a);
  $hx = hexfmt($x);
  $hi = hexfmt($i);
  $nhi = hexfmt(-$i);
'

run_mlr --opprint --from $indir/ten.dkvp put '
  $hi = "0x".fmtnum($i, "%04x");
  $ex = fmtnum($x, "%8.3e");
'
