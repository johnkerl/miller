mlr seqgen --start 1 --stop 10 then put '
  func f(a, b) {                          # function arguments a and b
      r = 0.0;                            # local r scoped to the function
      for (local i = 0; i < 6; i += 1) {  # local i scoped to the for-loop
          local u = urand();              # local u scoped to the for-loop
          r += u;                         # updates r from the enclosing scope
      }
      r /= 6;
      return a + (b - a) * r;
  }
  local o = f(10, 20);                    # local to the top-level scope
  $o = o;
'
