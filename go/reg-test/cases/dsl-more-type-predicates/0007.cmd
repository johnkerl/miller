mlr --from reg-test/input/abixy-het put -q 'for (k,v in $*) {if (is_bool(NR==2)) {@bool[NR][k]    = "NR==2"}} end{ emit @bool,    "NR", "k" }'
