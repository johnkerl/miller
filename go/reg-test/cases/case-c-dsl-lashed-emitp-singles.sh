# TODO: bracketed-oosvar emits aren't in the mlrgo BNF
run_mlr -n put '
end {
  @a = 111;
  emitp @a
}'
run_mlr -n put '
end {
  @a = 111;
  emitp (@a)
}'

run_mlr -n put '
end {
  @a[111] = 222;
  emitp @a, "s"
}'
run_mlr -n put '
end {
  @a[111] = 222;
  emitp (@a), "s"
}'

run_mlr -n put '
end {
  @a[111] = 222;
  @a[333] = 444;
  emitp @a, "s"
}'

run_mlr -n put '
end {
  @a[111] = 222;
  @a[333] = 444;
  emitp (@a), "s"
}'

run_mlr -n put '
end {
  @a[111][222] = 333;
  emitp @a, "s"
}'
run_mlr -n put '
end {
  @a[111][222] = 333;
  emitp (@a), "s"
}'

run_mlr -n put '
end {
  @a[111][222] = 333;
  @a[444][555] = 666;
  emitp @a, "s"
}'
run_mlr -n put '
end {
  @a[111][222] = 333;
  @a[444][555] = 666;
  emitp (@a), "s"
}'

run_mlr -n put '
end {
  @a[111][222] = 333;
  emitp @a, "s", "t"
}'
run_mlr -n put '
end {
  @a[111][222] = 333;
  emitp (@a), "s", "t"
}'

run_mlr -n put '
end {
  @a[111][222] = 333;
  emitp @a[111], "t"
}'

run_mlr -n put '
end {
  @a[111][222] = 333;
  emitp (@a[111]), "t"
}'

# ----------------------------------------------------------------
announce LASHED EMIT SINGLES

run_mlr -n put '
end {
  @a = 111;
  emit @a
}'
run_mlr -n put '
end {
  @a = 111;
  emit (@a)
}'

run_mlr -n put '
end {
  @a[111] = 222;
  emit @a, "s"
}'
run_mlr -n put '
end {
  @a[111] = 222;
  emit (@a), "s"
}'

run_mlr -n put '
end {
  @a[111] = 222;
  @a[333] = 444;
  emit @a, "s"
}'
run_mlr -n put '
end {
  @a[111] = 222;
  @a[333] = 444;
  emit (@a), "s"
}'

run_mlr -n put '
end {
  @a[111][222] = 333;
  emit @a, "s"
}'
run_mlr -n put '
end {
  @a[111][222] = 333;
  emit (@a), "s"
}'

# ----------------------------------------------------------------
announce LASHED EMITP PAIRS

run_mlr -n put '
end {
  @a = 111;
  @b = 222;
  emitp (@a, @b)
}'

run_mlr -n put '
end {
  @a[1] = 111;
  @b[1] = 222;
  emitp (@a[1], @b[1])
}'

run_mlr -n put '
end {
  @a[1][2][3] = 4;
  @b[1][2][3] = 8;
  emitp (@a, @b), "s", "t", "u"
}'

run_mlr -n put '
end {
  @a[1][2][3] = 4;
  @b[5][2][3] = 8;
  emitp (@a[1], @b[5]), "t", "u"
}'

run_mlr -n put '
end {
  @a[1][2][3] = 4;
  @b[5][6][3] = 8;
  emitp (@a[1][2], @b[5][6]), "u"
}'

# ----------------------------------------------------------------
announce LASHED EMIT PAIRS

run_mlr -n put '
end {
  @a = 111;
  @b = 222;
  emit (@a, @b)
}'

run_mlr -n put '
end {
  @a[1] = 111;
  @b[1] = 222;
  emit (@a[1], @b[1])
}'

run_mlr -n put '
end {
  @a[1][2][3] = 4;
  @b[1][2][3] = 8;
  emit (@a, @b), "s", "t", "u"
}'

run_mlr -n put '
end {
  @a[1][2][3] = 4;
  @b[5][2][3] = 8;
  emit (@a[1], @b[5]), "t", "u"
}'

run_mlr -n put '
end {
  @a[1][2][3] = 4;
  @b[5][6][3] = 8;
  emit (@a[1][2], @b[5][6]), "u"
}'

