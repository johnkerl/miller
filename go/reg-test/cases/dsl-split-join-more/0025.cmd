mlr --oxtab --from reg-test/input/abixy-het put 's = joinkv({1:2, "abc":4, 5:"xyz"}, ":", ";"); $* = splitkv(s, ":", ";"); for (k,v in $*) { print k.":".typeof(k)." ".v.":".typeof(v)}'
