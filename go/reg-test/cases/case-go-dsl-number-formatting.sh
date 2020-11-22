run_mlr --opprint --from $indir/ten.dkvp put '
  $ha = hexfmt($a);
  $hx = hexfmt($x);
  $hi = hexfmt($i);
  $nhi = hexfmt(-$i);
'

run_mlr --opprint --from $indir/ten.dkvp put '
  $hi = "0x".fmtnum($i, "%04x");
  $ex = fmtnum($x, "%8.3e");
'
