run_mlr --from $indir/abixy put -q '@x={"a":NR}; @y={"a":-NR}; emit (@x, @y), "k"'
run_mlr --from $indir/abixy put -q '@x={"a":NR}; @y={"b":-NR}; emit (@x, @y), "k"'
run_mlr --from $indir/abixy put -q '@x={"b":NR}; @y={"a":-NR}; emit (@x, @y), "k"'
run_mlr --from $indir/abixy put -q '@x={"b":NR}; @y={"b":-NR}; emit (@x, @y), "k"'

run_mlr --from $indir/abixy-het put -q '@x={"a":NR}; @y={"a":-NR}; emit (@x, @y), "k"'
run_mlr --from $indir/abixy-het put -q '@x={"a":NR}; @y={"b":-NR}; emit (@x, @y), "k"'
run_mlr --from $indir/abixy-het put -q '@x={"b":NR}; @y={"a":-NR}; emit (@x, @y), "k"'
run_mlr --from $indir/abixy-het put -q '@x={"b":NR}; @y={"b":-NR}; emit (@x, @y), "k"'
