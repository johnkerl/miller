
run_mlr --opprint put '$f=is_absent($x)'                          $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_absent($y)'                          $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_absent($z)'                          $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_absent($nosuch)'                     $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_absent(@nosuch)'                     $indir/nullvals.dkvp
run_mlr --opprint put '@somesuch=1;$f=is_absent(@somesuch)'       $indir/nullvals.dkvp

run_mlr --opprint put '$f=is_bool($x>1)'                          $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_bool($y>1)'                          $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_bool($z>1)'                          $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_bool($nosuch>1)'                     $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_bool(@nosuch>1)'                     $indir/nullvals.dkvp
run_mlr --opprint put '@somesuch=1;$f=is_bool(@somesuch>1)'       $indir/nullvals.dkvp

run_mlr --opprint put '$f=is_boolean($x>1)'                       $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_boolean($y>1)'                       $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_boolean($z>1)'                       $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_boolean($nosuch>1)'                  $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_boolean(@nosuch>1)'                  $indir/nullvals.dkvp
run_mlr --opprint put '@somesuch=1;$f=is_boolean(@somesuch>1)'    $indir/nullvals.dkvp

run_mlr --opprint put '$f=is_empty($x)'                           $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_empty($y)'                           $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_empty($z)'                           $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_empty($nosuch)'                      $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_empty(@nosuch)'                      $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_empty($*)'                           $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_empty({1:2})'                        $indir/nullvals.dkvp
run_mlr --opprint put '@somesuch=1;$f=is_empty(@somesuch)'        $indir/nullvals.dkvp

run_mlr --opprint put '$f=is_empty_map($x)'                       $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_empty_map($y)'                       $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_empty_map($z)'                       $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_empty_map($nosuch)'                  $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_empty_map(@nosuch)'                  $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_empty_map($*)'                       $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_empty_map({1:2})'                    $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_empty_map({})'                       $indir/nullvals.dkvp
run_mlr --opprint put '@somesuch=1;$f=is_empty_map(@somesuch)'    $indir/nullvals.dkvp

run_mlr --opprint put '$f=is_float($x)'                           $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_float($y)'                           $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_float($z)'                           $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_float($nosuch)'                      $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_float(@nosuch)'                      $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_float($*)'                           $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_float({1:2})'                        $indir/nullvals.dkvp
run_mlr --opprint put '@somesuch=1;$f=is_float(@somesuch)'        $indir/nullvals.dkvp

run_mlr --opprint put '$f=is_int($x)'                             $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_int($y)'                             $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_int($z)'                             $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_int($nosuch)'                        $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_int(@nosuch)'                        $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_int($*)'                             $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_int({1:2})'                          $indir/nullvals.dkvp
run_mlr --opprint put '@somesuch=1;$f=is_int(@somesuch)'          $indir/nullvals.dkvp

run_mlr --opprint put '$f=is_map($x)'                             $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_map($y)'                             $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_map($z)'                             $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_map($nosuch)'                        $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_map(@nosuch)'                        $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_map($*)'                             $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_map({1:2})'                          $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_map({})'                             $indir/nullvals.dkvp
run_mlr --opprint put '@somesuch=1;$f=is_map(@somesuch)'          $indir/nullvals.dkvp

run_mlr --opprint put '$f=is_nonempty_map($x)'                    $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_nonempty_map($y)'                    $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_nonempty_map($z)'                    $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_nonempty_map($nosuch)'               $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_nonempty_map(@nosuch)'               $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_nonempty_map($*)'                    $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_nonempty_map({1:2})'                 $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_nonempty_map({})'                    $indir/nullvals.dkvp
run_mlr --opprint put '@somesuch=1;$f=is_nonempty_map(@somesuch)' $indir/nullvals.dkvp

run_mlr --opprint put '$f=is_not_empty($x)'                       $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_not_empty($y)'                       $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_not_empty($z)'                       $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_not_empty($nosuch)'                  $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_not_empty(@nosuch)'                  $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_not_empty($*)'                       $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_not_empty({1:2})'                    $indir/nullvals.dkvp
run_mlr --opprint put '@somesuch=1;$f=is_not_empty(@somesuch)'    $indir/nullvals.dkvp

run_mlr --opprint put '$f=is_not_map($x)'                         $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_not_map($y)'                         $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_not_map($z)'                         $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_not_map($nosuch)'                    $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_not_map(@nosuch)'                    $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_not_map($*)'                         $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_not_map({1:2})'                      $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_not_map({})'                         $indir/nullvals.dkvp
run_mlr --opprint put '@somesuch=1;$f=is_not_map(@somesuch)'      $indir/nullvals.dkvp

run_mlr --opprint put '$f=is_not_null($x)'                        $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_not_null($y)'                        $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_not_null($z)'                        $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_not_null($nosuch)'                   $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_not_null(@nosuch)'                   $indir/nullvals.dkvp
run_mlr --opprint put '@somesuch=1;$f=is_not_null(@somesuch)'     $indir/nullvals.dkvp

run_mlr --opprint put '$f=is_null($x)'                            $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_null($y)'                            $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_null($z)'                            $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_null($nosuch)'                       $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_null(@nosuch)'                       $indir/nullvals.dkvp
run_mlr --opprint put '@somesuch=1;$f=is_null(@somesuch)'         $indir/nullvals.dkvp

run_mlr --opprint put '$f=is_numeric($x)'                         $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_numeric($y)'                         $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_numeric($z)'                         $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_numeric($nosuch)'                    $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_numeric(@nosuch)'                    $indir/nullvals.dkvp
run_mlr --opprint put '@somesuch=1;$f=is_numeric(@somesuch)'      $indir/nullvals.dkvp

run_mlr --opprint put '$f=is_present($x)'                         $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_present($y)'                         $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_present($z)'                         $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_present($nosuch)'                    $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_present(@nosuch)'                    $indir/nullvals.dkvp
run_mlr --opprint put '@somesuch=1;$f=is_present(@somesuch)'      $indir/nullvals.dkvp

run_mlr --opprint put '$f=is_string($x)'                          $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_string($y)'                          $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_string($z)'                          $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_string($nosuch)'                     $indir/nullvals.dkvp
run_mlr --opprint put '$f=is_string(@nosuch)'                     $indir/nullvals.dkvp
run_mlr --opprint put '@somesuch=1;$f=is_string(@somesuch)'       $indir/nullvals.dkvp
