echo x=1,y=2 | mlr put '
  while (NF < 10) {
    $[NF+1] = ""
  }
  $foo = "bar"
'
