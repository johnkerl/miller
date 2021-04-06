mlr --from reg-test/input/abixy-het put -q 'for (k,v in $*) {if (is_int(v))      {@int[NR][k]     = v}}    end{ emit @int,     "NR", "k" }'
