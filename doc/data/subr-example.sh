mlr --opprint --from data/small put -q '
  begin {
    @call_count = 0;
  }
  subr s(n) {
    @call_count += 1;
    if (is_numeric(n)) {
      if (n > 1) {
        call s(n-1);
      } else {
        print "numcalls=" . @call_count;
      }
    }
  }
  print "NR=" . NR;
  call s(NR);
'
