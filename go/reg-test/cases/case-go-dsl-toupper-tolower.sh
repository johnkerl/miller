run_mlr --from $indir/s.dkvp --opprint put '
  $u = toupper($a);
  $l = tolower($u);
  $c = capitalize($l);
'
