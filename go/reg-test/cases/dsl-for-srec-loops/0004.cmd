mlr --from reg-test/input/abixy put -v 'for(k,v in $*) {if (k != "x") {unset $[k]}}; $j = NR'
