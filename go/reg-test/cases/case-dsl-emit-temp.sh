announce "1. NON-LASHED NON-INDEXED NAMEDVAR"

run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a] += $x;
  @counts[$a] += 1;
  end {
    emitp @sumx;
    emit  @sumx;
  }
'

announce "1. NON-LASHED NON-INDEXED MAP"

run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a] += $x;
  @counts[$a] += 1;
  end {
    emitp @sums;
    emit  @sums;
  }
'

announce "2. LASHED NON-INDEXED NAMEDVAR"
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a] += $x;
  @counts[$a] += 1;
  end {
    emitp (@sumx, @countx);
    emit  (@sumx, @countx);
  }
'

announce "2. LASHED NON-INDEXED MAP"
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a] += $x;
  @counts[$a] += 1;
  end {
    emitp (@sums, @counts);
    emit  (@sums, @counts);
  }
'

announce "3. NON-LASHED INDEXED MAP"
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a] += $x;
  @counts[$a] += 1;
  end {
    emitp @sums, "a";
    emit  @sums, "a";
  }
'

announce "4. LASHED INDEXED MAP"
run_mlr --ojson --jvstack --from $indir/abixy put -q '
  @sumx += $x;
  @countx += 1;
  @sums[$a] += $x;
  @counts[$a] += 1;
  end {
    emitp (@sums, @counts), "a";
    emit  (@sums, @counts), "a";
  }
'




