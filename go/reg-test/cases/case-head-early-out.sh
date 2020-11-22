# Test early-out for unkeyed head
run_mlr head -n 2 then put 'end{ print "Final NR is ".NR}' $indir/abixy-wide
run_mlr head -n 2 -g a then put 'end{ print "Final NR is ".NR}' $indir/abixy-wide
run_mlr cat then head -n 2 then put 'end{ print "Final NR is ".NR}' $indir/abixy-wide
run_mlr tac then head -n 2 then put 'end{ print "Final NR is ".NR}' $indir/abixy-wide
run_mlr head -n 2 then put 'end{ print "Final NR is ".NR}' $indir/abixy-wide $indir/abixy-wide $indir/abixy-wide
