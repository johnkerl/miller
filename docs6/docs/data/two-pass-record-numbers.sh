mlr --opprint --from data/small put -q '
  @records[NR] = $*;
  end {
    for((I,k),v in @records) {
      @records[I]["I"] = I;
      @records[I]["N"] = NR;
      @records[I]["PCT"] = 100*I/NR
    }
    emit @records,"I"
  }
' then reorder -f I,N,PCT
