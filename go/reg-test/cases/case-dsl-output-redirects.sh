# ----------------------------------------------------------------
mention print

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; print' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; print > stdout' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; print > stderr' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; print > ENV["outdir"]."/foo.dat"' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt
run_cat $outdir/foo.dat

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; print @x' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; print > stdout, @x' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; print > stderr, @x' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; print > ENV["outdir"]."/foo.dat", @x' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt
run_cat $outdir/foo.dat

# ----------------------------------------------------------------
mention printn

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; printn' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; printn > stdout' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; printn > stderr' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; printn > ENV["outdir"]."/foo.dat"' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt
run_cat $outdir/foo.dat

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; printn @x' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; printn > stdout, @x' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; printn > stderr, @x' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; printn > ENV["outdir"]."/foo.dat", @x' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt
run_cat $outdir/foo.dat

# ----------------------------------------------------------------
mention eprint

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; eprint' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; eprint @x' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

# ----------------------------------------------------------------
mention eprintn

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; eprintn' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; eprintn @x' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

# ----------------------------------------------------------------
mention dump

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; dump' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; dump > stdout' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; dump > stderr' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; dump > ENV["outdir"]."/foo.dat"' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt
run_cat $outdir/foo.dat

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; dump @x' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; dump > stdout, @x' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; dump > stderr, @x' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; dump > ENV["outdir"]."/foo.dat", @x' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt
run_cat $outdir/foo.dat

# ----------------------------------------------------------------
mention edump

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; edump' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; edump @x' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

# ----------------------------------------------------------------
mention tee

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; tee > stdout, $*' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; tee > stderr, $*' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; tee > ENV["outdir"]."/foo.dat", $*' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt
run_cat $outdir/foo.dat

# ----------------------------------------------------------------
mention emitf

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; emitf @x' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; emitf > stdout, @x' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; emitf > stderr, @x' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; emitf > ENV["outdir"]."/foo.dat", @x' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt
run_cat $outdir/foo.dat

# ----------------------------------------------------------------
mention emit

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; emit @x' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; emit > stdout, @x' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; emit > stderr, @x' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; emit > ENV["outdir"]."/foo.dat", @x' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt
run_cat $outdir/foo.dat

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; @y=2; emit (@x, @y)' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; @y=2; emit > stdout, (@x, @y)' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; @y=2; emit > stderr, (@x, @y)' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; @y=2; emit > ENV["outdir"]."/foo.dat", (@x, @y)' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt
run_cat $outdir/foo.dat

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x={"a":1}; emit @x, "a"' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x={"a":1}; emit > stdout, @x, "a"' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x={"a":1}; emit > stderr, @x, "a"' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x={"a":1}; emit > ENV["outdir"]."/foo.dat", @x, "a"' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt
run_cat $outdir/foo.dat

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x={"a":1}; @y={"a":2}; emit (@x, @y), "a"' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x={"a":1}; @y={"a":2}; emit > stdout, (@x, @y), "a"' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x={"a":1}; @y={"a":2}; emit > stderr, (@x, @y), "a"' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x={"a":1}; @y={"a":2}; emit > ENV["outdir"]."/foo.dat", (@x, @y), "a"' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt
run_cat $outdir/foo.dat

# ----------------------------------------------------------------
mention emitp

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; emitp @x' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; emitp > stdout, @x' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; emitp > stderr, @x' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; emitp > ENV["outdir"]."/foo.dat", @x' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt
run_cat $outdir/foo.dat

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; @y=2; emitp (@x, @y)' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; @y=2; emitp > stdout, (@x, @y)' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; @y=2; emitp > stderr, (@x, @y)' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x=1; @y=2; emitp > ENV["outdir"]."/foo.dat", (@x, @y)' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt
run_cat $outdir/foo.dat

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x={"a":1}; emitp @x, "a"' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x={"a":1}; emitp > stdout, @x, "a"' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x={"a":1}; emitp > stderr, @x, "a"' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x={"a":1}; emitp > ENV["outdir"]."/foo.dat", @x, "a"' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt
run_cat $outdir/foo.dat

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x={"a":1}; @y={"a":2}; emitp (@x, @y), "a"' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x={"a":1}; @y={"a":2}; emitp > stdout, (@x, @y), "a"' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x={"a":1}; @y={"a":2}; emitp > stderr, (@x, @y), "a"' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt

rm -f $outdir/redirect-out.txt $outdir/redirect-err.txt $outdir/foo.dat
run_mlr_externally_redirected --from $indir/2.dkvp put '@x={"a":1}; @y={"a":2}; emitp > ENV["outdir"]."/foo.dat", (@x, @y), "a"' \
  1> $outdir/redirect-out.txt 2> $outdir/redirect-err.txt
run_cat $outdir/redirect-out.txt
run_cat $outdir/redirect-err.txt
run_cat $outdir/foo.dat
