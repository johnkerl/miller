# Check for BNF/AST errors, with minimal CST involvement
mlr_expect_fail -n put -v 'call s()'
mlr_expect_fail -n put -v 'call s'
mlr_expect_fail -n put -v 'call s(1,2,3)'
run_mlr -n put -v 'subr s() {}'
run_mlr -n put -v 'subr s() {x=1}'
run_mlr -n put -v 'subr s() {return}'
mlr_expect_fail -n put -v 'subr s() {return 2}'
mlr_expect_fail -n put 'subr s()'
run_mlr -n put -v 'subr s() {}; call s()'
run_mlr -n put -v 'call s(); subr s() {}'

run_mlr -n put -v 'subr s(a) {print "HELLO, ".a."!"} call s("WORLD")'

# Check for CST invovlement
run_mlr --from $indir/2.dkvp put 'subr s() {}'
run_mlr --from $indir/2.dkvp put 'subr s() {return}'

mlr_expect_fail --from $indir/2.dkvp put 'call s()'

mlr_expect_fail -n put 'subr s(){} subr s(){}'

run_mlr --from $indir/2.dkvp put -v -q '
  func s(x) {
    return x*2;
  }
  subr s(a) {
    print "HELLO, ".a."!"
  }
  print s(NR);
  call s("WORLD");
'

# Test recursion
run_mlr --opprint --from $indir/abixy head -n 5 then put '
    subr s(n) {
        print "n = " . n;
        if (is_numeric(n)) {
            if (n > 0) {
                call s(n-1)
            }
        }
    }
    print "";
    call s(NR)
'
