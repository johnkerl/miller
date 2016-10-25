mlr --from data/small --opprint put '
  num suma = 0;
  num sumb = 0;
  for (num a = 1, num b = 1; a <= NR; a += 1, b *= 2) {
    suma += a;
    sumb += b;
  }
  $suma = suma;
  $sumb = sumb;
'
