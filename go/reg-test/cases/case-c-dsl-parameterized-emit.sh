run_mlr put -q '@sum[$a] = $x; end{ emitp @sum; }'         $indir/abixy
run_mlr put -q '@sum[$a] = $x; end{ emitp @sum,"a"; }'     $indir/abixy
run_mlr put -q '@sum[$a] = $x; end{ emitp @sum,"a","b"; }' $indir/abixy

run_mlr put -q '@sum[$a][$b] = $x; end{ emitp @sum; }'         $indir/abixy
run_mlr put -q '@sum[$a][$b] = $x; end{ emitp @sum,"a"; }'     $indir/abixy
run_mlr put -q '@sum[$a][$b] = $x; end{ emitp @sum,"a","b"; }' $indir/abixy

run_mlr put -q '@v = $a;        end {emitf @v }' $indir/abixy
run_mlr put -q '@v = $i;        end {emitf @v }' $indir/abixy
run_mlr put -q '@v = $x;        end {emitf @v }' $indir/abixy
run_mlr put -q '@v = $nonesuch; end {emitf @v }' $indir/abixy

run_mlr put -q '@v = $a;        end {emitp @v }' $indir/abixy
run_mlr put -q '@v = $i;        end {emitp @v }' $indir/abixy
run_mlr put -q '@v = $x;        end {emitp @v }' $indir/abixy
run_mlr put -q '@v = $nonesuch; end {emitp @v }' $indir/abixy

run_mlr put -q '@sum += $i;        end {emitf @sum }' $indir/abixy
run_mlr put -q '@sum += $x;        end {emitf @sum }' $indir/abixy
run_mlr put -q '@sum += $nonesuch; end {emitf @sum }' $indir/abixy

run_mlr put -q '@sum += $i;        end {emitp  @sum          }' $indir/abixy
run_mlr put -q '@sum += $x;        end {emitp  @sum          }' $indir/abixy
run_mlr put -q '@sum += $nonesuch; end {emitp  @sum          }' $indir/abixy
run_mlr put -q '@sum += $i;        end {emitp  @sum, "extra" }' $indir/abixy
run_mlr put -q '@sum += $x;        end {emitp  @sum, "extra" }' $indir/abixy
run_mlr put -q '@sum += $nonesuch; end {emitp  @sum, "extra" }' $indir/abixy

run_mlr put -q '@sum[$a] += $i;        end {emitp  @sum               }' $indir/abixy
run_mlr put -q '@sum[$a] += $x;        end {emitp  @sum               }' $indir/abixy
run_mlr put -q '@sum[$a] += $nonesuch; end {emitp  @sum               }' $indir/abixy
run_mlr put -q '@sum[$a] += $i;        end {emitp  @sum, "a"          }' $indir/abixy
run_mlr put -q '@sum[$a] += $x;        end {emitp  @sum, "a"          }' $indir/abixy
run_mlr put -q '@sum[$a] += $nonesuch; end {emitp  @sum, "a"          }' $indir/abixy
run_mlr put -q '@sum[$a] += $i;        end {emitp  @sum, "a", "extra" }' $indir/abixy
run_mlr put -q '@sum[$a] += $x;        end {emitp  @sum, "a", "extra" }' $indir/abixy
run_mlr put -q '@sum[$a] += $nonesuch; end {emitp  @sum, "a", "extra" }' $indir/abixy

run_mlr put -q '@sum[$a][$b] += $i;        end {emitp  @sum                    }' $indir/abixy
run_mlr put -q '@sum[$a][$b] += $x;        end {emitp  @sum                    }' $indir/abixy
run_mlr put -q '@sum[$a][$b] += $nonesuch; end {emitp  @sum                    }' $indir/abixy
run_mlr put -q '@sum[$a][$b] += $i;        end {emitp  @sum, "a"               }' $indir/abixy
run_mlr put -q '@sum[$a][$b] += $x;        end {emitp  @sum, "a"               }' $indir/abixy
run_mlr put -q '@sum[$a][$b] += $nonesuch; end {emitp  @sum, "a"               }' $indir/abixy
run_mlr put -q '@sum[$a][$b] += $i;        end {emitp  @sum, "a", "b"          }' $indir/abixy
run_mlr put -q '@sum[$a][$b] += $x;        end {emitp  @sum, "a", "b"          }' $indir/abixy
run_mlr put -q '@sum[$a][$b] += $nonesuch; end {emitp  @sum, "a", "b"          }' $indir/abixy
run_mlr put -q '@sum[$a][$b] += $i;        end {emitp  @sum, "a", "b", "extra" }' $indir/abixy
run_mlr put -q '@sum[$a][$b] += $x;        end {emitp  @sum, "a", "b", "extra" }' $indir/abixy
run_mlr put -q '@sum[$a][$b] += $nonesuch; end {emitp  @sum, "a", "b", "extra" }' $indir/abixy

run_mlr --oxtab put -q '@sum[$a][$b] += $i; NR == 3 { @x = $x }; NR == 7 { @v = $* }; end {emitp all}' $indir/abixy-het

run_mlr --opprint put -q '@v[NR]=$*; end{emitp @v}'     $indir/abixy
run_mlr --opprint put -q '@v[NR]=$*; end{emitp @v,"X"}' $indir/abixy

run_mlr --opprint put -q '@v[NR]=$*; end{emitp @v[1]}'         $indir/abixy
run_mlr --opprint put -q '@v[NR]=$*; end{emitp @v[1],"X"}'     $indir/abixy
run_mlr --opprint put -q '@v[NR]=$*; end{emitp @v[1],"X","Y"}' $indir/abixy

run_mlr --opprint put -q '@v[NR][NR]=$*; end{emitp @v[1]}'         $indir/abixy
run_mlr --opprint put -q '@v[NR][NR]=$*; end{emitp @v[1],"X"}'     $indir/abixy
run_mlr --opprint put -q '@v[NR][NR]=$*; end{emitp @v[1],"X","Y"}' $indir/abixy

run_mlr --opprint put -q '@v[NR][NR]=$*; end{emitp @v[1][1]}'         $indir/abixy
run_mlr --opprint put -q '@v[NR][NR]=$*; end{emitp @v[1][1],"X"}'     $indir/abixy
run_mlr --opprint put -q '@v[NR][NR]=$*; end{emitp @v[1][1],"X","Y"}' $indir/abixy

run_mlr put -q '@a=$a; @b[1]=$b; @c[1][2]=$x; end{emitp all}'                     $indir/abixy-het
run_mlr put -q '@a=$a; @b[1]=$b; @c[1][2]=$x; end{emitp all,"one"}'               $indir/abixy-het
run_mlr put -q '@a=$a; @b[1]=$b; @c[1][2]=$x; end{emitp all,"one","two"}'         $indir/abixy-het
run_mlr put -q '@a=$a; @b[1]=$b; @c[1][2]=$x; end{emitp all,"one","two","three"}' $indir/abixy-het

run_mlr --oxtab put --oflatsep @ -q '@a=$a; @b[1]=$b; @c[1][2]=$x; end{emitp all}' $indir/abixy-het
