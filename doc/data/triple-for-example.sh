mlr --from data/small --opprint put '
  local suma = 0;
  local sumb = 0;
  for (local a = 1, local b = 1; a <= NR; a += 1, b *= 2) {
    suma += a;
    sumb += b;
  }
  $suma = suma;
  $sumb = sumb;
'
