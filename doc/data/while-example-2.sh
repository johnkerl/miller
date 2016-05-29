echo x=1,y=2 | mlr put '
  do {
    $[NF+1] = "";
    if (NF == 5) {
      break
    }
  } while (NF < 10);
  $foo = "bar"
'
