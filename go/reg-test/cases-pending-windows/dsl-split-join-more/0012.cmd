mlr --from reg-test/input/abixy-het put '$* = splitkv("a=1,b=2,c", OPS, OFS); print ">>".ORS."<<"; for (k, v in $*) {print k.":".typeof(k)." ".v.":".typeof(v)}'
