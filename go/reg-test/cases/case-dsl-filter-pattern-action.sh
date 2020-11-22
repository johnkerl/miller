run_mlr --opprint put '         $x > 0.5;  $z = "flag"'  $indir/abixy
run_mlr --opprint put '       !($x > 0.5); $z = "flag"'  $indir/abixy
run_mlr --opprint put 'filter   $x > 0.5;  $z = "flag"'  $indir/abixy
run_mlr --opprint put '         $x > 0.5  {$z = "flag"}' $indir/abixy
run_mlr --opprint put 'filter !($x > 0.5); $z = "flag"'  $indir/abixy
run_mlr --opprint put '       !($x > 0.5) {$z = "flag"}' $indir/abixy
 
# ----------------------------------------------------------------
announce DSL SUB/GSUB/REGEX_EXTRACT

run_mlr --opprint put '$y = sub($x, "e.*l",        "")' $indir/sub.dat
run_mlr --opprint put '$y = sub($x, "e.*l"i,       "")' $indir/sub.dat
run_mlr --opprint put '$y = sub($x, "e.*"."l",     "")' $indir/sub.dat
run_mlr --opprint put '$y = sub($x, "e.*l",        "y123y")' $indir/sub.dat
run_mlr --opprint put '$y = sub($x, "e.*l"i,       "y123y")' $indir/sub.dat
run_mlr --opprint put '$y = sub($x, "e.*"."l",     "y123y")' $indir/sub.dat
run_mlr --opprint put '$y = sub($x, "([hg])e.*l(.)", "y\1y123\2y")' $indir/sub.dat
run_mlr --opprint put '$y = sub($x, "([hg])e.*l.",   "y\1y123\2y")' $indir/sub.dat
run_mlr --opprint put '$y = sub($x, "([hg])e.*l(.)", "y\1y123.y")'  $indir/sub.dat

run_mlr --opprint put '$y = sub($x,  "a",    "aa")'   $indir/gsub.dat
run_mlr --opprint put '$y = gsub($x, "a",    "aa")'   $indir/gsub.dat
run_mlr --opprint put '$y = gsub($x, "A",    "Aa")'   $indir/gsub.dat
run_mlr --opprint put '$y = gsub($x, "a"i,   "Aa")'   $indir/gsub.dat
run_mlr --opprint put '$y = gsub($x, "A"i,   "Aa")'   $indir/gsub.dat
run_mlr --opprint put '$y = gsub($x, "a(.)", "aa\1\1\1")' $indir/gsub.dat

run_mlr --opprint put '$y = sub($x,  "a",    "")'   $indir/gsub.dat
run_mlr --opprint put '$y = gsub($x, "a",    "")'   $indir/gsub.dat
run_mlr --opprint put '$y = gsub($x, "A",    "")'   $indir/gsub.dat
run_mlr --opprint put '$y = gsub($x, "a"i,   "")'   $indir/gsub.dat
run_mlr --opprint put '$y = gsub($x, "A"i,   "")'   $indir/gsub.dat

run_mlr --oxtab cat                       $indir/subtab.dkvp
run_mlr --oxtab put -f $indir/subtab1.mlr $indir/subtab.dkvp
run_mlr --oxtab put -f $indir/subtab2.mlr $indir/subtab.dkvp
run_mlr --oxtab put -f $indir/subtab3.mlr $indir/subtab.dkvp
run_mlr --oxtab put -f $indir/subtab4.mlr $indir/subtab.dkvp

