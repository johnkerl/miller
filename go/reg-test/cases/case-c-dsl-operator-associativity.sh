# Note: filter -v and put -v print the AST.

run_mlr put    -v '$x = 1 || 2 || 3'   /dev/null
run_mlr filter -v '     1 || 2 || 3'   /dev/null
run_mlr put    -v '$x = 1 ^^ 2 ^^ 3'   /dev/null
run_mlr filter -v '     1 ^^ 2 ^^ 3'   /dev/null
run_mlr put    -v '$x = 1 && 2 && 3'   /dev/null
run_mlr filter -v '     1 && 2 && 3'   /dev/null

run_mlr put    -v '$x = 1  == 2  == 3' /dev/null
run_mlr filter -v '     1  == 2  == 3' /dev/null
run_mlr put    -v '$x = 1  != 2  != 3' /dev/null
run_mlr filter -v '     1  != 2  != 3' /dev/null
run_mlr put    -v '$x = 1  =~ 2  =~ 3' /dev/null
run_mlr filter -v '     1  =~ 2  =~ 3' /dev/null
run_mlr put    -v '$x = 1 !=~ 2 !=~ 3' /dev/null
run_mlr filter -v '     1 !=~ 2 !=~ 3' /dev/null
run_mlr put    -v '$x = 1  == 2  != 3' /dev/null
run_mlr filter -v '     1  == 2  != 3' /dev/null
run_mlr put    -v '$x = 1  != 2  == 3' /dev/null
run_mlr filter -v '     1  != 2  == 3' /dev/null

run_mlr put    -v '$x = 1  <  2  <  3' /dev/null
run_mlr filter -v '     1  <  2  <  3' /dev/null
run_mlr put    -v '$x = 1  <= 2  <= 3' /dev/null
run_mlr filter -v '     1  <= 2  <= 3' /dev/null
run_mlr put    -v '$x = 1  >  2  >  3' /dev/null
run_mlr filter -v '     1  >  2  >  3' /dev/null
run_mlr put    -v '$x = 1  >= 2  >= 3' /dev/null
run_mlr filter -v '     1  >= 2  >= 3' /dev/null
run_mlr put    -v '$x = 1  <  2  <= 3' /dev/null
run_mlr filter -v '     1  <  2  <= 3' /dev/null
run_mlr put    -v '$x = 1  <= 2  <  3' /dev/null
run_mlr filter -v '     1  <= 2  <  3' /dev/null

run_mlr put    -v '$x = 1 |  2 |  3'   /dev/null
run_mlr filter -v '     1 |  2 |  3'   /dev/null
run_mlr put    -v '$x = 1 ^  2 ^  3'   /dev/null
run_mlr filter -v '     1 ^  2 ^  3'   /dev/null
run_mlr put    -v '$x = 1 &  2 &  3'   /dev/null
run_mlr filter -v '     1 &  2 &  3'   /dev/null
run_mlr put    -v '$x = 1 |  2 &  3'   /dev/null
run_mlr filter -v '     1 |  2 &  3'   /dev/null
run_mlr put    -v '$x = 1 |  2 ^  3'   /dev/null
run_mlr filter -v '     1 |  2 ^  3'   /dev/null
run_mlr put    -v '$x = 1 ^  2 |  3'   /dev/null
run_mlr filter -v '     1 ^  2 |  3'   /dev/null
run_mlr put    -v '$x = 1 ^  2 &  3'   /dev/null
run_mlr filter -v '     1 ^  2 &  3'   /dev/null
run_mlr put    -v '$x = 1 &  2 |  3'   /dev/null
run_mlr filter -v '     1 &  2 |  3'   /dev/null
run_mlr put    -v '$x = 1 &  2 ^  3'   /dev/null
run_mlr filter -v '     1 &  2 ^  3'   /dev/null

run_mlr put    -v '$x = 1  << 2  << 3' /dev/null
run_mlr filter -v '     1  << 2  << 3' /dev/null
run_mlr put    -v '$x = 1  >> 2  >> 3' /dev/null
run_mlr filter -v '     1  >> 2  >> 3' /dev/null
run_mlr put    -v '$x = 1  << 2  >> 3' /dev/null
run_mlr filter -v '     1  << 2  >> 3' /dev/null
run_mlr put    -v '$x = 1  >> 2  << 3' /dev/null
run_mlr filter -v '     1  >> 2  << 3' /dev/null

run_mlr put    -v '$x = 1 + 2 + 3'   /dev/null
run_mlr filter -v '     1 + 2 + 3'   /dev/null
run_mlr put    -v '$x = 1 - 2 - 3'   /dev/null
run_mlr filter -v '     1 - 2 - 3'   /dev/null
run_mlr put    -v '$x = 1 + 2 - 3'   /dev/null
run_mlr filter -v '     1 + 2 - 3'   /dev/null
run_mlr put    -v '$x = 1 - 2 + 3'   /dev/null
run_mlr filter -v '     1 - 2 + 3'   /dev/null
run_mlr put    -v '$x = 1 . 2 . 3'   /dev/null
run_mlr filter -v '     1 . 2 . 3'   /dev/null

