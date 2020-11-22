# ----------------------------------------------------------------
announce MAPVARS IN SCALAR FUNCTION-CALL CONTEXTS

mlr_expect_fail --from $indir/abixy put '$z=strlen($*)'
mlr_expect_fail --from $indir/abixy put '$z=strlen({})'
run_mlr --from $indir/abixy put 'a={}; $z=strlen(a)'
