mlr --from reg-test/input/s.dkvp head -n 2 then put -q 'for (k,v in $*) { emit { "foo" : v } }'
