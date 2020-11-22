# ----------------------------------------------------------------
announce DSL SCIENTIFIC NOTATION IN FIELD VALUES

run_mlr --opprint put '$y=$x+1' $indir/scinot.dkvp

# ----------------------------------------------------------------
announce DSL SCIENTIFIC NOTATION IN EXPRESSION LITERALS

run_mlr --opprint put '$y = 123     + $i' $indir/scinot1.dkvp
run_mlr --opprint put '$y = 123.    + $i' $indir/scinot1.dkvp
run_mlr --opprint put '$y = 123.4   + $i' $indir/scinot1.dkvp
run_mlr --opprint put '$y = .234    + $i' $indir/scinot1.dkvp
run_mlr --opprint put '$y = 1e2     + $i' $indir/scinot1.dkvp
run_mlr --opprint put '$y = 1e-2    + $i' $indir/scinot1.dkvp
run_mlr --opprint put '$y = 1.2e3   + $i' $indir/scinot1.dkvp
run_mlr --opprint put '$y = 1.e3    + $i' $indir/scinot1.dkvp
run_mlr --opprint put '$y = 1.2e-3  + $i' $indir/scinot1.dkvp
run_mlr --opprint put '$y = 1.e-3   + $i' $indir/scinot1.dkvp
run_mlr --opprint put '$y = .2e3    + $i' $indir/scinot1.dkvp
run_mlr --opprint put '$y = .2e-3   + $i' $indir/scinot1.dkvp
run_mlr --opprint put '$y = 1.e-3   + $i' $indir/scinot1.dkvp
run_mlr --opprint put '$y = -123    + $i' $indir/scinot1.dkvp
run_mlr --opprint put '$y = -123.   + $i' $indir/scinot1.dkvp
run_mlr --opprint put '$y = -123.4  + $i' $indir/scinot1.dkvp
run_mlr --opprint put '$y = -.234   + $i' $indir/scinot1.dkvp
run_mlr --opprint put '$y = -1e2    + $i' $indir/scinot1.dkvp
run_mlr --opprint put '$y = -1e-2   + $i' $indir/scinot1.dkvp
run_mlr --opprint put '$y = -1.2e3  + $i' $indir/scinot1.dkvp
run_mlr --opprint put '$y = -1.e3   + $i' $indir/scinot1.dkvp
run_mlr --opprint put '$y = -1.2e-3 + $i' $indir/scinot1.dkvp
run_mlr --opprint put '$y = -1.e-3  + $i' $indir/scinot1.dkvp
run_mlr --opprint put '$y = -.2e3   + $i' $indir/scinot1.dkvp
run_mlr --opprint put '$y = -.2e-3  + $i' $indir/scinot1.dkvp
run_mlr --opprint put '$y = -1.e-3  + $i' $indir/scinot1.dkvp