# ----------------------------------------------------------------
announce LASHED EMIT WITH VARYING DEPTH

run_mlr -n put '
end {
  @a[1][1] = 1;
  @a[1][2] = 2;
  @a[2][1] = 3;
  @a[2][2] = 4;
  @b[1][1] = 5;
  @b[1][2] = 6;
  @b[2][1] = 7;
  @b[2][2] = 8;
  emit (@a[1], @b[2]), "t"
}'

run_mlr -n put '
end {
  @a[1][1] = 1;
  @a[1][2] = 2;
  @a[2][1] = 3;
  @a[2][2] = 4;
  @b[1][1] = 5;
  @b[1][2] = 6;
  @b[2][1] = 7;
  @b[2][2] = 8;
  emit (@a, @b), "s", "t"
}'

run_mlr -n put '
end {
  @a[1][1] = 1;
  @a[1][2] = 2;
  @a[2][1] = 3;
  @a[2][2] = 4;
  @a[3] = 10;
  @a[4] = 11;
  @a[5][6][7] = 12;
  @b[1][1] = 5;
  @b[1][2] = 6;
  @b[2][1] = 7;
  @b[2][2] = 8;
  emit (@a, @b), "s", "t"
}'

run_mlr -n put '
end {
  @a[1][1] = 1;
  @a[1][2] = 2;
  @a[2][1] = 3;
  @a[2][2] = 4;
  @a[3] = 10;
  @a[4] = 11;
  @a[5][6][7] = 12;
  @b[1][1] = 5;
  @b[1][2] = 6;
  @b[2][1] = 7;
  @b[2][2] = 8;
  emit (@b, @a), "s", "t"
}'


run_mlr -n put '
end {
  @a[1][2][3] = 4;
  @b[5][2][3] = 8;
  emit (@a[1], @b[3]), "t", "u"
}'

run_mlr -n put '
end {
  @a[1][2][3] = 4;
  @b[5][2][3] = 8;
  emit (@a[1][2], @b[5][9]), "t", "u"
}'

run_mlr -n put '
end {
  @a[1][2][3] = 4;
  @b[5][2][3] = 8;
  emit (@a[1][2], @b[9][2]), "t", "u"
}'

run_mlr -n put '
end {
  @a[1][2][3] = 4;
  @b[5][2][3] = 8;
  emit (@a[9], @b[5]), "t", "u"
}'

# ----------------------------------------------------------------
announce LASHED EMITP WITH VARYING DEPTH

run_mlr -n put '
end {
  @a[1][1] = 1;
  @a[1][2] = 2;
  @a[2][1] = 3;
  @a[2][2] = 4;
  @b[1][1] = 5;
  @b[1][2] = 6;
  @b[2][1] = 7;
  @b[2][2] = 8;
  emitp (@a[1], @b[2]), "t"
}'

run_mlr -n put '
end {
  @a[1][1] = 1;
  @a[1][2] = 2;
  @a[2][1] = 3;
  @a[2][2] = 4;
  @b[1][1] = 5;
  @b[1][2] = 6;
  @b[2][1] = 7;
  @b[2][2] = 8;
  emitp (@a, @b), "s", "t"
}'

run_mlr -n put '
end {
  @a[1][1] = 1;
  @a[1][2] = 2;
  @a[2][1] = 3;
  @a[2][2] = 4;
  @a[3] = 10;
  @a[4] = 11;
  @a[5][6][7] = 12;
  @b[1][1] = 5;
  @b[1][2] = 6;
  @b[2][1] = 7;
  @b[2][2] = 8;
  emitp (@a, @b), "s", "t"
}'

run_mlr -n put '
end {
  @a[1][1] = 1;
  @a[1][2] = 2;
  @a[2][1] = 3;
  @a[2][2] = 4;
  @a[3] = 10;
  @a[4] = 11;
  @a[5][6][7] = 12;
  @b[1][1] = 5;
  @b[1][2] = 6;
  @b[2][1] = 7;
  @b[2][2] = 8;
  emitp (@b, @a), "s", "t"
}'


run_mlr -n put '
end {
  @a[1][2][3] = 4;
  @b[5][2][3] = 8;
  emitp (@a[1], @b[3]), "t", "u"
}'

