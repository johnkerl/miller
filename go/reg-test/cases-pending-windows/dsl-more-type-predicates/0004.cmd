mlr --from reg-test/input/abixy-het put -q 'for (k,v in $*) {if (is_numeric(v))  {@numeric[NR][k] = v}}    end{ emit @numeric, "NR", "k" }'
