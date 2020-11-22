# in1, in2, out all same formats
run_mlr --json    join -s -j x -f $indir/multi-format-join-a.json $indir/multi-format-join-b.json
run_mlr --dkvp    join -s -j x -f $indir/multi-format-join-a.dkvp $indir/multi-format-join-b.dkvp
run_mlr --csvlite join -s -j x -f $indir/multi-format-join-a.csv  $indir/multi-format-join-b.csv

run_mlr --json    join -j x -f $indir/multi-format-join-a.json $indir/multi-format-join-b.json
run_mlr --dkvp    join -j x -f $indir/multi-format-join-a.dkvp $indir/multi-format-join-b.dkvp
run_mlr --csvlite join -j x -f $indir/multi-format-join-a.csv  $indir/multi-format-join-b.csv

# in2 different format
run_mlr --json    join -s -i csvlite -j x -f $indir/multi-format-join-a.csv $indir/multi-format-join-b.json
run_mlr --dkvp    join -s -i csvlite -j x -f $indir/multi-format-join-a.csv $indir/multi-format-join-b.dkvp
run_mlr --csvlite join -s -i csvlite -j x -f $indir/multi-format-join-a.csv $indir/multi-format-join-b.csv

run_mlr --json    join    -i csvlite -j x -f $indir/multi-format-join-a.csv $indir/multi-format-join-b.json
run_mlr --dkvp    join    -i csvlite -j x -f $indir/multi-format-join-a.csv $indir/multi-format-join-b.dkvp
run_mlr --csvlite join    -i csvlite -j x -f $indir/multi-format-join-a.csv $indir/multi-format-join-b.csv

run_mlr --json    join -s -i dkvp  -j x -f $indir/multi-format-join-a.dkvp $indir/multi-format-join-b.json
run_mlr --dkvp    join -s -i dkvp  -j x -f $indir/multi-format-join-a.dkvp $indir/multi-format-join-b.dkvp
run_mlr --csvlite join -s -i dkvp  -j x -f $indir/multi-format-join-a.dkvp $indir/multi-format-join-b.csv

run_mlr --json    join    -i dkvp  -j x -f $indir/multi-format-join-a.dkvp $indir/multi-format-join-b.json
run_mlr --dkvp    join    -i dkvp  -j x -f $indir/multi-format-join-a.dkvp $indir/multi-format-join-b.dkvp
run_mlr --csvlite join    -i dkvp  -j x -f $indir/multi-format-join-a.dkvp $indir/multi-format-join-b.csv

run_mlr --json    join -s -i json  -j x -f $indir/multi-format-join-a.json $indir/multi-format-join-b.json
run_mlr --dkvp    join -s -i json  -j x -f $indir/multi-format-join-a.json $indir/multi-format-join-b.dkvp
run_mlr --csvlite join -s -i json  -j x -f $indir/multi-format-join-a.json $indir/multi-format-join-b.csv

run_mlr --json    join    -i json  -j x -f $indir/multi-format-join-a.json $indir/multi-format-join-b.json
run_mlr --dkvp    join    -i json  -j x -f $indir/multi-format-join-a.json $indir/multi-format-join-b.dkvp
run_mlr --csvlite join    -i json  -j x -f $indir/multi-format-join-a.json $indir/multi-format-join-b.csv

# vary all three
run_mlr --ijson    --ojson join -i csvlite -j x -f $indir/multi-format-join-a.csv $indir/multi-format-join-b.json
run_mlr --idkvp    --ojson join -i csvlite -j x -f $indir/multi-format-join-a.csv $indir/multi-format-join-b.dkvp
run_mlr --icsvlite --ojson join -i csvlite -j x -f $indir/multi-format-join-a.csv $indir/multi-format-join-b.csv

run_mlr --ijson    --ojson join -i dkvp     -j x -f $indir/multi-format-join-a.dkvp $indir/multi-format-join-b.json
run_mlr --idkvp    --ojson join -i dkvp     -j x -f $indir/multi-format-join-a.dkvp $indir/multi-format-join-b.dkvp
run_mlr --icsvlite --ojson join -i dkvp     -j x -f $indir/multi-format-join-a.dkvp $indir/multi-format-join-b.csv

run_mlr --ijson    --ojson join -i json     -j x -f $indir/multi-format-join-a.json $indir/multi-format-join-b.json
run_mlr --idkvp    --ojson join -i json     -j x -f $indir/multi-format-join-a.json $indir/multi-format-join-b.dkvp
run_mlr --icsvlite --ojson join -i json     -j x -f $indir/multi-format-join-a.json $indir/multi-format-join-b.csv


run_mlr --ijson    --odkvp join -i csvlite -j x -f $indir/multi-format-join-a.csv $indir/multi-format-join-b.json
run_mlr --idkvp    --odkvp join -i csvlite -j x -f $indir/multi-format-join-a.csv $indir/multi-format-join-b.dkvp
run_mlr --icsvlite --odkvp join -i csvlite -j x -f $indir/multi-format-join-a.csv $indir/multi-format-join-b.csv

run_mlr --ijson    --odkvp join -i dkvp     -j x -f $indir/multi-format-join-a.dkvp $indir/multi-format-join-b.json
run_mlr --idkvp    --odkvp join -i dkvp     -j x -f $indir/multi-format-join-a.dkvp $indir/multi-format-join-b.dkvp
run_mlr --icsvlite --odkvp join -i dkvp     -j x -f $indir/multi-format-join-a.dkvp $indir/multi-format-join-b.csv

run_mlr --ijson    --odkvp join -i json     -j x -f $indir/multi-format-join-a.json $indir/multi-format-join-b.json
run_mlr --idkvp    --odkvp join -i json     -j x -f $indir/multi-format-join-a.json $indir/multi-format-join-b.dkvp
run_mlr --icsvlite --odkvp join -i json     -j x -f $indir/multi-format-join-a.json $indir/multi-format-join-b.csv


run_mlr --ijson    --ocsvlite join -i csvlite -j x -f $indir/multi-format-join-a.csv $indir/multi-format-join-b.json
run_mlr --idkvp    --ocsvlite join -i csvlite -j x -f $indir/multi-format-join-a.csv $indir/multi-format-join-b.dkvp
run_mlr --icsvlite --ocsvlite join -i csvlite -j x -f $indir/multi-format-join-a.csv $indir/multi-format-join-b.csv

run_mlr --ijson    --ocsvlite join -i dkvp     -j x -f $indir/multi-format-join-a.dkvp $indir/multi-format-join-b.json
run_mlr --idkvp    --ocsvlite join -i dkvp     -j x -f $indir/multi-format-join-a.dkvp $indir/multi-format-join-b.dkvp
run_mlr --icsvlite --ocsvlite join -i dkvp     -j x -f $indir/multi-format-join-a.dkvp $indir/multi-format-join-b.csv

run_mlr --ijson    --ocsvlite join -i json     -j x -f $indir/multi-format-join-a.json $indir/multi-format-join-b.json
run_mlr --idkvp    --ocsvlite join -i json     -j x -f $indir/multi-format-join-a.json $indir/multi-format-join-b.dkvp
run_mlr --icsvlite --ocsvlite join -i json     -j x -f $indir/multi-format-join-a.json $indir/multi-format-join-b.csv