run_mlr put    -v '$x = 1 * 2 * 3'   /dev/null
run_mlr filter -v '     1 * 2 * 3'   /dev/null
run_mlr put    -v '$x = 1 / 2 / 3'   /dev/null
run_mlr filter -v '     1 / 2 / 3'   /dev/null
run_mlr put    -v '$x = 1 // 2 // 3' /dev/null
run_mlr filter -v '     1 // 2 // 3' /dev/null
run_mlr put    -v '$x = 1 % 2 % 3'   /dev/null
run_mlr filter -v '     1 % 2 % 3'   /dev/null
run_mlr put    -v '$x = 1 ** 2 ** 3' /dev/null
run_mlr filter -v '     1 ** 2 ** 3' /dev/null


run_mlr put    -v '$x = 1 *  2 /  3'   /dev/null
run_mlr filter -v '     1 *  2 /  3'   /dev/null
run_mlr put    -v '$x = 1 *  2 // 3'   /dev/null
run_mlr filter -v '     1 *  2 // 3'   /dev/null
run_mlr put    -v '$x = 1 *  2 %  3'   /dev/null
run_mlr filter -v '     1 *  2 %  3'   /dev/null
run_mlr put    -v '$x = 1 *  2 ** 3'   /dev/null
run_mlr filter -v '     1 *  2 ** 3'   /dev/null

run_mlr put    -v '$x = 1 /  2 *  3'   /dev/null
run_mlr filter -v '     1 /  2 *  3'   /dev/null
run_mlr put    -v '$x = 1 /  2 // 3'   /dev/null
run_mlr filter -v '     1 /  2 // 3'   /dev/null
run_mlr put    -v '$x = 1 /  2 %  3'   /dev/null
run_mlr filter -v '     1 /  2 %  3'   /dev/null
run_mlr put    -v '$x = 1 /  2 ** 3'   /dev/null
run_mlr filter -v '     1 /  2 ** 3'   /dev/null

run_mlr put    -v '$x = 1 // 2 *  3'   /dev/null
run_mlr filter -v '     1 // 2 *  3'   /dev/null
run_mlr put    -v '$x = 1 // 2 /  3'   /dev/null
run_mlr filter -v '     1 // 2 /  3'   /dev/null
run_mlr put    -v '$x = 1 // 2 %  3'   /dev/null
run_mlr filter -v '     1 // 2 %  3'   /dev/null
run_mlr put    -v '$x = 1 // 2 ** 3'   /dev/null
run_mlr filter -v '     1 // 2 ** 3'   /dev/null

run_mlr put    -v '$x = 1 %  2 *  3'   /dev/null
run_mlr filter -v '     1 %  2 *  3'   /dev/null
run_mlr put    -v '$x = 1 %  2 /  3'   /dev/null
run_mlr filter -v '     1 %  2 /  3'   /dev/null
run_mlr put    -v '$x = 1 %  2 // 3'   /dev/null
run_mlr filter -v '     1 %  2 // 3'   /dev/null
run_mlr put    -v '$x = 1 %  2 ** 3'   /dev/null
run_mlr filter -v '     1 %  2 ** 3'   /dev/null

run_mlr put    -v '$x = 1 ** 2 *  3'   /dev/null
run_mlr filter -v '     1 ** 2 *  3'   /dev/null
run_mlr put    -v '$x = 1 ** 2 /  3'   /dev/null
run_mlr filter -v '     1 ** 2 /  3'   /dev/null
run_mlr put    -v '$x = 1 ** 2 // 3'   /dev/null
run_mlr filter -v '     1 ** 2 // 3'   /dev/null
run_mlr put    -v '$x = 1 ** 2 %  3'   /dev/null
run_mlr filter -v '     1 ** 2 %  3'   /dev/null

run_mlr put    -v '$x = ++1'   /dev/null
run_mlr filter -v '     ++1'   /dev/null
run_mlr put    -v '$x = --1'   /dev/null
run_mlr filter -v '     --1'   /dev/null
run_mlr put    -v '$x = !!1'   /dev/null
run_mlr filter -v '     !!1'   /dev/null
run_mlr put    -v '$x = ~~1'   /dev/null
run_mlr filter -v '     ~~1'   /dev/null

run_mlr put    -v '$x = 1 ? 2 : 3'         /dev/null
run_mlr filter -v '     1 ? 2 : 3'         /dev/null
run_mlr put    -v '$x = 1 ? 2 ? 3 : 4 : 5' /dev/null
run_mlr filter -v '     1 ? 2 ? 3 : 4 : 5' /dev/null
run_mlr put    -v '$x = 1 ? 2 : 3 ? 4 : 5' /dev/null
run_mlr filter -v '     1 ? 2 : 3 ? 4 : 5' /dev/null
