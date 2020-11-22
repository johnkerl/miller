run_mlr --opprint unsparsify $indir/abixy
run_mlr --opprint unsparsify $indir/abixy-het

run_mlr --opprint unsparsify -f nonesuch $indir/abixy-het
run_mlr --opprint unsparsify -f a,b,i,x,y $indir/abixy-het
run_mlr --opprint unsparsify -f aaa,bbb,xxx,iii,yyy $indir/abixy-het
run_mlr --opprint unsparsify -f a,b,i,x,y,aaa,bbb,xxx,iii,yyy $indir/abixy-het
run_mlr --opprint unsparsify -f a,b,i,x,y,aaa,bbb,xxx,iii,yyy then regularize $indir/abixy-het
