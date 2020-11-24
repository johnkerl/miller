# ----------------------------------------------------------------
announce MAPVARS IN SCALAR FUNCTION-CALL CONTEXTS

run_mlr --from $indir/abixy put '$z=strlen($*)'
run_mlr --from $indir/abixy put '$z=strlen({})'
run_mlr --from $indir/abixy put 'a={}; $z=strlen(a)'
