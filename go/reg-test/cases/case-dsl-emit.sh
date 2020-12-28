# ================================================================
run_mlr --from $indir/s.dkvp --ojson put -q '@sum += $i; emit {"sum": @sum}'
run_mlr --from $indir/s.dkvp --ojson put -q '@sum[$a] += $i; emit {"sum": @sum}'

# ================================================================
announce "1. NON-LASHED NON-INDEXED NAMEDVAR EMITP"
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emitp @sumx;
  }
'

announce "1. NON-LASHED NON-INDEXED MAP EMITP"
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emitp @sums;
  }
'

announce "1. NON-LASHED NON-INDEXED NAMEDVAR EMIT"
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emit  @sumx;
  }
'

announce "1. NON-LASHED NON-INDEXED MAP EMIT"
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emit  @sums;
  }
'

# ================================================================
announce "2. LASHED NON-INDEXED NAMEDVAR EMITP"
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emitp (@sumx, @countx);
  }
'

announce "2. LASHED NON-INDEXED MAP EMITP"
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emitp (@sums, @counts);
  }
'

announce "2. LASHED NON-INDEXED NAMEDVAR EMIT"
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emit  (@sumx, @countx);
  }
'

announce "2. LASHED NON-INDEXED MAP EMIT"
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emit  (@sums, @counts);
  }
'

# ================================================================
announce "3. NON-LASHED INDEXED MAP EMITP"
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a] += $x;
  @counts[$a] += 1;
  end {
    emitp @sums, "a";
  }
'

announce "3. NON-LASHED INDEXED MAP EMIT"
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a] += $x;
  @counts[$a] += 1;
  end {
    emit  @sums, "a";
  }
'

announce "3. NON-LASHED UNDER-INDEXED MAP EMITP"
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emitp @sums, "a";
  }
'

announce "3. NON-LASHED AT-INDEXED MAP EMITP"
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emitp @sums, "a", "b";
  }
'

announce "3. NON-LASHED OVER-INDEXED MAP EMITP"
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emitp @sums, "a", "b", "c";
  }
'

announce "3. NON-LASHED UNDER-INDEXED MAP EMIT"
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emit  @sums, "a";
  }
'

announce "3. NON-LASHED AT-INDEXED MAP EMIT"
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emit  @sums, "a", "b";
  }
'

announce "3. NON-LASHED OVER-INDEXED MAP EMIT"
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emit  @sums, "a", "b", "c";
  }
'

# ================================================================
announce "4. LASHED INDEXED MAP EMITP"
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a] += $x;
  @counts[$a] += 1;
  end {
    emitp (@sums, @counts), "a";
  }
'

announce "4. LASHED INDEXED MAP EMIT"
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a] += $x;
  @counts[$a] += 1;
  end {
    emit  (@sums, @counts), "a";
  }
'

announce "4. LASHED UNDER-INDEXED MAP EMITP"
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emitp (@sums, @counts), "a";
  }
'

announce "4. LASHED AT-INDEXED MAP EMITP"
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emitp (@sums, @counts), "a", "b";
  }
'

announce "4. LASHED OVER-INDEXED MAP EMITP"
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emitp (@sums, @counts), "a", "b", "c";
  }
'

announce "4. LASHED UNDER-INDEXED MAP EMIT"
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emit  (@sums, @counts), "a";
  }
'

announce "4. LASHED AT-INDEXED MAP EMIT"
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emit  (@sums, @counts), "a", "b";
  }
'

announce "4. LASHED OVER-INDEXED MAP EMIT"
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emit  (@sums, @counts), "a", "b", "c";
  }
'
