# ----------------------------------------------------------------
announce USER-DEFINED SUBROUTINES

# Test recursion
run_mlr --opprint --from $indir/abixy head -n 5 then put '
    subr s(n) {
        print "n = " . n;
        if (is_numeric(n)) {
            if (n > 0) {
                call s(n-1)
            }
        }
    }
    print "";
    call s(NR)
'
