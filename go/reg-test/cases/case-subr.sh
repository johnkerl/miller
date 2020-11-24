# Check for BNF/AST errors, with minimal CST involvement
run_mlr -n put -v 'call s()'
run_mlr -n put -v 'call s(1,2,3)'
run_mlr -n put -v 'subr s() {}'
run_mlr -n put -v 'subr s() {x=1}'
run_mlr -n put -v 'subr s() {return}'

# Check for CST invovlement
