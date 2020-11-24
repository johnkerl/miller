run_mlr --icsvlite --odkvp cat $indir/a.csv
run_mlr --icsvlite --odkvp cat $indir/b.csv
run_mlr --icsvlite --odkvp cat $indir/c.csv
run_mlr --icsvlite --odkvp cat $indir/d.csv
run_mlr --icsvlite --odkvp cat $indir/e.csv
run_mlr --icsvlite --odkvp cat $indir/f.csv
run_mlr --icsvlite --odkvp cat $indir/g.csv

run_mlr --icsvlite --odkvp cat $indir/a.csv $indir/a.csv
run_mlr --icsvlite --odkvp cat $indir/b.csv $indir/b.csv
run_mlr --icsvlite --odkvp cat $indir/c.csv $indir/c.csv
run_mlr --icsvlite --odkvp cat $indir/d.csv $indir/d.csv
run_mlr --icsvlite --odkvp cat $indir/e.csv $indir/e.csv
run_mlr --icsvlite --odkvp cat $indir/f.csv $indir/f.csv
run_mlr --icsvlite --odkvp cat $indir/g.csv $indir/g.csv

run_mlr --icsvlite --odkvp cat $indir/a.csv $indir/b.csv
run_mlr --icsvlite --odkvp cat $indir/b.csv $indir/c.csv
run_mlr --icsvlite --odkvp cat $indir/c.csv $indir/d.csv
run_mlr --icsvlite --odkvp cat $indir/d.csv $indir/e.csv
run_mlr --icsvlite --odkvp cat $indir/e.csv $indir/f.csv
run_mlr --icsvlite --odkvp cat $indir/f.csv $indir/g.csv

run_mlr --icsvlite --odkvp cat $indir/a.csv $indir/b.csv \
  $indir/c.csv $indir/d.csv $indir/e.csv $indir/f.csv $indir/g.csv

run_mlr --icsvlite --odkvp tac $indir/het.csv

run_mlr --headerless-csv-output --csvlite tac $indir/a.csv
run_mlr --headerless-csv-output --csvlite tac $indir/c.csv
run_mlr --headerless-csv-output --csvlite tac $indir/a.csv $indir/c.csv
run_mlr --headerless-csv-output --csvlite tac $indir/het.csv
run_mlr --headerless-csv-output --csvlite group-like $indir/het.csv
