mlr --from reg-test/input/abixy-het put -q 'for (k,v in $*) {if (is_bool(v))     {@bool[NR][k]    = v}}    end{ emit @bool,    "NR", "k" }'
