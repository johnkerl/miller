#mlrs="mlr5 ~/tmp/miller/mlr ./mlr"
#reps="1"

mlrs="mlr5 ./mlr"
reps="1 2 3"

echo; for mlr in $mlrs; do for k in $reps; do justtime $mlr --csv --from ~/tmp/big.csv check > /dev/null;  done; done
echo; for mlr in $mlrs; do for k in $reps; do justtime $mlr --csv --from ~/tmp/big.csv cat   > /dev/null;  done; done
echo; for mlr in $mlrs; do for k in $reps; do justtime $mlr --csv --from ~/tmp/big.csv head  > /dev/null;  done; done
echo; for mlr in $mlrs; do for k in $reps; do justtime $mlr --csv --from ~/tmp/big.csv tail  > /dev/null;  done; done
echo; for mlr in $mlrs; do for k in $reps; do justtime $mlr --csv --from ~/tmp/big.csv tac   > /dev/null;  done; done
echo; for mlr in $mlrs; do for k in $reps; do justtime $mlr --csv --from ~/tmp/big.csv sort -f shape    > /dev/null; done; done
echo; for mlr in $mlrs; do for k in $reps; do justtime $mlr --csv --from ~/tmp/big.csv sort -n quantity > /dev/null; done; done
