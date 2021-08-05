mlr --from data/small --opprint put '
  num suma = 0;
  for (a = 1; a <= NR; a += 1) {
    suma += a;
  }
  $suma = suma;
'
