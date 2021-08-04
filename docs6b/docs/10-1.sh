grep op=cache log.txt \
  | mlr --idkvp --opprint stats1 -a mean -f hit -g type then sort -f type
