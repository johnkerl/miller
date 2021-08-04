mlr --from data/ragged.csv --fs comma --nidx put '
  @maxnf = max(@maxnf, NF);
  while(NF < @maxnf) {
    $[NF+1] = "";
  }
'
