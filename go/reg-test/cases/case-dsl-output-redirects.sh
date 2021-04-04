# ----------------------------------------------------------------
mention print

run_mlr --from $indir/2.dkvp put '@x=1; print'

run_mlr --from $indir/2.dkvp put '@x=1; print > stdout'

run_mlr --from $indir/2.dkvp put '@x=1; print > stderr'

run_mlr --from $indir/2.dkvp put '@x=1; print > ENV["outdir"]."/foo.dat"'

run_mlr --from $indir/2.dkvp put '@x=1; print @x'

run_mlr --from $indir/2.dkvp put '@x=1; print > stdout, @x'

run_mlr --from $indir/2.dkvp put '@x=1; print > stderr, @x'

run_mlr --from $indir/2.dkvp put '@x=1; print > ENV["outdir"]."/foo.dat", @x'

# ----------------------------------------------------------------
mention printn

run_mlr --from $indir/2.dkvp put '@x=1; printn'

run_mlr --from $indir/2.dkvp put '@x=1; printn > stdout'

run_mlr --from $indir/2.dkvp put '@x=1; printn > stderr'

run_mlr --from $indir/2.dkvp put '@x=1; printn > ENV["outdir"]."/foo.dat"'

run_mlr --from $indir/2.dkvp put '@x=1; printn @x'

run_mlr --from $indir/2.dkvp put '@x=1; printn > stdout, @x'

run_mlr --from $indir/2.dkvp put '@x=1; printn > stderr, @x'

run_mlr --from $indir/2.dkvp put '@x=1; printn > ENV["outdir"]."/foo.dat", @x'

# ----------------------------------------------------------------
mention eprint

run_mlr --from $indir/2.dkvp put '@x=1; eprint'

run_mlr --from $indir/2.dkvp put '@x=1; eprint @x'

# ----------------------------------------------------------------
mention eprintn

run_mlr --from $indir/2.dkvp put '@x=1; eprintn'

run_mlr --from $indir/2.dkvp put '@x=1; eprintn @x'

# ----------------------------------------------------------------
mention dump

run_mlr --from $indir/2.dkvp put '@x=1; dump'

run_mlr --from $indir/2.dkvp put '@x=1; dump > stdout'

run_mlr --from $indir/2.dkvp put '@x=1; dump > stderr'

run_mlr --from $indir/2.dkvp put '@x=1; dump > ENV["outdir"]."/foo.dat"'

run_mlr --from $indir/2.dkvp put '@x=1; dump @x'

run_mlr --from $indir/2.dkvp put '@x=1; dump > stdout, @x'

run_mlr --from $indir/2.dkvp put '@x=1; dump > stderr, @x'

run_mlr --from $indir/2.dkvp put '@x=1; dump > ENV["outdir"]."/foo.dat", @x'

# ----------------------------------------------------------------
mention edump

run_mlr --from $indir/2.dkvp put '@x=1; edump'

run_mlr --from $indir/2.dkvp put '@x=1; edump @x'

# ----------------------------------------------------------------
mention tee

# Do either 'put -q' (no record-stream output) or use --opprint (record-stream
# output is all at end of stream) since '> stdout' redirection decouples
# record-stream output from print output, resulting in non-deterministic
# output, which makes regtests fail.
run_mlr --from $indir/2.dkvp put -q '@x=1; tee > stdout, $*'

run_mlr --from $indir/2.dkvp put '@x=1; tee > stderr, $*'

run_mlr --from $indir/2.dkvp put '@x=1; tee > ENV["outdir"]."/foo.dat", $*'

# ----------------------------------------------------------------
mention emitf

run_mlr --from $indir/2.dkvp put '@x=1; emitf @x'

# Do either 'put -q' (no record-stream output) or use --opprint (record-stream
# output is all at end of stream) since '> stdout' redirection decouples
# record-stream output from print output, resulting in non-deterministic
# output, which makes regtests fail.
run_mlr --from $indir/2.dkvp put -q '@x=1; emitf > stdout, @x'

run_mlr --from $indir/2.dkvp put '@x=1; emitf > stderr, @x'

run_mlr --from $indir/2.dkvp put '@x=1; emitf > ENV["outdir"]."/foo.dat", @x'

# ----------------------------------------------------------------
mention emit

run_mlr --from $indir/2.dkvp put '@x=1; emit @x'

# Do either 'put -q' (no record-stream output) or use --opprint (record-stream
# output is all at end of stream) since '> stdout' redirection decouples
# record-stream output from print output, resulting in non-deterministic
# output, which makes regtests fail.
run_mlr --from $indir/2.dkvp put -q '@x=1; emit > stdout, @x'

run_mlr --from $indir/2.dkvp put '@x=1; emit > stderr, @x'

run_mlr --from $indir/2.dkvp put '@x=1; emit > ENV["outdir"]."/foo.dat", @x'

run_mlr --from $indir/2.dkvp put '@x=1; @y=2; emit (@x, @y)'

