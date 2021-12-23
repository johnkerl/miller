mlrs="mlr5 ~/tmp/miller/mlr ./mlr"
reps="1"

#mlrs="mlr5 ./mlr"
#reps="1 2 3"

echo; for mlr in $mlrs; do for k in $reps; do justtime $mlr --csv --from ~/tmp/big.csv check |  md5sum;  done; done
echo; for mlr in $mlrs; do for k in $reps; do justtime $mlr --csv --from ~/tmp/big.csv cat   |  md5sum;  done; done
echo; for mlr in $mlrs; do for k in $reps; do justtime $mlr --csv --from ~/tmp/big.csv head  |  md5sum;  done; done
echo; for mlr in $mlrs; do for k in $reps; do justtime $mlr --csv --from ~/tmp/big.csv tail  |  md5sum;  done; done
echo; for mlr in $mlrs; do for k in $reps; do justtime $mlr --csv --from ~/tmp/big.csv tac   |  md5sum;  done; done
echo; for mlr in $mlrs; do for k in $reps; do justtime $mlr --csv --from ~/tmp/big.csv sort -f shape    | md5sum; done; done
echo; for mlr in $mlrs; do for k in $reps; do justtime $mlr --csv --from ~/tmp/big.csv sort -n quantity | md5sum; done; done
