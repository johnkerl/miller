# Check for BNF/AST errors, with minimal CST involvement
run_mlr -n put -v 'call s()'
mlr_expect_fail -n put -v 'call s'
run_mlr -n put -v 'call s(1,2,3)'
run_mlr -n put -v 'subr s() {}'
run_mlr -n put -v 'subr s() {x=1}'
run_mlr -n put -v 'subr s() {return}'
mlr_expect_fail -n put -v 'subr s() {return 2}'
mlr_expect_fail -n put 'subr s()'

run_mlr -n put -v 'subr s(a) {print "HELLO, ".a."!"} call s("WORLD")'

# Check for CST invovlement
run_mlr --from $indir/2.dkvp put 'subr s() {}'
run_mlr --from $indir/2.dkvp put 'subr s() {return}'

mlr_expect_fail --from $indir/2.dkvp put 'call s()'

run_mlr --from $indir/2.dkvp put 'subr s(a) {print "HELLO, ".a."!"} call s("WORLD")'

#run_mlr --from $indir/2.dkvp put 'call s()'
