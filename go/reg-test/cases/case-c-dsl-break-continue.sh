# ----------------------------------------------------------------
announce DSL BREAK/CONTINUE IN SINGLE WHILE/DO-WHILE

run_mlr --opprint --from $indir/abixy put '
  while ($i < 5) {
    $i += 1;
    break;
    $a = "ERROR";
  }
'

run_mlr --opprint --from $indir/abixy put '
  while ($i < 5) {
    $i += 1;
    continue;
    $a = "ERROR";
  }
'

run_mlr --opprint --from $indir/abixy put '
  do {
    $i += 1;
    break;
    $a = "ERROR";
  } while ($i < 5);
'

run_mlr --opprint --from $indir/abixy put '
  do {
    $i += 1;
    continue;
    $a = "ERROR";
  } while ($i < 5);
'

run_mlr --opprint --from $indir/abixy put '
  $NR = NR;
  while ($i < 5) {
    $i += 1;
    if (NR == 2) {
      break;
    }
    $a = "reached";
  }
'

run_mlr --opprint --from $indir/abixy put '
  $NR = NR;
  while ($i < 5) {
    $i += 1;
    if (NR == 2) {
      continue;
    }
    $a = "reached";
  }
'

run_mlr --opprint --from $indir/abixy put '
$NR = NR;
  do {
    $i += 1;
    if (NR == 2) {
      break;
    }
    $a = "reached";
  } while ($i < 5);
'

run_mlr --opprint --from $indir/abixy put '
  $NR = NR;
  do {
    $i += 1;
    if (NR == 2) {
      continue;
    }
    $a = "reached";
  } while ($i < 5);
'

# ----------------------------------------------------------------
announce DSL BREAK/CONTINUE IN NESTED WHILE/DO-WHILE

run_mlr --opprint --from $indir/abixy put '
  $j = NR;
  while ($j < 4) {
    $k = NR;
    $j += 1;
    break;
    while ($k < 7) {
      $k += 1
    }
  }
'

run_mlr --opprint --from $indir/abixy put '
  $j = NR;
  while ($j < 4) {
    $k = NR;
    $j += 1;
    continue;
    while ($k < 7) {
      $k += 1
    }
  }
'

run_mlr --opprint --from $indir/abixy put '
  $j = NR;
  while ($j < 4) {
    $k = NR;
    $j += 1;
    while ($k < 7) {
      $k += 1;
      break;
      $k += 10000;
    }
  }
'

run_mlr --opprint --from $indir/abixy put '
  $j = NR;
  while ($j < 4) {
    $k = NR;
    $j += 1;
    while ($k < 7) {
      $k += 1;
      continue;
      $k += 10000;
    }
  }
'


run_mlr --opprint --from $indir/abixy put '
  $j = NR;
  while ($j < 4) {
    $k = NR;
    $j += 1;
    if (NR == 2 || NR == 8) {
      break;
    }
    while ($k < 7) {
      $k += 1
    }
  }
'

run_mlr --opprint --from $indir/abixy put '
  $j = NR;
  while ($j < 4) {
    $k = NR;
    $j += 1;
    if (NR == 2 || NR == 8) {
      continue;
    }
    while ($k < 7) {
      $k += 1
    }
  }
'

run_mlr --opprint --from $indir/abixy put '
  $j = NR;
  while ($j < 4) {
    $k = NR;
    $j += 1;
    while ($k < 7) {
      $k += 1;
      if (NR == 2 || NR == 8) {
        break;
      }
      $k += 10000;
    }
  }
'

run_mlr --opprint --from $indir/abixy put '
  $j = NR;
  while ($j < 4) {
    $k = NR;
    $j += 1;
    while ($k < 7) {
      $k += 1;
      if (NR == 2 || NR == 8) {
        continue;
      }
      $k += 10000;
    }
  }
'


run_mlr --opprint --from $indir/abixy put '
  $j = NR;
  do {
    $k = NR;
    $j += 1;
    break;
    do {
      $k += 1
    } while ($k < 7);
  } while ($j < 4);
