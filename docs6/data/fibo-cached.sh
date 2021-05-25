mlr --ofmt '%.9lf' --opprint seqgen --start 1 --stop 28 then put '
  func f(n) {
    @fcount += 1;                 # count number of calls to the function
    if (is_present(@fcache[n])) { # cache hit
      return @fcache[n]
    } else {                      # cache miss
      num rv = 1;
      if (n >= 2) {
        rv = f(n-1) + f(n-2)      # recurse
      }
      @fcache[n] = rv;
      return rv
    }
  }
  @fcount = 0;
  $o = f($i);
  $fcount = @fcount;
' then put '$seconds=systime()' then step -a delta -f seconds then cut -x -f seconds
