mlr --icsv --opprint \
  join -j color --ul --ur -f data/prevtemp.csv \
  then unsparsify --fill-with 0 \
  then put '$count_delta = $current_count - $previous_count' \
  data/currtemp.csv
