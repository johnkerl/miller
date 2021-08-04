mlr --ofmt '%.9lf' --opprint seqgen --start 1 --stop 28 then put '
  func f(n) {
      @fcount += 1;              # count number of calls to the function
      if (n < 2) {
          return 1
      } else {
          return f(n-1) + f(n-2) # recurse
      }
  }

  @fcount = 0;
  $o = f($i);
  $fcount = @fcount;

' then put '$seconds=systime()' then step -a delta -f seconds then cut -x -f seconds
