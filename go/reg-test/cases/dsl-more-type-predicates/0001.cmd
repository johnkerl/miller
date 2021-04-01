mlr --from reg-test/input/abixy --opprint put ' for (k, v in $*) { $[k."_type"]      = typeof(v)     } '