# Do either 'put -q' (no record-stream output) or use --opprint (record-stream
# output is all at end of stream) since '> stdout' redirection decouples
# record-stream output from print output, resulting in non-deterministic
# output, which makes regtests fail.
run_mlr --from $indir/2.dkvp put -q '@x=1; @y=2; emit > stdout, (@x, @y)'

run_mlr --from $indir/2.dkvp put '@x=1; @y=2; emit > stderr, (@x, @y)'

run_mlr --from $indir/2.dkvp put '@x=1; @y=2; emit > ENV["outdir"]."/foo.dat", (@x, @y)'

run_mlr --from $indir/2.dkvp put '@x={"a":1}; emit @x, "a"'

# Do either 'put -q' (no record-stream output) or use --opprint (record-stream
# output is all at end of stream) since '> stdout' redirection decouples
# record-stream output from print output, resulting in non-deterministic
# output, which makes regtests fail.
run_mlr --from $indir/2.dkvp put -q '@x={"a":1}; emit > stdout, @x, "a"'

run_mlr --from $indir/2.dkvp put '@x={"a":1}; emit > stderr, @x, "a"'

run_mlr --from $indir/2.dkvp put '@x={"a":1}; emit > ENV["outdir"]."/foo.dat", @x, "a"'

run_mlr --from $indir/2.dkvp put '@x={"a":1}; @y={"a":2}; emit (@x, @y), "a"'

# Do either 'put -q' (no record-stream output) or use --opprint (record-stream
# output is all at end of stream) since '> stdout' redirection decouples
# record-stream output from print output, resulting in non-deterministic
# output, which makes regtests fail.
run_mlr --from $indir/2.dkvp put -q '@x={"a":1}; @y={"a":2}; emit > stdout, (@x, @y), "a"'

run_mlr --from $indir/2.dkvp put '@x={"a":1}; @y={"a":2}; emit > stderr, (@x, @y), "a"'

run_mlr --from $indir/2.dkvp put '@x={"a":1}; @y={"a":2}; emit > ENV["outdir"]."/foo.dat", (@x, @y), "a"'

# ----------------------------------------------------------------
mention emitp

run_mlr --from $indir/2.dkvp put '@x=1; emitp @x'

# Do either 'put -q' (no record-stream output) or use --opprint (record-stream
# output is all at end of stream) since '> stdout' redirection decouples
# record-stream output from print output, resulting in non-deterministic
# output, which makes regtests fail.
run_mlr --from $indir/2.dkvp put -q '@x=1; emitp > stdout, @x'

run_mlr --from $indir/2.dkvp put '@x=1; emitp > stderr, @x'

run_mlr --from $indir/2.dkvp put '@x=1; emitp > ENV["outdir"]."/foo.dat", @x'

run_mlr --from $indir/2.dkvp put '@x=1; @y=2; emitp (@x, @y)'

# Do either 'put -q' (no record-stream output) or use --opprint (record-stream
# output is all at end of stream) since '> stdout' redirection decouples
# record-stream output from print output, resulting in non-deterministic
# output, which makes regtests fail.
run_mlr --from $indir/2.dkvp put -q '@x=1; @y=2; emitp > stdout, (@x, @y)'

run_mlr --from $indir/2.dkvp put '@x=1; @y=2; emitp > stderr, (@x, @y)'

run_mlr --from $indir/2.dkvp put '@x=1; @y=2; emitp > ENV["outdir"]."/foo.dat", (@x, @y)'

run_mlr --from $indir/2.dkvp put '@x={"a":1}; emitp @x, "a"'

# Do either 'put -q' (no record-stream output) or use --opprint (record-stream
# output is all at end of stream) since '> stdout' redirection decouples
# record-stream output from print output, resulting in non-deterministic
# output, which makes regtests fail.
run_mlr --from $indir/2.dkvp put -q '@x={"a":1}; emitp > stdout, @x, "a"'

run_mlr --from $indir/2.dkvp put '@x={"a":1}; emitp > stderr, @x, "a"'

run_mlr --from $indir/2.dkvp put '@x={"a":1}; emitp > ENV["outdir"]."/foo.dat", @x, "a"'

run_mlr --from $indir/2.dkvp put '@x={"a":1}; @y={"a":2}; emitp (@x, @y), "a"'

# Do either 'put -q' (no record-stream output) or use --opprint (record-stream
# output is all at end of stream) since '> stdout' redirection decouples
# record-stream output from print output, resulting in non-deterministic
# output, which makes regtests fail.
run_mlr --from $indir/2.dkvp put -q '@x={"a":1}; @y={"a":2}; emitp > stdout, (@x, @y), "a"'

run_mlr --from $indir/2.dkvp put '@x={"a":1}; @y={"a":2}; emitp > stderr, (@x, @y), "a"'

run_mlr --from $indir/2.dkvp put '@x={"a":1}; @y={"a":2}; emitp > ENV["outdir"]."/foo.dat", (@x, @y), "a"'

# ----------------------------------------------------------------
# Test separate output formats for mlr main and mlr put.
run_mlr --from $indir/abixy --opprint put --ojson '@x=NR; emit > stdout, @x'
run_mlr --from $indir/abixy --ojson put --opprint '@x=NR; emit > stdout, @x'

