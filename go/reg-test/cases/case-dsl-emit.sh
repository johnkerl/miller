# ================================================================
run_mlr --from $indir/s.dkvp --ojson put -q '@sum += $i; emit {"sum": @sum}'
run_mlr --from $indir/s.dkvp --ojson put -q '@sum[$a] += $i; emit {"sum": @sum}'

# ================================================================
announce "1. NON-LASHED NON-INDEXED NAMEDVAR EMITP"
for path_to_mlr in "../c/mlr" "../go/mlr"; do
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emitp @sumx;
  }
'
done

announce "1. NON-LASHED NON-INDEXED MAP EMITP"
for path_to_mlr in "../c/mlr" "../go/mlr"; do
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emitp @sums;
  }
'
done

announce "1. NON-LASHED NON-INDEXED NAMEDVAR EMIT"
for path_to_mlr in "../c/mlr" "../go/mlr"; do
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emit  @sumx;
  }
'
done

announce "1. NON-LASHED NON-INDEXED MAP EMIT"
for path_to_mlr in "../c/mlr" "../go/mlr"; do
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emit  @sums;
  }
'
done

# ================================================================
announce "2. LASHED NON-INDEXED NAMEDVAR EMITP"
for path_to_mlr in "../c/mlr" "../go/mlr"; do
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emitp (@sumx, @countx);
  }
'
done

announce "2. LASHED NON-INDEXED MAP EMITP"
for path_to_mlr in "../c/mlr" "../go/mlr"; do
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emitp (@sums, @counts);
  }
'
done

announce "2. LASHED NON-INDEXED NAMEDVAR EMIT"
for path_to_mlr in "../c/mlr" "../go/mlr"; do
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emit  (@sumx, @countx);
  }
'
done

announce "2. LASHED NON-INDEXED MAP EMIT"
for path_to_mlr in "../c/mlr" "../go/mlr"; do
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emit  (@sums, @counts);
  }
'
done

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
for path_to_mlr in "../c/mlr" "../go/mlr"; do
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emitp @sums, "a";
  }
'
done

announce "3. NON-LASHED AT-INDEXED MAP EMITP"
for path_to_mlr in "../c/mlr" "../go/mlr"; do
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emitp @sums, "a", "b";
  }
'
done

announce "3. NON-LASHED OVER-INDEXED MAP EMITP"
for path_to_mlr in "../c/mlr" "../go/mlr"; do
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emitp @sums, "a", "b", "c";
  }
'
done

announce "3. NON-LASHED UNDER-INDEXED MAP EMIT"
for path_to_mlr in "../c/mlr" "../go/mlr"; do
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emit  @sums, "a";
  }
'
done

announce "3. NON-LASHED AT-INDEXED MAP EMIT"
for path_to_mlr in "../c/mlr" "../go/mlr"; do
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emit  @sums, "a", "b";
  }
'
done

announce "3. NON-LASHED OVER-INDEXED MAP EMIT"
for path_to_mlr in "../c/mlr" "../go/mlr"; do
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emit  @sums, "a", "b", "c";
  }
'
done

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
for path_to_mlr in "../c/mlr" "../go/mlr"; do
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emitp (@sums, @counts), "a";
  }
'
done

announce "4. LASHED AT-INDEXED MAP EMITP"
for path_to_mlr in "../c/mlr" "../go/mlr"; do
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emitp (@sums, @counts), "a", "b";
  }
'
done

announce "4. LASHED OVER-INDEXED MAP EMITP"
for path_to_mlr in "../c/mlr" "../go/mlr"; do
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emitp (@sums, @counts), "a", "b", "c";
  }
'
done

announce "4. LASHED UNDER-INDEXED MAP EMIT"
for path_to_mlr in "../c/mlr" "../go/mlr"; do
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emit  (@sums, @counts), "a";
  }
'
done

announce "4. LASHED AT-INDEXED MAP EMIT"
for path_to_mlr in "../c/mlr" "../go/mlr"; do
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emit  (@sums, @counts), "a", "b";
  }
'
done

announce "4. LASHED OVER-INDEXED MAP EMIT"
for path_to_mlr in "../c/mlr" "../go/mlr"; do
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a][$b] += $x;
  @counts[$a][$b] += 1;
  end {
    emit  (@sums, @counts), "a", "b", "c";
  }
'
done
