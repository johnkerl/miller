
run_mlr         --opprint put '$f=asserting_absent($nosuch)'                     $indir/nullvals.dkvp
run_mlr         --opprint put '$f=asserting_absent(@nosuch)'                     $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_absent($x)'                          $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_absent($y)'                          $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_absent($z)'                          $indir/nullvals.dkvp
mlr_expect_fail --opprint put '@somesuch=1;$f=asserting_absent(@somesuch)'       $indir/nullvals.dkvp
mlr_expect_fail --opprint put 'foo=asserting_absent($*)'                         $indir/nullvals.dkvp
mlr_expect_fail --opprint put 'foo=asserting_absent({1:2})'                      $indir/nullvals.dkvp

run_mlr         --opprint put '$f=asserting_empty($z)'                           $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_empty($x)'                           $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_empty($y)'                           $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_empty($nosuch)'                      $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_empty(@nosuch)'                      $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_empty({1:2})'                        $indir/nullvals.dkvp
mlr_expect_fail --opprint put '@somesuch=1;$f=asserting_empty(@somesuch)'        $indir/nullvals.dkvp
mlr_expect_fail --opprint put 'foo=asserting_empty($*)'                          $indir/nullvals.dkvp
mlr_expect_fail --opprint put 'foo=asserting_empty({1:2})'                       $indir/nullvals.dkvp

run_mlr         --opprint put '$f=asserting_empty_map({})'                       $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_empty_map($*)'                       $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_empty_map($x)'                       $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_empty_map($y)'                       $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_empty_map($z)'                       $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_empty_map($nosuch)'                  $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_empty_map(@nosuch)'                  $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_empty_map({1:2})'                    $indir/nullvals.dkvp
mlr_expect_fail --opprint put '@somesuch=1;$f=asserting_empty_map(@somesuch)'    $indir/nullvals.dkvp

run_mlr         --opprint put '$f=asserting_map($*)'                             $indir/nullvals.dkvp
run_mlr         --opprint put '$f=asserting_map({1:2})'                          $indir/nullvals.dkvp
run_mlr         --opprint put '$f=asserting_map({})'                             $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_map($x)'                             $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_map($y)'                             $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_map($z)'                             $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_map($nosuch)'                        $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_map(@nosuch)'                        $indir/nullvals.dkvp
mlr_expect_fail --opprint put '@somesuch=1;$f=asserting_map(@somesuch)'          $indir/nullvals.dkvp

run_mlr         --opprint put '$f=asserting_nonempty_map($*)'                    $indir/nullvals.dkvp
run_mlr         --opprint put '$f=asserting_nonempty_map({1:2})'                 $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_nonempty_map($x)'                    $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_nonempty_map($y)'                    $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_nonempty_map($z)'                    $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_nonempty_map($nosuch)'               $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_nonempty_map(@nosuch)'               $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_nonempty_map({})'                    $indir/nullvals.dkvp
mlr_expect_fail --opprint put '@somesuch=1;$f=asserting_nonempty_map(@somesuch)' $indir/nullvals.dkvp

run_mlr         --opprint put '$*=asserting_not_empty($*)'                       $indir/nullvals.dkvp
run_mlr         --opprint put '$f=asserting_not_empty($nosuch)'                  $indir/nullvals.dkvp
run_mlr         --opprint put '$f=asserting_not_empty(@nosuch)'                  $indir/nullvals.dkvp
run_mlr         --opprint put '$f=asserting_not_empty({1:2})'                    $indir/nullvals.dkvp
run_mlr         --opprint put '$nosuch=asserting_not_empty($nosuch)'             $indir/nullvals.dkvp
run_mlr         --opprint put '@somesuch=1;$f=asserting_not_empty(@somesuch)'    $indir/nullvals.dkvp
run_mlr         --opprint put '$*=asserting_not_empty($*)'                       $indir/nullvals.dkvp
run_mlr         --opprint put '$*=asserting_not_empty({1:2})'                    $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_not_empty($x)'                       $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_not_empty($y)'                       $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_not_empty($z)'                       $indir/nullvals.dkvp

run_mlr         --opprint put '$f=asserting_not_map($x)'                         $indir/nullvals.dkvp
run_mlr         --opprint put '$f=asserting_not_map($y)'                         $indir/nullvals.dkvp
run_mlr         --opprint put '$f=asserting_not_map($z)'                         $indir/nullvals.dkvp
run_mlr         --opprint put '$f=asserting_not_map($nosuch)'                    $indir/nullvals.dkvp
run_mlr         --opprint put '$f=asserting_not_map(@nosuch)'                    $indir/nullvals.dkvp
run_mlr         --opprint put '@somesuch=1;$f=asserting_not_map(@somesuch)'      $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_not_map($*)'                         $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_not_map({1:2})'                      $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_not_map({})'                         $indir/nullvals.dkvp

run_mlr         --opprint put '@somesuch=1;$f=asserting_not_null(@somesuch)'     $indir/nullvals.dkvp
run_mlr         --opprint put '$*=asserting_not_null($*)'                        $indir/nullvals.dkvp
run_mlr         --opprint put '$*=asserting_not_null({1:2})'                     $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_not_null($x)'                        $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_not_null($y)'                        $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_not_null($z)'                        $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_not_null($nosuch)'                   $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_not_null(@nosuch)'                   $indir/nullvals.dkvp

run_mlr         --opprint put '$f=asserting_null($z)'                            $indir/nullvals.dkvp
run_mlr         --opprint put '$f=asserting_null($nosuch)'                       $indir/nullvals.dkvp
run_mlr         --opprint put '$f=asserting_null(@nosuch)'                       $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_null($x)'                            $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_null($y)'                            $indir/nullvals.dkvp
mlr_expect_fail --opprint put '@somesuch=1;$f=asserting_null(@somesuch)'         $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$*=asserting_null($*)'                            $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$*=asserting_null({1:2})'                         $indir/nullvals.dkvp

mlr_expect_fail --opprint put '$f=asserting_numeric($x)'                         $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_numeric($y)'                         $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_numeric($z)'                         $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$*=asserting_numeric($*)'                         $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$*=asserting_numeric({1:2})'                      $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_numeric($nosuch)'                    $indir/nullvals.dkvp

run_mlr         --opprint put '$f=asserting_present($x)'                         $indir/nullvals.dkvp
run_mlr         --opprint put '$f=asserting_present($y)'                         $indir/nullvals.dkvp
run_mlr         --opprint put '$f=asserting_present($z)'                         $indir/nullvals.dkvp
run_mlr         --opprint put '@somesuch=1;$f=asserting_present(@somesuch)'      $indir/nullvals.dkvp
run_mlr         --opprint put '$*=asserting_present($*)'                         $indir/nullvals.dkvp
run_mlr         --opprint put '$*=asserting_present({1:2})'                      $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_present($nosuch)'                    $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_present(@nosuch)'                    $indir/nullvals.dkvp

run_mlr         --opprint put '$f=asserting_string($z)'                          $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$*=asserting_string($*)'                          $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$*=asserting_string({1:2})'                       $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_string($x)'                          $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_string($y)'                          $indir/nullvals.dkvp
mlr_expect_fail --opprint put '$f=asserting_string($nosuch)'                     $indir/nullvals.dkvp
