
run_mlr --from $indir/abixy --opprint put ' for (k, v in $*) { $[k."_type"]      = typeof(v)     } '

run_mlr --from $indir/abixy-het put -q 'for (k,v in $*) {if (is_float(v))    {@float[NR][k]   = v}}    end{ emit @float,   "NR", "k" }'
run_mlr --from $indir/abixy-het put -q 'for (k,v in $*) {if (is_int(v))      {@int[NR][k]     = v}}    end{ emit @int,     "NR", "k" }'
run_mlr --from $indir/abixy-het put -q 'for (k,v in $*) {if (is_numeric(v))  {@numeric[NR][k] = v}}    end{ emit @numeric, "NR", "k" }'
run_mlr --from $indir/abixy-het put -q 'for (k,v in $*) {if (is_string(v))   {@string[NR][k]  = v}}    end{ emit @string,  "NR", "k" }'
run_mlr --from $indir/abixy-het put -q 'for (k,v in $*) {if (is_bool(v))     {@bool[NR][k]    = v}}    end{ emit @bool,    "NR", "k" }'
run_mlr --from $indir/abixy-het put -q 'for (k,v in $*) {if (is_bool(NR==2)) {@bool[NR][k]    = "NR==2"}} end{ emit @bool,    "NR", "k" }'
