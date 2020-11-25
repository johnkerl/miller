# ----------------------------------------------------------------
announce USER-DEFINED FUNCTIONS

run_mlr --opprint --from $indir/abixy put 'func f(u,v){return u+v} $o=f(NR*1000,$x)'

mlr_expect_fail --opprint --from $indir/abixy put 'func f(x,y,z){$nnn=999; int n=10; return $y} $o=f($i,$x,$y)'

# general programming-language stuff
run_mlr -n put -q -f $indir/sieve.mlr
run_mlr -n put -q -f $indir/mand.mlr -e 'begin {@verbose = true}'

# not use all args (for valgrind)
run_mlr --from $indir/abixy put 'func f(x,y) { return 2*y } $o = f($x,$y) '

# Test variable-clear at scope exit; test read of unset locals.
run_mlr --opprint --from $indir/abixy put '$o1 = a; int a=1000+NR; $o2 = a; a=2000+NR; $o3 = a'

# Test recursion
run_mlr --opprint --from $indir/abixy put '
    func f(n) {
        if (is_numeric(n)) {
            if (n > 0) {
                return n * f(n-1)
            } else {
                return 1
            }
        }
        # implicitly return absent
    }
    $o = f(NR)
'

run_mlr --from $indir/abixy --opprint put '
  func f(n) {
    return n+1;
  }
  $o1 = f(NR);
  $o2 = f(f(NR));
  $o3 = f(f(f(NR)));
  $o4 = f(f(f(f(NR))));
  $o5 = f(f(f(f(f(NR)))));
  $o6 = f(f(f(f(f(f(NR))))));
'

run_mlr --from $indir/abixy --opprint put '
  func f(n) {
      if (n < 2) {
          return 1
      } else {
          return f(n-1) + f(n-2)
      }
  }
  $o = f(NR)
'

run_mlr --from $indir/abixy --opprint put '
  func f(n) {
    str sn = string(n); # for map keys
    if (is_present(@fcache[sn])) {
      return @fcache[sn]
    } else {
      num rv = 1;
      if (n >= 2) {
        rv = f(n-1) + f(n-2)
      }
      @fcache[sn] = rv;
      return rv
    }
  }
  $o = f(NR)
'
