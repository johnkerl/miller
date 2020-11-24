run_mlr --pprint reshape -i X,Y,Z   -o item,price $indir/reshape-wide.tbl
run_mlr --pprint reshape -i X,Z     -o item,price $indir/reshape-wide.tbl
run_mlr --pprint reshape -r '[X-Z]' -o item,price $indir/reshape-wide.tbl
run_mlr --pprint reshape -r '[XZ]'  -o item,price $indir/reshape-wide.tbl

run_mlr --pprint reshape -s item,price $indir/reshape-long.tbl

run_mlr --pprint reshape -i X,Y,Z -o item,price then reshape -s item,price $indir/reshape-wide.tbl
run_mlr --pprint reshape -s item,price then reshape -i X,Y,Z -o item,price $indir/reshape-long.tbl

run_mlr reshape -i X,Y,Z   -o item,price $indir/reshape-wide-ragged.dkvp
run_mlr reshape -i X,Z     -o item,price $indir/reshape-wide-ragged.dkvp
run_mlr reshape -r '[X-Z]' -o item,price $indir/reshape-wide-ragged.dkvp
run_mlr reshape -r '[XZ]'  -o item,price $indir/reshape-wide-ragged.dkvp

run_mlr reshape -s item,price $indir/reshape-long-ragged.dkvp

run_mlr --json reshape -i x,y -o item,value $indir/small-non-nested.json
