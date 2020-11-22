
# cat reg-test/input/capture.dkvp
# FIELD=ABC123
# FIELD=ABC..123
# FIELD=..ABC..123..
# FIELD=none of the above

run_mlr --opprint put '$FIELD =~ "([A-Z]+)([0-9]+)";         $F1="\1"; $F2="\2"; $F3="\3"' $indir/capture.dkvp
run_mlr --opprint put '$FIELD =~ "([A-Z]+)[^0-9]*([0-9]+)";  $F1="\1"; $F2="\2"; $F3="\3"' $indir/capture.dkvp

run_mlr --opprint put '$FIELD =~ "([A-Z]+)([0-9]+)"'         then put '$F1="\1"; $F2="\2"; $F3="\3"' $indir/capture.dkvp
run_mlr --opprint put '$FIELD =~ "([A-Z]+)[^0-9]*([0-9]+)"'  then put '$F1="\1"; $F2="\2"; $F3="\3"' $indir/capture.dkvp

# cat reg-test/input/capture-lengths.dkvp
# FIELD=
# FIELD=a
# FIELD=ab
# FIELD=abc
# FIELD=abcd
# FIELD=abcde
# FIELD=abcdef
# FIELD=abcdefg
# FIELD=abcdefgh
# FIELD=abcdefghi
# FIELD=abcdefghij
# FIELD=abcdefghijk
# FIELD=abcdefghijkl

run_mlr --opprint put '       $FIELD =~ "....."; $F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"'  $indir/capture-lengths.dkvp
run_mlr --opprint put '       $FIELD =~ "....." {$F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"}' $indir/capture-lengths.dkvp
run_mlr --opprint put 'filter $FIELD =~ "....."; $F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"'  $indir/capture-lengths.dkvp

run_mlr --opprint put '$FIELD =~ "(.)";                            $F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"' $indir/capture-lengths.dkvp
run_mlr --opprint put '$FIELD =~ "(.)(.)";                         $F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"' $indir/capture-lengths.dkvp
run_mlr --opprint put '$FIELD =~ "(.)(.)(.)";                      $F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"' $indir/capture-lengths.dkvp
run_mlr --opprint put '$FIELD =~ "(.)(.)(.)(.)";                   $F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"' $indir/capture-lengths.dkvp
run_mlr --opprint put '$FIELD =~ "(.)(.)(.)(.)(.)";                $F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"' $indir/capture-lengths.dkvp
run_mlr --opprint put '$FIELD =~ "(.)(.)(.)(.)(.)(.)";             $F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"' $indir/capture-lengths.dkvp
run_mlr --opprint put '$FIELD =~ "(.)(.)(.)(.)(.)(.)(.)";          $F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"' $indir/capture-lengths.dkvp
run_mlr --opprint put '$FIELD =~ "(.)(.)(.)(.)(.)(.)(.)(.)";       $F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"' $indir/capture-lengths.dkvp
run_mlr --opprint put '$FIELD =~ "(.)(.)(.)(.)(.)(.)(.)(.)(.)";    $F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"' $indir/capture-lengths.dkvp
run_mlr --opprint put '$FIELD =~ "(.)(.)(.)(.)(.)(.)(.)(.)(.)(.)"; $F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"' $indir/capture-lengths.dkvp

run_mlr --opprint put '$FIELD =~ "(.)"                            {$F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"}' $indir/capture-lengths.dkvp
run_mlr --opprint put '$FIELD =~ "(.)(.)"                         {$F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"}' $indir/capture-lengths.dkvp
run_mlr --opprint put '$FIELD =~ "(.)(.)(.)"                      {$F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"}' $indir/capture-lengths.dkvp
run_mlr --opprint put '$FIELD =~ "(.)(.)(.)(.)"                   {$F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"}' $indir/capture-lengths.dkvp
run_mlr --opprint put '$FIELD =~ "(.)(.)(.)(.)(.)"                {$F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"}' $indir/capture-lengths.dkvp
run_mlr --opprint put '$FIELD =~ "(.)(.)(.)(.)(.)(.)"             {$F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"}' $indir/capture-lengths.dkvp
run_mlr --opprint put '$FIELD =~ "(.)(.)(.)(.)(.)(.)(.)"          {$F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"}' $indir/capture-lengths.dkvp
run_mlr --opprint put '$FIELD =~ "(.)(.)(.)(.)(.)(.)(.)(.)"       {$F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"}' $indir/capture-lengths.dkvp
run_mlr --opprint put '$FIELD =~ "(.)(.)(.)(.)(.)(.)(.)(.)(.)"    {$F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"}' $indir/capture-lengths.dkvp
run_mlr --opprint put '$FIELD =~ "(.)(.)(.)(.)(.)(.)(.)(.)(.)(.)" {$F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"}' $indir/capture-lengths.dkvp

run_mlr --opprint put 'filter $FIELD =~ "(.)";                            $F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"' $indir/capture-lengths.dkvp
run_mlr --opprint put 'filter $FIELD =~ "(.)(.)";                         $F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"' $indir/capture-lengths.dkvp
run_mlr --opprint put 'filter $FIELD =~ "(.)(.)(.)";                      $F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"' $indir/capture-lengths.dkvp
run_mlr --opprint put 'filter $FIELD =~ "(.)(.)(.)(.)";                   $F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"' $indir/capture-lengths.dkvp
run_mlr --opprint put 'filter $FIELD =~ "(.)(.)(.)(.)(.)";                $F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"' $indir/capture-lengths.dkvp
run_mlr --opprint put 'filter $FIELD =~ "(.)(.)(.)(.)(.)(.)";             $F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"' $indir/capture-lengths.dkvp
run_mlr --opprint put 'filter $FIELD =~ "(.)(.)(.)(.)(.)(.)(.)";          $F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"' $indir/capture-lengths.dkvp
run_mlr --opprint put 'filter $FIELD =~ "(.)(.)(.)(.)(.)(.)(.)(.)";       $F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"' $indir/capture-lengths.dkvp
run_mlr --opprint put 'filter $FIELD =~ "(.)(.)(.)(.)(.)(.)(.)(.)(.)";    $F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"' $indir/capture-lengths.dkvp
run_mlr --opprint put 'filter $FIELD =~ "(.)(.)(.)(.)(.)(.)(.)(.)(.)(.)"; $F0="\0";$F1="\1";$F2="\2";$F3="\3";$F4="\4";$F5="\5";$F6="\6";$F7="\7";$F8="\8";$F9="\9"' $indir/capture-lengths.dkvp

echo 'abcdefg' | run_mlr --inidx --odkvp put '$1 =~ "ab(.)d(..)g"  { $c1 = "\1"; $c2 = "\2"}'
echo 'abcdefg' | run_mlr --inidx --odkvp put '$1 =~ "ab(.)?d(..)g" { $c1 = "\1"; $c2 = "\2"}'
echo 'abXdefg' | run_mlr --inidx --odkvp put '$1 =~ "ab(c)?d(..)g" { $c1 = "\1"; $c2 = "\2"}'
echo 'abdefg'  | run_mlr --inidx --odkvp put '$1 =~ "ab(c)?d(..)g" { $c1 = "\1"; $c2 = "\2"}'
