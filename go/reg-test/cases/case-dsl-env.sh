

run_mlr -n put -v '$z = ENV["HOME"]'
run_mlr -n put -v '$z = ENV["HOME"][1]'
run_mlr -n put -v '$ENV["FOO"] = "bar"'
run_mlr -n put -v '$ENV["FOO"][2] = "bar"'

export FOO=BAR
run_mlr --from $indir/s.dkvp --opprint head -n 2 then put '$z = ENV["FOO"]'
run_mlr --from $indir/s.dkvp --opprint head -n 2 then put '$z = ENV["FOO"][1]'
run_mlr --from $indir/s.dkvp --opprint head -n 2 then put 'ENV["FOO"] = "bar"'
mlr_expect_fail --from $indir/s.dkvp --opprint head -n 2 then put 'ENV["FOO"][2] = "bar"'
