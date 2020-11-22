# Note: filter -v and put -v print the AST.

run_mlr put    -v '$x = 1 || 2 ^^ 3'   /dev/null
run_mlr filter -v '     1 || 2 ^^ 3'   /dev/null
run_mlr put    -v '$x = 1 || 2 && 3'   /dev/null
run_mlr filter -v '     1 || 2 && 3'   /dev/null

run_mlr put    -v '$x = 1 ^^ 2 || 3'   /dev/null
run_mlr filter -v '     1 ^^ 2 || 3'   /dev/null
run_mlr put    -v '$x = 1 ^^ 2 && 3'   /dev/null
run_mlr filter -v '     1 ^^ 2 && 3'   /dev/null

run_mlr put    -v '$x = 1 && 2 || 3'   /dev/null
run_mlr filter -v '     1 && 2 || 3'   /dev/null
run_mlr put    -v '$x = 1 && 2 ^^ 3'   /dev/null
run_mlr filter -v '     1 && 2 ^^ 3'   /dev/null

run_mlr put    -v '$x =  1 == 2 <= 3'  /dev/null
run_mlr filter -v '      1 == 2 <= 3'  /dev/null
run_mlr put    -v '$x =  1 <= 2 == 3'  /dev/null
run_mlr filter -v '      1 <= 2 == 3'  /dev/null

run_mlr put    -v '$x =  1 <= 2 |  3'  /dev/null
run_mlr filter -v '      1 <= 2 |  3'  /dev/null
run_mlr put    -v '$x =  1 |  2 <= 3'  /dev/null
run_mlr filter -v '      1 |  2 <= 3'  /dev/null

run_mlr put    -v '$x =  1 |  2 ^  3'  /dev/null
run_mlr filter -v '      1 |  2 ^  3'  /dev/null
run_mlr put    -v '$x =  1 ^  2 |  3'  /dev/null
run_mlr filter -v '      1 ^  2 |  3'  /dev/null

run_mlr put    -v '$x =  1 ^  2 &  3'  /dev/null
run_mlr filter -v '      1 ^  2 &  3'  /dev/null
run_mlr put    -v '$x =  1 &  2 ^  3'  /dev/null
run_mlr filter -v '      1 &  2 ^  3'  /dev/null

run_mlr put    -v '$x =  1 &  2 << 3'  /dev/null
run_mlr filter -v '      1 &  2 << 3'  /dev/null
run_mlr put    -v '$x =  1 << 2 &  3'  /dev/null
run_mlr filter -v '      1 << 2 &  3'  /dev/null

run_mlr put    -v '$x =  1 +  2 * 3'   /dev/null
run_mlr filter -v '      1 +  2 * 3'   /dev/null
run_mlr put    -v '$x =  1 *  2 + 3'   /dev/null
run_mlr filter -v '      1 *  2 + 3'   /dev/null
run_mlr put    -v '$x =  1 + (2 * 3)'  /dev/null
run_mlr filter -v '      1 + (2 * 3)'  /dev/null
run_mlr put    -v '$x =  1 * (2 + 3)'  /dev/null
run_mlr filter -v '      1 * (2 + 3)'  /dev/null
run_mlr put    -v '$x = (1 + 2) * 3'   /dev/null
run_mlr filter -v '     (1 + 2) * 3'   /dev/null
run_mlr put    -v '$x = (1 * 2) + 3'   /dev/null
run_mlr filter -v '     (1 * 2) + 3'   /dev/null

run_mlr put    -v '$x =  1 +   2 ** 3'  /dev/null
run_mlr filter -v '      1 +   2 ** 3'  /dev/null
run_mlr put    -v '$x =  1 **  2 +  3'  /dev/null
run_mlr filter -v '      1 **  2 +  3'  /dev/null
run_mlr put    -v '$x =  1 +  (2 ** 3)' /dev/null
run_mlr filter -v '      1 +  (2 ** 3)' /dev/null
run_mlr put    -v '$x =  1 ** (2 +  3)' /dev/null
run_mlr filter -v '      1 ** (2 +  3)' /dev/null
run_mlr put    -v '$x = (1 +  2) ** 3'  /dev/null
run_mlr filter -v '     (1 +  2) ** 3'  /dev/null
run_mlr put    -v '$x = (1 ** 2) +  3'  /dev/null
run_mlr filter -v '     (1 ** 2) +  3'  /dev/null

run_mlr put    -v '$x =  1 *   2 ** 3'  /dev/null
run_mlr filter -v '      1 *   2 ** 3'  /dev/null
run_mlr put    -v '$x =  1 **  2 *  3'  /dev/null
run_mlr filter -v '      1 **  2 *  3'  /dev/null
run_mlr put    -v '$x =  1 *  (2 ** 3)' /dev/null
run_mlr filter -v '      1 *  (2 ** 3)' /dev/null
run_mlr put    -v '$x =  1 ** (2 *  3)' /dev/null
run_mlr filter -v '      1 ** (2 *  3)' /dev/null
run_mlr put    -v '$x = (1 *  2) ** 3'  /dev/null
run_mlr filter -v '     (1 *  2) ** 3'  /dev/null
run_mlr put    -v '$x = (1 ** 2) *  3'  /dev/null
run_mlr filter -v '     (1 ** 2) *  3'  /dev/null

