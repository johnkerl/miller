mlr --from data/rect.txt put -q '
  is_present($outer) {
    unset @r
  }
  for (k, v in $*) {
    @r[k] = v
  }
  is_present($inner1) {
    emit @r
  }'