'

run_mlr --opprint --from $indir/abixy put '
  $j = NR;
  do {
    $k = NR;
    $j += 1;
    continue;
    do {
      $k += 1
    } while ($k < 7);
  } while ($j < 4);
'

run_mlr --opprint --from $indir/abixy put '
  $j = NR;
  do {
    $k = NR;
    $j += 1;
    do {
      $k += 1;
      break;
      $k += 10000;
    } while ($k < 7);
  } while ($j < 4);
'

run_mlr --opprint --from $indir/abixy put '
  $j = NR;
  do {
    $k = NR;
    $j += 1;
    do {
      $k += 1;
      continue;
      $k += 10000;
    } while ($k < 7);
  } while ($j < 4);
'


run_mlr --opprint --from $indir/abixy put '
  $j = NR;
  do {
    $k = NR;
    $j += 1;
    if (NR == 2 || NR == 8) {
      break;
    }
    do {
      $k += 1
    } while ($k < 7);
  } while ($j < 4);
'

run_mlr --opprint --from $indir/abixy put '
  $j = NR;
  do {
    $k = NR;
    $j += 1;
    if (NR == 2 || NR == 8) {
      continue;
    }
    do {
      $k += 1
    } while ($k < 7);
  } while ($j < 4);
'

run_mlr --opprint --from $indir/abixy put '
  $j = NR;
  do {
    $k = NR;
    $j += 1;
    do {
      $k += 1;
      if (NR == 2 || NR == 8) {
        break;
      }
      $k += 10000;
    } while ($k < 7);
  } while ($j < 4);
'

run_mlr --opprint --from $indir/abixy put '
  $j = NR;
  do {
    $k = NR;
    $j += 1;
    do {
      $k += 1;
      if (NR == 2 || NR == 8) {
        continue;
      }
      $k += 10000;
    } while ($k < 7);
  } while ($j < 4);
'

# ----------------------------------------------------------------
announce DSL BREAK/CONTINUE IN SINGLE FOR-SREC

