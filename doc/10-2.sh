mlr --from log.txt --opprint \
  filter 'ispresent($batch_size)' \
  then step -a delta -f time,num_filtered \
  then sec2gmt time

