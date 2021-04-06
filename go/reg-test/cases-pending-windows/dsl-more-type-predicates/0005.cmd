mlr --from reg-test/input/abixy-het put -q 'for (k,v in $*) {if (is_string(v))   {@string[NR][k]  = v}}    end{ emit @string,  "NR", "k" }'
