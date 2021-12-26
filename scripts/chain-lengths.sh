mlrs="mlr5 ~/tmp/miller/mlr ./mlr"
#mlrs="mlr5 ./mlr"

#reps="1"
reps="1 2 3"

echo; for mlr in $mlrs; do
  for k in $reps; do
    justtime $mlr --csv --from ~/tmp/big.csv \
      then put -f scripts/chain-1.mlr \
    > /dev/null
  done
done

echo; for mlr in $mlrs; do
  for k in $reps; do
    justtime $mlr --csv --from ~/tmp/big.csv \
      then put -f scripts/chain-1.mlr \
      then put -f scripts/chain-1.mlr \
    > /dev/null
  done
done

echo; for mlr in $mlrs; do
  for k in $reps; do
    justtime $mlr --csv --from ~/tmp/big.csv \
      then put -f scripts/chain-1.mlr \
      then put -f scripts/chain-1.mlr \
      then put -f scripts/chain-1.mlr \
    > /dev/null
  done
done

echo; for mlr in $mlrs; do
  for k in $reps; do
    justtime $mlr --csv --from ~/tmp/big.csv \
      then put -f scripts/chain-1.mlr \
      then put -f scripts/chain-1.mlr \
      then put -f scripts/chain-1.mlr \
      then put -f scripts/chain-1.mlr \
    > /dev/null
  done
done

echo; for mlr in $mlrs; do
  for k in $reps; do
    justtime $mlr --csv --from ~/tmp/big.csv \
      then put -f scripts/chain-1.mlr \
      then put -f scripts/chain-1.mlr \
      then put -f scripts/chain-1.mlr \
      then put -f scripts/chain-1.mlr \
      then put -f scripts/chain-1.mlr \
    > /dev/null
  done
done

echo; for mlr in $mlrs; do
  for k in $reps; do
    justtime $mlr --csv --from ~/tmp/big.csv \
      then put -f scripts/chain-1.mlr \
      then put -f scripts/chain-1.mlr \
      then put -f scripts/chain-1.mlr \
      then put -f scripts/chain-1.mlr \
      then put -f scripts/chain-1.mlr \
      then put -f scripts/chain-1.mlr \
    > /dev/null
  done
done
