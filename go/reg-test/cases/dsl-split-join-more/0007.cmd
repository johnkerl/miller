mlr --from reg-test/input/abixy-het put '$* = splitnv("a,b,c" , ","); for (k,v in $*) { print k.":".typeof(k)." ".v.":".typeof(v)}'
