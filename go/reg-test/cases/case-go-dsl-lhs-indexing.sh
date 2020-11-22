run_mlr --ojson --from $indir/2.dkvp put '$abc[FILENAME] = "def"'
run_mlr --ojson --from $indir/2.dkvp put '$abc[NR] = "def"'
run_mlr --ojson --from $indir/2.dkvp put '$abc[FILENAME][NR] = "def"'
run_mlr --ojson --from $indir/2.dkvp put '$abc[NR][FILENAME] = "def"'

run_mlr --ojson --from $indir/2.dkvp put '@abc[FILENAME] = "def"; $ghi = @abc'
run_mlr --ojson --from $indir/2.dkvp put '@abc[NR] = "def"; $ghi = @abc'
run_mlr --ojson --from $indir/2.dkvp put '@abc[FILENAME][NR] = "def"; $ghi = @abc'
run_mlr --ojson --from $indir/2.dkvp put '@abc[NR][FILENAME] = "def"; $ghi = @abc'

run_mlr --from $indir/2.dkvp --ojson put '@a = 3; $new=@a'
run_mlr --from $indir/2.dkvp --ojson put '@a = 3; @a[1]=4; $new=@a'
run_mlr --from $indir/2.dkvp --ojson put '@a = 3; @a[1]=4;@a[1][1]=5; $new=@a'

run_mlr --from $indir/2.dkvp --ojson put '@a = 3; @a["x"]=4; $new=@a'
run_mlr --from $indir/2.dkvp --ojson put '@a = 3; @a["x"]=4;@a["x"]["x"]=5; $new=@a'
