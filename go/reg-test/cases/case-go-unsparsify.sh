run_mlr --opprint --from $indir/needs-unsparsify.dkvp unsparsify
run_mlr --opprint --from $indir/needs-unsparsify.dkvp unsparsify --fill-with X
run_mlr --opprint --from $indir/abixy-het unsparsify
run_mlr --opprint --from $indir/abixy-het unsparsify -f a,b,i,x,y
run_mlr --opprint --from $indir/abixy-het unsparsify -f a,b,i,x,y then regularize
run_mlr --opprint --from $indir/abixy-het unsparsify -f aaa,bbb,iii,xxx,yyy
run_mlr --opprint --from $indir/abixy-het unsparsify -f aaa,bbb,iii,xxx,yyy then regularize
run_mlr --opprint --from $indir/abixy-het unsparsify -f a,b,i,x,y,aaa,bbb,iii,xxx,yyy
run_mlr --opprint --from $indir/abixy-het unsparsify -f a,b,i,x,y,aaa,bbb,iii,xxx,yyy then regularize
