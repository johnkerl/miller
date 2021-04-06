mlr --from reg-test/input/abixy-het put -q 'for (k,v in $*) {if (is_float(v))    {@float[NR][k]   = v}}    end{ emit @float,   "NR", "k" }'