run_mlr --opprint put '$y = ssub($x, "HE",       "")'           $indir/sub.dat
run_mlr --opprint put '$y = ssub($x, "HE",       "HE")'         $indir/sub.dat
run_mlr --opprint put '$y = ssub($x, "HE",       "12345")'      $indir/sub.dat
run_mlr --opprint put '$y = ssub($x, "LL",       "1")'          $indir/sub.dat
run_mlr --opprint put '$y = ssub($x, "LL",       "12")'         $indir/sub.dat
run_mlr --opprint put '$y = ssub($x, "LL",       "12345")'      $indir/sub.dat
run_mlr --opprint put '$y = ssub($x, "LLO",      "")'           $indir/sub.dat
run_mlr --opprint put '$y = ssub($x, "LLO",      "12")'         $indir/sub.dat
run_mlr --opprint put '$y = ssub($x, "LLO",      "123")'        $indir/sub.dat
run_mlr --opprint put '$y = ssub($x, "LLO",      "123456")'     $indir/sub.dat
run_mlr --opprint put '$y = ssub($x, "HELLO",    "")'           $indir/sub.dat
run_mlr --opprint put '$y = ssub($x, "HELLO",    "1234")'       $indir/sub.dat
run_mlr --opprint put '$y = ssub($x, "HELLO",    "12345")'      $indir/sub.dat
run_mlr --opprint put '$y = ssub($x, "HELLO",    "1234678")'    $indir/sub.dat
run_mlr --opprint put '$y = ssub($x, "nonesuch", "")'           $indir/sub.dat
run_mlr --opprint put '$y = ssub($x, "nonesuch", "1234")'       $indir/sub.dat
run_mlr --opprint put '$y = ssub($x, "nonesuch", "1234567890")' $indir/sub.dat

run_mlr --oxtab put '$y = regextract($x, "[A-Z]+")' $indir/sub.dat
run_mlr --oxtab put '$y = regextract($x, "[A-Z]*")' $indir/sub.dat
run_mlr --oxtab put '$y = regextract($x, "[a-z]+")' $indir/sub.dat
run_mlr --oxtab put '$y = regextract($x, "[a-z]*")' $indir/sub.dat
run_mlr --oxtab put '$y = regextract($x, "[0-9]+")' $indir/sub.dat
run_mlr --oxtab put '$y = regextract($x, "[0-9]*")' $indir/sub.dat

run_mlr --oxtab put '$y = regextract($x, "[ef]+")' $indir/sub.dat
run_mlr --oxtab put '$y = regextract($x, "[ef]*")' $indir/sub.dat
run_mlr --oxtab put '$y = regextract($x, "[hi]+")' $indir/sub.dat
run_mlr --oxtab put '$y = regextract($x, "[hi]*")' $indir/sub.dat
run_mlr --oxtab put '$y = regextract($x, "[op]+")' $indir/sub.dat
run_mlr --oxtab put '$y = regextract($x, "[op]*")' $indir/sub.dat

run_mlr --oxtab put '$y = regextract_or_else($x, "[A-Z]+", "DEFAULT")' $indir/sub.dat
run_mlr --oxtab put '$y = regextract_or_else($x, "[A-Z]*", "DEFAULT")' $indir/sub.dat
run_mlr --oxtab put '$y = regextract_or_else($x, "[a-z]+", "DEFAULT")' $indir/sub.dat
run_mlr --oxtab put '$y = regextract_or_else($x, "[a-z]*", "DEFAULT")' $indir/sub.dat
run_mlr --oxtab put '$y = regextract_or_else($x, "[0-9]+", "DEFAULT")' $indir/sub.dat
run_mlr --oxtab put '$y = regextract_or_else($x, "[0-9]*", "DEFAULT")' $indir/sub.dat

run_mlr --oxtab put '$y = regextract_or_else($x, "[ef]+", "DEFAULT")' $indir/sub.dat
run_mlr --oxtab put '$y = regextract_or_else($x, "[ef]*", "DEFAULT")' $indir/sub.dat
run_mlr --oxtab put '$y = regextract_or_else($x, "[hi]+", "DEFAULT")' $indir/sub.dat
run_mlr --oxtab put '$y = regextract_or_else($x, "[hi]*", "DEFAULT")' $indir/sub.dat
run_mlr --oxtab put '$y = regextract_or_else($x, "[op]+", "DEFAULT")' $indir/sub.dat
run_mlr --oxtab put '$y = regextract_or_else($x, "[op]*", "DEFAULT")' $indir/sub.dat

echo 'abcdefg' | run_mlr --nidx put '$1 = sub($1, "ab(.)d(..)g",  "ab<<\1>>d<<\2>>g")'
echo 'abcdefg' | run_mlr --nidx put '$1 = sub($1, "ab(c)?d(..)g", "ab<<\1>>d<<\2>>g")'
echo 'abXdefg' | run_mlr --nidx put '$1 = sub($1, "ab(c)?d(..)g", "ab<<\1>>d<<\2>>g")'
echo 'abdefg'  | run_mlr --nidx put '$1 = sub($1, "ab(c)?d(..)g", "ab<<\1>>d<<\2>>g")'
