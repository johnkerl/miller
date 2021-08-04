mlr -n put '
  end {
    o = {1:2, 3:{4:5}};
    for (key in o) {
      print "  key:" . key . "  valuetype:" . typeof(o[key]);
    }
  }
'
