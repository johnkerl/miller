# These should produce no warnings
run_mlr -n put    'y =  1'
run_mlr -n put -w 'y =  1'
run_mlr -n put -W 'y =  1'

run_mlr -n put    'x = 3; y = x'
run_mlr -n put -w 'x = 3; y = x'
run_mlr -n put -W 'x = 3; y = x'

run_mlr -n put -W 'for (k in $*) { print k }'
run_mlr -n put -W 'for (k,v in $*) { print k,v }'
run_mlr -n put -W 'for ((k1,k2),v in $*) { print k1,k2,v }'

# Presence of $x not knowable until runtime; only catchable by a 'strict mode'.
run_mlr -n put    'y = $x'
run_mlr -n put -w 'y = $x'
run_mlr -n put -W 'y = $x'

# x should be flagged as an uninitialized read
run_mlr         -n put    'y = x'
run_mlr         -n put -w 'y = x'
mlr_expect_fail -n put -W 'y = x'

# Both x and y should be flagged as uninitialized reads
run_mlr         -n put    'z = x + y'
run_mlr         -n put -w 'z = x + y'
mlr_expect_fail -n put -W 'z = x + y'

# This should produce no warnings
run_mlr -n put -W 'i = 0; z[i] = 1'

# i should be flagged as an uninitialized read
mlr_expect_fail -n put -W 'z[i] = 1'

# This should produce no warnings
run_mlr -n put -W 'func f(n) { return n+1 }'

# m should be flagged as an uninitialized read
mlr_expect_fail -n put -W 'func f(n) { return m+1 }'
mlr_expect_fail -n put -W 'm = 0; func f(n) { return m+1 }'
mlr_expect_fail -n put -W 'subr f(n) { print m+1 }'
mlr_expect_fail -n put -W 'm = 0; subr f(n) { print m+1 }'

# The uninit-warner isn't too smart. We expect this to not warn.  (It would be
# great if someday it got smarter and warned here -- this test would need to be
# re-accepted.)
run_mlr -n put -W 'if (false) {x = 1}; print x'