run_mlr -n put '
end {
  @a[1][2][3] = 4;
  @b[5][2][3] = 8;
  emitp (@a[1][2], @b[5][9]), "t", "u"
}'

run_mlr -n put '
end {
  @a[1][2][3] = 4;
  @b[5][2][3] = 8;
  emitp (@a[1][2], @b[9][2]), "t", "u"
}'

run_mlr -n put '
end {
  @a[1][2][3] = 4;
  @b[5][2][3] = 8;
  emitp (@a[9], @b[5]), "t", "u"
}'

# ----------------------------------------------------------------
announce CANONICAL LASHED EMIT

run_mlr --from $indir/abixy-wide --opprint put -q '
  @count[$a] += 1;
  @sum[$a] += $x;
  end {
      for (a, c in @count) {
          @mean[a] = @sum[a] / @count[a]
      }
      emit (@sum, @count, @mean), "a"
  }
'

run_mlr --from $indir/abixy-wide --opprint put -q '
  @count[$a][$b] += 1;
  @sum[$a][$b] += $x;
  end {
      for ((a, b), c in @count) {
          @mean[a][b] = @sum[a][b] / @count[a][b]
      }
      emit (@sum, @count, @mean), "a", "b"
  }
'

# ----------------------------------------------------------------
announce MAP-VARIANT EMITS

run_mlr         --from $indir/abixy-het --opprint put -q 'o=$a.$b; emit o'
run_mlr         --from $indir/abixy-het --opprint put -q 'o={"ab":$a.$b}; emit o'

run_mlr         --from $indir/abixy-het --opprint put -q '@o=$a.$b; emit @o'
run_mlr         --from $indir/abixy-het --opprint put -q '@o={"ab":$a.$b}; emit @o'

run_mlr         --from $indir/abixy-het --opprint put -q '@o=$a.$b; emit @*'
run_mlr         --from $indir/abixy-het --opprint put -q '@o={"ab":$a.$b}; emit @*'

mlr_expect_fail --from $indir/abixy-het --opprint put -q 'emit $a.$b'
run_mlr         --from $indir/abixy-het --opprint put -q 'emit {"ab":$a.$b}'

run_mlr         --from $indir/abixy-het --opprint put -q 'func f(a,b) { return a.b } o = f($a, $b); emit o'
run_mlr         --from $indir/abixy-het --opprint put -q 'func f(a,b) { return a.b } emit f($a, $b)'

run_mlr         --from $indir/abixy-het --opprint put -q 'func f(a,b) { return {"ab": a.b} } o = f($a, $b); emit o'
run_mlr         --from $indir/abixy-het --opprint put -q 'func f(a,b) { return {"ab": a.b} } emit f($a, $b)'

# ----------------------------------------------------------------
announce MAP-VARIANT LASHED EMITS

# scalar emits aren't reasonable ... but they shouldn't segv

mention scalar lashed emits
mlr_expect_fail --from $indir/abixy-het --opprint put -q 'emit ($a . "_" . $b, $x . "_" . $y)'
run_mlr         --from $indir/abixy-het --opprint put -q ' o = $a . "_" . $b;  p = $x . "_" . $y; emit  (o,  p)'
run_mlr         --from $indir/abixy-het --opprint put -q '@o = $a . "_" . $b; @p = $x . "_" . $y; emit (@o, @p)'
run_mlr         --from $indir/abixy-het --opprint put -q 'func f(a, b) { return a . "_" . b }  o = f($a, $b);  p = f($x, $y); emit  (o,  p)'
run_mlr         --from $indir/abixy-het --opprint put -q 'func f(a, b) { return a . "_" . b } @o = f($a, $b); @p = f($x, $y); emit (@o, @p)'
run_mlr         --from $indir/abixy-het --opprint put -q 'func f(a, b) { return a . "_" . b } emit (f($a, $b), f($x, $y))'

