mlr seqgen --start 1 --stop 10 then put '
  func f(a, b) {
      r = 0.0;
      for (local i = 0; i < 6; i += 1) {
          local u = urand();
          r += u;
      }
      r /= 6;
      return a + (b - a) * r;
  }
  local o = f(10, 20);
  $o = o;
'
