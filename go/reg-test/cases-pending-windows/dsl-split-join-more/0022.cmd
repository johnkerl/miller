mlr --oxtab --from reg-test/input/abixy-het put 's = joinkv($*, ":", ";"); $* = splitkvx(s, ":", ";"); for (k,v in $*) { print k.":".typeof(k)." ".v.":".typeof(v)}'