run_mlr put    -v '$x = -1 +  2 *  3'  /dev/null
run_mlr filter -v '     -1 +  2 *  3'  /dev/null
run_mlr put    -v '$x = -1 *  2 +  3'  /dev/null
run_mlr filter -v '     -1 *  2 +  3'  /dev/null
run_mlr put    -v '$x =  1 + -2 *  3'  /dev/null
run_mlr filter -v '      1 + -2 *  3'  /dev/null
run_mlr put    -v '$x =  1 * -2 +  3'  /dev/null
run_mlr filter -v '      1 * -2 +  3'  /dev/null
run_mlr put    -v '$x =  1 +  2 * -3'  /dev/null
run_mlr filter -v '      1 +  2 * -3'  /dev/null
run_mlr put    -v '$x =  1 *  2 + -3'  /dev/null
run_mlr filter -v '      1 *  2 + -3'  /dev/null

run_mlr put    -v '$x = ~1 |  2 &  3'  /dev/null
run_mlr filter -v '     ~1 |  2 &  3'  /dev/null
run_mlr put    -v '$x = ~1 &  2 |  3'  /dev/null
run_mlr filter -v '     ~1 &  2 |  3'  /dev/null
run_mlr put    -v '$x =  1 | ~2 &  3'  /dev/null
run_mlr filter -v '      1 | ~2 &  3'  /dev/null
run_mlr put    -v '$x =  1 & ~2 |  3'  /dev/null
run_mlr filter -v '      1 & ~2 |  3'  /dev/null
run_mlr put    -v '$x =  1 |  2 & ~3'  /dev/null
run_mlr filter -v '      1 |  2 & ~3'  /dev/null
run_mlr put    -v '$x =  1 &  2 | ~3'  /dev/null
run_mlr filter -v '      1 &  2 | ~3'  /dev/null

run_mlr put    -v '$x = $a==1 && $b == 1 && $c == 1' /dev/null
run_mlr filter -v '     $a==1 && $b == 1 && $c == 1' /dev/null
run_mlr put    -v '$x = $a==1 || $b == 1 && $c == 1' /dev/null
run_mlr filter -v '     $a==1 || $b == 1 && $c == 1' /dev/null
run_mlr put    -v '$x = $a==1 || $b == 1 || $c == 1' /dev/null
run_mlr filter -v '     $a==1 || $b == 1 || $c == 1' /dev/null
run_mlr put    -v '$x = $a==1 && $b == 1 || $c == 1' /dev/null
run_mlr filter -v '     $a==1 && $b == 1 || $c == 1' /dev/null

run_mlr put    -v '$x = $a==1 ? $b == 2 : $c == 3' /dev/null
run_mlr filter -v '     $a==1 ? $b == 2 : $c == 3' /dev/null

run_mlr put    -v '$x = true' /dev/null
run_mlr filter -v '     true' /dev/null

run_mlr put    -v 'true || 1==0; $x = 3' /dev/null
run_mlr filter -v '        true || 1==0' /dev/null

run_mlr put    -v '1==0 || false; $x = 3' /dev/null
run_mlr filter -v '        1==0 || false' /dev/null

run_mlr put    -v 'true && false; $x = 3' /dev/null
run_mlr filter -v '        true && false' /dev/null

run_mlr put    -v 'true && false && true; $x = 3' /dev/null
run_mlr filter -v '        true && false && true' /dev/null

run_mlr put    -v '$y += $x + 3'  /dev/null
run_mlr put    -v '$y += $x * 3'  /dev/null

run_mlr put -v '$y ||= $x' /dev/null
run_mlr put -v '$y ^^= $x' /dev/null
run_mlr put -v '$y &&= $x' /dev/null
run_mlr put -v '$y |=  $x' /dev/null
run_mlr put -v '$y ^=  $x' /dev/null
run_mlr put -v '$y &=  $x' /dev/null
run_mlr put -v '$y <<= $x' /dev/null
run_mlr put -v '$y >>= $x' /dev/null
run_mlr put -v '$y +=  $x' /dev/null
run_mlr put -v '$y -=  $x' /dev/null
run_mlr put -v '$y .=  $x' /dev/null
run_mlr put -v '$y *=  $x' /dev/null
run_mlr put -v '$y /=  $x' /dev/null
run_mlr put -v '$y //= $x' /dev/null
run_mlr put -v '$y %=  $x' /dev/null
run_mlr put -v '$y **= $x' /dev/null