run_mlr --opprint --from $indir/abixy put -q '
  for (k,v in $*) {
      @logging1[NR][k] = v;
      if (k == "x") {
          break;
      }
  }
  end {
    emitp @logging1, "NR", "k";
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  for (k,v in $*) {
      if (k == "x") {
          break;
      }
      @logging2[NR][k] = v;
  }
  end {
    emitp @logging2, "NR", "k";
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  for (k,v in $*) {
      @logging3[NR][k] = v;
      if (k == "x") {
          continue;
      }
  }
  end {
    emitp @logging3, "NR", "k";
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  for (k,v in $*) {
      if (k == "x") {
          continue;
      }
      @logging4[NR][k] = v;
  }
  end {
    emitp @logging4, "NR", "k"
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  for (k,v in $*) {
      @logging1[NR][k] = v;
      if (k == "x") {
          break;
      }
  }

  for (k,v in $*) {
      if (k == "x") {
          break;
      }
      @logging2[NR][k] = v;
  }

  for (k,v in $*) {
      @logging3[NR][k] = v;
      if (k == "x") {
          continue;
      }
  }

  for (k,v in $*) {
      if (k == "x") {
          continue;
      }
      @logging4[NR][k] = v;
  }

  end {
    emitp @logging1, "NR", "k";
    emitp @logging2, "NR", "k";
    emitp @logging3, "NR", "k";
    emitp @logging4, "NR", "k"
  }
'

# ----------------------------------------------------------------
announce DSL BREAK/CONTINUE IN NESTED FOR-SREC

run_mlr --opprint --from $indir/abixy put -q '
  for (k1, v1 in $*) {
    @output1[NR][k1] = "before";
    break;
    @output1[NR][k1] = v1;
    for (k2, v2 in $*) {
      @output2[NR][k1."_".k2] = "before";
      @output2[NR][k1."_".k2] = v2;
    }
  }
  end {
    emit @output1, "NR", "name";
    emit @output2, "NR", "names";
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  for (k1, v1 in $*) {
    @output1[NR][k1] = "before";
    continue;
    @output1[NR][k1] = v1;
    for (k2, v2 in $*) {
      @output2[NR][k1."_".k2] = "before";
      @output2[NR][k1."_".k2] = v2;
    }
  }
  end {
    emit @output1, "NR", "name";
    emit @output2, "NR", "names";
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  for (k1, v1 in $*) {
    @output1[NR][k1] = "before";
    @output1[NR][k1] = v1;
    for (k2, v2 in $*) {
      @output2[NR][k1."_".k2] = "before";
      break;
      @output2[NR][k1."_".k2] = v2;
    }
  }
  end {
    emit @output1, "NR", "name";
    emit @output2, "NR", "names";
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  for (k1, v1 in $*) {
    @output1[NR][k1] = "before";
    @output1[NR][k1] = v1;
    for (k2, v2 in $*) {
      @output2[NR][k1."_".k2] = "before";
      continue;
      @output2[NR][k1."_".k2] = v2;
    }
  }
  end {
    emit @output1, "NR", "name";
    emit @output2, "NR", "names";
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  for (k1, v1 in $*) {
    @output1[NR][k1] = "before";
    break;
    @output1[NR][k1] = v1;
    for (k2, v2 in $*) {
      @output2[NR][k1."_".k2] = "before";
      break;
      @output2[NR][k1."_".k2] = v2;
    }
  }
  end {
    emit @output1, "NR", "name";
    emit @output2, "NR", "names";
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  for (k1, v1 in $*) {
    @output1[NR][k1] = "before";
    continue;
    @output1[NR][k1] = v1;
    for (k2, v2 in $*) {
      @output2[NR][k1."_".k2] = "before";
      break;
      @output2[NR][k1."_".k2] = v2;
    }
  }
  end {
    emit @output1, "NR", "name";
    emit @output2, "NR", "names";
  }
'


run_mlr --opprint --from $indir/abixy put -q '
  for (k1, v1 in $*) {
    @output1[NR][k1] = "before";
    break;
    @output1[NR][k1] = v1;
    for (k2, v2 in $*) {
      @output2[NR][k1."_".k2] = "before";
      continue;
      @output2[NR][k1."_".k2] = v2;
    }
  }
  end {
    emit @output1, "NR", "name";
    emit @output2, "NR", "names";
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  for (k1, v1 in $*) {
    @output1[NR][k1] = "before";
    continue;
    @output1[NR][k1] = v1;
    for (k2, v2 in $*) {
      @output2[NR][k1."_".k2] = "before";
      continue;
      @output2[NR][k1."_".k2] = v2;
    }
  }
  end {
    emit @output1, "NR", "name";
    emit @output2, "NR", "names";
  }
'


run_mlr --opprint --from $indir/abixy put -q '
  for (k1, v1 in $*) {
    @output1[NR][k1] = "before";
    if (k1 == "b") {
        break
    }
    @output1[NR][k1] = v1;
    for (k2, v2 in $*) {
      @output2[NR][k1."_".k2] = "before";
      @output2[NR][k1."_".k2] = v2;
    }
  }
  end {
    emit @output1, "NR", "name";
    emit @output2, "NR", "names";
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  for (k1, v1 in $*) {
    @output1[NR][k1] = "before";
    if (k1 == "b") {
        continue
    }
    @output1[NR][k1] = v1;
    for (k2, v2 in $*) {
      @output2[NR][k1."_".k2] = "before";
      @output2[NR][k1."_".k2] = v2;
    }
  }
  end {
    emit @output1, "NR", "name";
    emit @output2, "NR", "names";
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  for (k1, v1 in $*) {
    @output1[NR][k1] = "before";
    @output1[NR][k1] = v1;
    for (k2, v2 in $*) {
      @output2[NR][k1."_".k2] = "before";
      if (k2 == "a") {
          break
      }
      @output2[NR][k1."_".k2] = v2;
    }
  }
  end {
    emit @output1, "NR", "name";
    emit @output2, "NR", "names";
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  for (k1, v1 in $*) {
    @output1[NR][k1] = "before";
    @output1[NR][k1] = v1;
    for (k2, v2 in $*) {
      @output2[NR][k1."_".k2] = "before";
      if (k2 == "b") {
          continue
      }
      @output2[NR][k1."_".k2] = v2;
    }
  }
  end {
    emit @output1, "NR", "name";
    emit @output2, "NR", "names";
  }
'


run_mlr --opprint --from $indir/abixy put -q '
  for (k1, v1 in $*) {
    @output1[NR][k1] = "before";
    if (k1 == "b") {
        break
    }
    @output1[NR][k1] = v1;
    for (k2, v2 in $*) {
      @output2[NR][k1."_".k2] = "before";
      if (k2 == "a") {
          break
      }
      @output2[NR][k1."_".k2] = v2;
    }
  }
  end {
    emit @output1, "NR", "name";
    emit @output2, "NR", "names";
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  for (k1, v1 in $*) {
    @output1[NR][k1] = "before";
    if (k1 == "b") {
        continue
    }
    @output1[NR][k1] = v1;
    for (k2, v2 in $*) {
      @output2[NR][k1."_".k2] = "before";
      if (k2 == "a") {
          break
      }
      @output2[NR][k1."_".k2] = v2;
    }
  }
  end {
    emit @output1, "NR", "name";
    emit @output2, "NR", "names";
  }
'


run_mlr --opprint --from $indir/abixy put -q '
  for (k1, v1 in $*) {
    @output1[NR][k1] = "before";
    if (k1 == "b") {
        break
    }
    @output1[NR][k1] = v1;
    for (k2, v2 in $*) {
      @output2[NR][k1."_".k2] = "before";
      if (k2 == "a") {
          continue
      }
      @output2[NR][k1."_".k2] = v2;
    }
  }
  end {
    emit @output1, "NR", "name";
    emit @output2, "NR", "names";
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  for (k1, v1 in $*) {
    @output1[NR][k1] = "before";
    if (k1 == "b") {
        continue
    }
    @output1[NR][k1] = v1;
    for (k2, v2 in $*) {
      @output2[NR][k1."_".k2] = "before";
      if (k2 == "a") {
          continue
      }
      @output2[NR][k1."_".k2] = v2;
    }
  }
  end {
    emit @output1, "NR", "name";
    emit @output2, "NR", "names";
  }
'
# ----------------------------------------------------------------

announce DSL BREAK/CONTINUE IN SINGLE FOR-OOSVAR

mention single-key tests, direct break/continue

run_mlr --opprint --from $indir/abixy put -q '
  @logging[NR] = $*;
  end {
    for (k1, v in @logging[2]) {
        break;
        @output[k1] = v;
    }
    emit @output, "NR", "name"
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  @logging[NR] = $*;
  end {
    for (k1, v in @logging[2]) {
        @output[k1] = v;
        break;
        @output[k1] = "ERROR";
    }
    emit @output, "NR", "name"
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  @logging[NR] = $*;
  end {
    for (k1, v in @logging[2]) {
        continue;
        @output[k1] = v
    }
    emit @output, "NR", "name"
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  @logging[NR] = $*;
  end {
    for (k1, v in @logging[2]) {
        @output[k1] = v;
        continue;
        @output[k1] = "ERROR";
    }
    emit @output, "NR", "name"
  }
'

mention single-key tests, indirect break/continue

run_mlr --opprint --from $indir/abixy put -q '
  @logging[NR] = $*;
  end {
    for (k1, v in @logging[2]) {
        if (k1 == "i") {
          break;
        }
        @output[k1] = v;
    }
    emit @output, "NR", "name"
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  @logging[NR] = $*;
  end {
    for (k1, v in @logging[2]) {
        @output[k1] = v;
        if (k1 == "i") {
          break;
        }
    }
    emit @output, "NR", "name"
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  @logging[NR] = $*;
  end {
    for (k1, v in @logging[2]) {
        if (k1 == "i") {
          continue;
        }
        @output[k1] = v
    }
    emit @output, "NR", "name"
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  @logging[NR] = $*;
  end {
    for (k1, v in @logging[2]) {
        @output[k1] = v;
        if (k1 == "i") {
          continue;
        }
        @output[k1] = "reached";
    }
    emit @output, "NR", "name"
  }
'

mention multiple-key tests, direct break/continue

run_mlr --opprint --from $indir/abixy put -q '
  @logging[NR] = $*;
  end {
    for ((k1, k2), v in @logging) {
        break;
        @output[k1][k2] = v;
    }
    emit @output, "NR", "name"
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  @logging[NR] = $*;
  end {
    for ((k1, k2), v in @logging) {
        @output[k1][k2] = v;
        break;
        @output[k1][k2] = "ERROR"
    }
    emit @output, "NR", "name"
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  @logging[NR] = $*;
  end {
    for ((k1, k2), v in @logging) {
        continue;
        @output[k1][k2] = v
    }
    emit @output, "NR", "name"
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  @logging[NR] = $*;
  end {
    for ((k1, k2), v in @logging) {
        @output[k1][k2] = v;
        continue;
        @output[k1][k2] = "ERROR";
    }
    emit @output, "NR", "name"
  }
'

mention multiple-key tests, indirect break/continue

run_mlr --opprint --from $indir/abixy put -q '
  @logging[NR] = $*;
  end {
    for ((k1, k2), v in @logging) {
        if (k1 == 5) {
          break;
        }
        @output[k1][k2] = v;
    }
    emit @output, "NR", "name"
  }
'
run_mlr --opprint --from $indir/abixy put -q '
  @logging[NR] = $*;
  end {
    for ((k1, k2), v in @logging) {
        if (k2 == "i") {
          break;
        }
        @output[k1][k2] = v;
    }
    emit @output, "NR", "name"
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  @logging[NR] = $*;
  end {
    for ((k1, k2), v in @logging) {
        @output[k1][k2] = v;
        if (k1 == 5) {
          break;
        }
    }
    emit @output, "NR", "name"
  }
'
run_mlr --opprint --from $indir/abixy put -q '
  @logging[NR] = $*;
  end {
    for ((k1, k2), v in @logging) {
        @output[k1][k2] = v;
        if (k2 == "i") {
          break;
        }
    }
    emit @output, "NR", "name"
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  @logging[NR] = $*;
  end {
    for ((k1, k2), v in @logging) {
        if (k1 == 5) {
          continue;
        }
        @output[k1][k2] = v
    }
    emit @output, "NR", "name"
  }
'
run_mlr --opprint --from $indir/abixy put -q '
  @logging[NR] = $*;
  end {
    for ((k1, k2), v in @logging) {
        if (k2 == "i") {
          continue;
        }
        @output[k1][k2] = v
    }
    emit @output, "NR", "name"
  }
'

run_mlr --opprint --from $indir/abixy put -q '
  @logging[NR] = $*;
  end {
    for ((k1, k2), v in @logging) {
        @output[k1][k2] = "before";
        if (k1 == 5) {
          continue;
        }
        @output[k1][k2] = v;
    }
    emit @output, "NR", "name"
  }
'
run_mlr --opprint --from $indir/abixy put -q '
  @logging[NR] = $*;
  end {
    for ((k1, k2), v in @logging) {
        @output[k1][k2] = "before";
        if (k2 == "i") {
          continue;
        }
        @output[k1][k2] = v;
    }
    emit @output, "NR", "name"
  }
'

# ----------------------------------------------------------------
announce DSL BREAK/CONTINUE IN NESTED FOR-OOSVAR

run_mlr --opprint --from $indir/abixy put -q '
  @logging[NR] = $*;
  end {
    for ((k1, k2), v in @logging) {
        if (k1 != 2) {
          continue
        }
        for ((k3, k4), v in @logging) {
          if (k3 != 4) {
            continue
          }
          @output[k1][k2][k3][k4] = v;
        }
    }
    emit @output, "NR1", "name1", "NR2", "name2"
  }
'
