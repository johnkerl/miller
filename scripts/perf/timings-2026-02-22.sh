datafile=timings-2026-02-22.dat

mlr --d2p --from $datafile \
  grep cat then reshape -s desc,seconds \
  | sed '1s/^/#/' \
  | pgr -cat -flabels -lul -lp -ms 5 -o cats.png &

mlr --d2p --from $datafile \
  grep chain then reshape -s desc,seconds \
  | sed '1s/^/#/' \
  | pgr -cat -flabels -lul -lp -ms 5 -o chains.png &

mlr --d2p --from $datafile \
  grep -v cat then grep -v chain then reshape -s desc,seconds \
  | sed '1s/^/#/'  \
  | pgr -cat -flabels -lul -lp -ms 5 -o verbs.png &
