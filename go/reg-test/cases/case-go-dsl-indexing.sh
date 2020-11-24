run_mlr --from $indir/s.dkvp --idkvp --opprint put '$z = {"a":$a,"b":$b,"i":$i,"x":$x,"y":$y}["b"]'

run_mlr --from $indir/s.dkvp --from $indir/t.dkvp --ojson put '$z=[1,2,[NR,[FILENAME,5],$x*$y]]'

run_mlr --from $indir/s.dkvp --idkvp --ojson put '$z = $*["a"]'
run_mlr --from $indir/s.dkvp --idkvp --ojson put '$z = $*'

run_mlr --from $indir/s.dkvp --idkvp --ojson put '$* = {"s": 7, "t": 8}'
run_mlr --from $indir/s.dkvp --idkvp --ojson put '$*["st"] = 78'
run_mlr --from $indir/s.dkvp --idkvp --ojson put '$*["a"] = 78'
run_mlr --from $indir/s.dkvp --idkvp --ojson put '$*["a"] = {}'
run_mlr --from $indir/s.dkvp --idkvp --opprint put -v '$new = $["a"]'
run_mlr --from $indir/s.dkvp --idkvp --opprint put -v '$["new"] = $a'
run_mlr --from $indir/s.dkvp --idkvp --opprint put -v '${new} = $a . $b'
run_mlr --from $indir/s.dkvp --idkvp --opprint put -v '$new = ${a} . ${b}'

run_mlr --from $indir/s.dkvp --idkvp --opprint put '@tmp = $a . $b; $ab = @tmp'
run_mlr --ojson --from $indir/s.dkvp put '@curi=$i; $curi = @curi; $lagi=@lagi; @lagi=$i'
run_mlr --from $indir/s.dkvp --ojson put '$z["abc"]["def"]["ghi"]=NR'

run_mlr --json put '$a=$a[2]["b"][1]' $indir/nested.json