mention non-scalar non-keyed lashed emits
run_mlr         --from $indir/abixy-het --opprint put -q 'emit ({"ab": $a . "_" . $b}, {"ab": $x . "_" . $y})'
run_mlr         --from $indir/abixy-het --opprint put -q ' o = {"ab": $a . "_" . $b};  p = {"ab": $x . "_" . $y}; emit  (o, p)'
run_mlr         --from $indir/abixy-het --opprint put -q '@o = {"ab": $a . "_" . $b}; @p = {"ab": $x . "_" . $y}; emit (@o, @p)'
run_mlr         --from $indir/abixy-het --opprint put -q 'func f(a, b) { return {"ab": a . "_" . b} }  o = f($a, $b);  p = f($x, $y); emit  (o, p)'
run_mlr         --from $indir/abixy-het --opprint put -q 'func f(a, b) { return {"ab": a . "_" . b} } @o = f($a, $b); @p = f($x, $y); emit (@o, @p)'
run_mlr         --from $indir/abixy-het --opprint put -q 'func f(a, b) { return {"ab": a . "_" . b} } emit (f($a, $b), f($x, $y))'

mention non-scalar non-keyed lashed emits
run_mlr         --from $indir/abixy-het --opprint put -q 'emitp ({"ab": $a . "_" . $b}, {"ab": $x . "_" . $y})'
run_mlr         --from $indir/abixy-het --opprint put -q ' o = {"ab": $a . "_" . $b};  p = {"ab": $x . "_" . $y}; emitp  (o, p)'
run_mlr         --from $indir/abixy-het --opprint put -q '@o = {"ab": $a . "_" . $b}; @p = {"ab": $x . "_" . $y}; emitp (@o, @p)'
run_mlr         --from $indir/abixy-het --opprint put -q 'func f(a, b) { return {"ab": a . "_" . b} }  o = f($a, $b);  p = f($x, $y); emitp  (o, p)'
run_mlr         --from $indir/abixy-het --opprint put -q 'func f(a, b) { return {"ab": a . "_" . b} } @o = f($a, $b); @p = f($x, $y); emitp (@o, @p)'
run_mlr         --from $indir/abixy-het --opprint put -q 'func f(a, b) { return {"ab": a . "_" . b} } emit (f($a, $b), f($x, $y))'

mention non-scalar keyed lashed emits
run_mlr         --from $indir/abixy-het --opprint put -q 'emit ({"ab": $a . "_" . $b}, {"ab": $x . "_" . $y}), "ab"'
run_mlr         --from $indir/abixy-het --opprint put -q ' o = {"ab": $a . "_" . $b};  p = {"ab": $x . "_" . $y}; emit  (o, p), "ab"'
run_mlr         --from $indir/abixy-het --opprint put -q '@o = {"ab": $a . "_" . $b}; @p = {"ab": $x . "_" . $y}; emit (@o, @p), "ab"'
run_mlr         --from $indir/abixy-het --opprint put -q 'func f(a, b) { return {"ab": a . "_" . b} }  o = f($a, $b);  p = f($x, $y); emit  (o, p), "ab"'
run_mlr         --from $indir/abixy-het --opprint put -q 'func f(a, b) { return {"ab": a . "_" . b} } @o = f($a, $b); @p = f($x, $y); emit (@o, @p), "ab"'
run_mlr         --from $indir/abixy-het --opprint put -q 'func f(a, b) { return {"ab": a . "_" . b} } emit (f($a, $b), f($x, $y)), "ab"'

mention non-scalar keyed lashed emits
run_mlr         --from $indir/abixy-het --opprint put -q 'emitp ({"ab": $a . "_" . $b}, {"ab": $x . "_" . $y}), "ab"'
run_mlr         --from $indir/abixy-het --opprint put -q ' o = {"ab": $a . "_" . $b};  p = {"ab": $x . "_" . $y}; emitp  (o, p), "ab"'
run_mlr         --from $indir/abixy-het --opprint put -q '@o = {"ab": $a . "_" . $b}; @p = {"ab": $x . "_" . $y}; emitp (@o, @p), "ab"'
run_mlr         --from $indir/abixy-het --opprint put -q 'func f(a, b) { return {"ab": a . "_" . b} }  o = f($a, $b);  p = f($x, $y); emitp  (o, p), "ab"'
run_mlr         --from $indir/abixy-het --opprint put -q 'func f(a, b) { return {"ab": a . "_" . b} } @o = f($a, $b); @p = f($x, $y); emitp (@o, @p), "ab"'
run_mlr         --from $indir/abixy-het --opprint put -q 'func f(a, b) { return {"ab": a . "_" . b} } emitp (f($a, $b), f($x, $y)), "ab"'
