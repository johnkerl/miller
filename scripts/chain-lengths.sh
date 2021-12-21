mlrs="mlr5 ~/tmp/miller/mlr ./mlr"

echo; for mlr in $mlrs; do
  justtime $mlr --csv --from ~/tmp/big.csv \
    then put -f scripts/chain-1.mlr \
  | md5sum;
done

echo; for mlr in $mlrs; do
  justtime $mlr --csv --from ~/tmp/big.csv \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
  | md5sum;
done

echo; for mlr in $mlrs; do
  justtime $mlr --csv --from ~/tmp/big.csv \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
  | md5sum;
done

echo; for mlr in $mlrs; do
  justtime $mlr --csv --from ~/tmp/big.csv \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
  | md5sum;
done

echo; for mlr in $mlrs; do
  justtime $mlr --csv --from ~/tmp/big.csv \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
  | md5sum;
done

echo; for mlr in $mlrs; do
  justtime $mlr --csv --from ~/tmp/big.csv \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
  | md5sum;
done
