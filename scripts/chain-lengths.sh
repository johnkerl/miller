echo; for m in mlr5 "./mlr -S"; do
  justtime $m --csv --from ~/tmp/big.csv \
    then put -f scripts/chain-1.mlr \
  | md5sum;
done

echo; for m in mlr5 "./mlr -S"; do
  justtime $m --csv --from ~/tmp/big.csv \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
  | md5sum;
done

echo; for m in mlr5 "./mlr -S"; do
  justtime $m --csv --from ~/tmp/big.csv \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
  | md5sum;
done

echo; for m in mlr5 "./mlr -S"; do
  justtime $m --csv --from ~/tmp/big.csv \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
  | md5sum;
done

echo; for m in mlr5 "./mlr -S"; do
  justtime $m --csv --from ~/tmp/big.csv \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
  | md5sum;
done

echo; for m in mlr5 "./mlr -S"; do
  justtime $m --csv --from ~/tmp/big.csv \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
    then put -f scripts/chain-1.mlr \
  | md5sum;
done
