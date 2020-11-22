run_mlr filter -v '$x =~ "bcd"'       $indir/regex.dkvp
run_mlr filter -v '$x =~ "^bcd"'      $indir/regex.dkvp
run_mlr filter -v '$x =~ "^abc"'      $indir/regex.dkvp
run_mlr filter -v '$x =~ "^abc$"'     $indir/regex.dkvp
run_mlr filter -v '$x =~ "^a.*d$"'    $indir/regex.dkvp
run_mlr filter -v '$x =~ "^a.*"."d$"' $indir/regex.dkvp
run_mlr filter -v '$y =~ "\".."'      $indir/regex.dkvp

run_mlr filter -v '$x =~ "bcd"i'       $indir/regex.dkvp
run_mlr filter -v '$x =~ "^bcd"i'      $indir/regex.dkvp
run_mlr filter -v '$x =~ "^abc"i'      $indir/regex.dkvp
run_mlr filter -v '$x =~ "^abc$"i'     $indir/regex.dkvp
run_mlr filter -v '$x =~ "^a.*d$"i'    $indir/regex.dkvp
run_mlr filter -v '$x =~ "^a.*"."d$"i' $indir/regex.dkvp

run_mlr --csv filter '$text =~ "."'    $indir/dot-match.csv
run_mlr --csv filter '$text =~ "\."'   $indir/dot-match.csv
