mlrs="mlr"
if [ $# -ge 1 ]; then
    mlrs="$@"
fi

#reps="1"
reps="1 2 3"
#reps="1 2 3 4 5 6 7 8 9 10"

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
