
run_mlr --ojson seqgen --start 1 --stop 1 then put '
  $x = arrayify({
    "1": "a",  
    "2": "b",  
    "3": "c",  
  })
'

run_mlr --ojson seqgen --start 1 --stop 1 then put '
  $x = arrayify({
    "0": "a",  
    "1": "b",  
    "2": "c",  
  })
'

run_mlr --ojson seqgen --start 1 --stop 1 then put '
  $x = arrayify({
    "1": "a",  
    "3": "b",  
    "5": "c",  
  })
'

run_mlr --ojson seqgen --start 1 --stop 1 then put '
  $x = arrayify({
    "s": {
      "1": "a",  
      "2": "b",  
      "3": "c",  
    }
  })
'

run_mlr --ojson seqgen --start 1 --stop 1 then put '
  $x = arrayify({
    "1": {
      "1": "a",  
      "2": "b",  
      "3": "c",  
    },
    "2": {
      "1": "d",  
      "2": "e",  
      "3": "f",  
    }
  })
'
