
# This important test validates the local-stack allocator: which variables are
# assigned which offsets in the stack, and how the local-extent contract is
# satisfied by the clear-at-enter-subframe logic.

run_mlr --from $indir/abixy put -f $indir/test-dsl-stack-allocation.mlr

# These test absent-null handing for as-yet-undefined localvars in expression RHSs.
run_mlr --from $indir/abixy put 'a=a; $oa = a'
run_mlr --from $indir/abixy put 'a=a; $oa = a; a = a; $ob = a'
run_mlr --from $indir/abixy put 'a=a; $oa = a; a = a; $ob = a; a = b; $oc = a'
run_mlr --from $indir/abixy put 'a=a; $oa = a; a = a; $ob = a; a = b; $oc = a; b = 3; b = a; $od = a'
run_mlr --from $indir/abixy put 'a=a; $oa = a; a = a; $ob = a; a = b; $oc = a; b = 3; b = a; $od = a; b = 4;a = b; $oe= a'
