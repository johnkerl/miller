mlr --from data/rect.txt put -q '
  ispresent($outer) {
    unset @r
  }
  for (k, v in $*) {
    @r[k] = v
  }
  ispresent($inner1) {
    emit @r
  }'
