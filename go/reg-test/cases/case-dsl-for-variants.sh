run_mlr --from $indir/s.dkvp put 'for (@i = 0; @i < NR; @i += 1) { $i += @i }'
run_mlr --from $indir/s.dkvp put 'i=999; for (i = 0; i < NR; i += 1) { $i += i }'
run_mlr --from $indir/s.dkvp put -v 'j = 2; for (i = 0; i < NR; i += 1) { $i += i }'

# The middle 'continuation' statement must be:
# * empty, in which case it evaluates to true;
# * one statement, which must be a bare-boolean;
# * more than one, in which case all but the last must be assignments and the
#   last must be a bare boolean.

# Zero continuations: true
run_mlr --from $indir/s.dkvp head -n 2 then put '
  for (int i = 0; ; i += 1) {
    print i;
    if (i > 5) {
      break
    }
  }
'

# One continuations: is bare-boolean
run_mlr --from $indir/s.dkvp head -n 2 then put '
  for (int i = 0; i < 5 ; i += 1) {
    print i;
  }
'

# One continuations: is not bare-boolean
mlr_expect_fail --from $indir/s.dkvp head -n 2 then put '
  for (int i = 0; j = 5 ; i += 1) {
    print i;
    if (i > 5) {
      brea
    }
  }
'

# Two continuations: assignment, bare-boolean
run_mlr --from $indir/s.dkvp head -n 2 then put '
  j = 20;
  for (int i = 0; j += 1, i < 5 ; i += 1) {
    print i;
    print j;
  }
'

# Two continuations: two assignments
mlr_expect_fail --from $indir/s.dkvp head -n 2 then put '
  j = 20;
  for (int i = 0; j += 1, i += 5 ; i += 1) {
    print i;
    if (i > 5) {
      break
    }
  }
'

# Two continuations: two bare-booleans
mlr_expect_fail --from $indir/s.dkvp head -n 2 then put '
  j = 20;
  for (int i = 0; j < 10, i < 10 ; i += 1) {
    print i;
    if (i > 5) {
      break
    }
  }
'


# Two continuations: bare-boolean, assignment
mlr_expect_fail --from $indir/s.dkvp head -n 2 then put '
  j = 20;
  for (int i = 0; i < 10, j = 10 ; i += 1) {
    print i;
    if (i > 5) {
      break
    }
  }
'

