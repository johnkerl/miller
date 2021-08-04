mlr --from data/small put '
  print "NR = ".NR;
  for (key in $*) {
    value = $[key];
    print "  key:" . key . "  value:".value;
  }

'
