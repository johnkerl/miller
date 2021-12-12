echo; for m in mlr5 "./mlr -S"; do
  justtime $m --csv --from ~/tmp/big.csv \
    check \
  | md5sum
done

echo; for m in mlr5 "./mlr -S"; do
  justtime $m --csv --from ~/tmp/big.csv \
    cat \
  | md5sum;
done

echo; for m in mlr5 "./mlr -S"; do
  justtime $m --csv --from ~/tmp/big.csv \
    head \
  | md5sum;
done

echo; for m in mlr5 "./mlr -S"; do
  justtime $m --csv --from ~/tmp/big.csv \
    tail \
  | md5sum;
done

echo; for m in mlr5 "./mlr -S"; do
  justtime $m --csv --from ~/tmp/big.csv \
    tac \
  | md5sum;
done

echo; for m in mlr5 "./mlr -S"; do
  justtime $m --csv --from ~/tmp/big.csv \
    sort -f shape \
  | md5sum;
done

echo; for m in mlr5 "./mlr -S"; do
  justtime $m --csv --from ~/tmp/big.csv \
    sort -n quantity \
  | md5sum;
done
