mlr --from data/miss-date.csv --icsv \
  cat -n \
  then put '$datestamp = strptime($date, "%Y-%m-%d")' \
  then step -a delta -f datestamp \
| head
