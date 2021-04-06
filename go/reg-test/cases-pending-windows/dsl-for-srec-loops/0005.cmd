mlr --from reg-test/input/abixy --opprint put -v 'for(k,v in $*) {$[k."_orig"]=v; $[k] = "other"}'
