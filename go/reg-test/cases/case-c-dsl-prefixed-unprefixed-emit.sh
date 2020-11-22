run_mlr --oxtab put -q '@sum     += $x; end{dump;emitp  @sum     }'  $indir/abixy-wide
run_mlr --oxtab put -q '@sum     += $x; end{dump;emit @sum     }'  $indir/abixy-wide


run_mlr --oxtab put -q '@sum[$a] += $x; end{dump;emitp  @sum     }'  $indir/abixy-wide
run_mlr --oxtab put -q '@sum[$a] += $x; end{dump;emit @sum     }'  $indir/abixy-wide

run_mlr --oxtab put -q '@sum[$a] += $x; end{dump;emitp  @sum, "a"}'  $indir/abixy-wide
run_mlr --oxtab put -q '@sum[$a] += $x; end{dump;emit @sum, "a"}'  $indir/abixy-wide


run_mlr --oxtab put -q '@sum[$a][$b] += $x; end{dump;emitp  @sum     }'  $indir/abixy-wide
run_mlr --oxtab put -q '@sum[$a][$b] += $x; end{dump;emit @sum     }'  $indir/abixy-wide

run_mlr --oxtab put -q '@sum[$a][$b] += $x; end{dump;emitp  @sum, "a"}'  $indir/abixy-wide
run_mlr --oxtab put -q '@sum[$a][$b] += $x; end{dump;emit @sum, "a"}'  $indir/abixy-wide

run_mlr --opprint put -q '@sum[$a][$b] += $x; end{dump;emitp  @sum, "a", "b"}'  $indir/abixy-wide
run_mlr --opprint put -q '@sum[$a][$b] += $x; end{dump;emit @sum, "a", "b"}'  $indir/abixy-wide


run_mlr --oxtab put -q '@sum[$a][$b][$a.$b] += $x; end{dump;emitp  @sum     }'  $indir/abixy-wide
run_mlr --oxtab put -q '@sum[$a][$b][$a.$b] += $x; end{dump;emit @sum     }'  $indir/abixy-wide

run_mlr --oxtab put -q '@sum[$a][$b][$a.$b] += $x; end{dump;emitp  @sum, "a"}'  $indir/abixy-wide
run_mlr --oxtab put -q '@sum[$a][$b][$a.$b] += $x; end{dump;emit @sum, "a"}'  $indir/abixy-wide

run_mlr --opprint put -q '@sum[$a][$b][$a.$b] += $x; end{dump;emitp  @sum, "a", "b"}'  $indir/abixy-wide
run_mlr --opprint put -q '@sum[$a][$b][$a.$b] += $x; end{dump;emit @sum, "a", "b"}'  $indir/abixy-wide

run_mlr --opprint put -q '@sum[$a][$b][$a.$b] += $x; end{dump;emitp  @sum, "a", "b", "ab"}'  $indir/abixy-wide
run_mlr --opprint put -q '@sum[$a][$b][$a.$b] += $x; end{dump;emit @sum, "a", "b", "ab"}'  $indir/abixy-wide



run_mlr --oxtab head -n 2  then put -q '@v       =  $*; end{dump;emitp  @v}'         $indir/abixy
run_mlr --oxtab head -n 2  then put -q '@v       =  $*; end{dump;emit @v}'         $indir/abixy


run_mlr --oxtab head -n 2  then put -q '@v[NR]   =  $*; end{dump;emitp  @v        }' $indir/abixy
run_mlr --oxtab head -n 2  then put -q '@v[NR]   =  $*; end{dump;emit @v        }' $indir/abixy

run_mlr --opprint head -n 2  then put -q '@v[NR]   =  $*; end{dump;emitp  @v,   "I"}' $indir/abixy
run_mlr --opprint head -n 2  then put -q '@v[NR]   =  $*; end{dump;emit @v,   "I"}' $indir/abixy


run_mlr --oxtab head -n 2  then put -q '@v[NR][NR]   =  $*; end{dump;emitp  @v        }' $indir/abixy
run_mlr --oxtab head -n 2  then put -q '@v[NR][NR]   =  $*; end{dump;emit @v        }' $indir/abixy

run_mlr --opprint head -n 2  then put -q '@v[NR][NR]   =  $*; end{dump;emitp  @v,   "I"}' $indir/abixy
run_mlr --opprint head -n 2  then put -q '@v[NR][NR]   =  $*; end{dump;emit @v,   "I"}' $indir/abixy

run_mlr --opprint head -n 2  then put -q '@v[NR][NR]   =  $*; end{dump;emitp  @v,   "I", "J"}' $indir/abixy
run_mlr --opprint head -n 2  then put -q '@v[NR][NR]   =  $*; end{dump;emit @v,   "I", "J"}' $indir/abixy


run_mlr --oxtab head -n 2  then put -q '@v[NR][NR][NR]   =  $*; end{dump;emitp  @v        }' $indir/abixy
run_mlr --oxtab head -n 2  then put -q '@v[NR][NR][NR]   =  $*; end{dump;emit @v        }' $indir/abixy

run_mlr --opprint head -n 2  then put -q '@v[NR][NR][NR]   =  $*; end{dump;emitp  @v,   "I"}' $indir/abixy
run_mlr --opprint head -n 2  then put -q '@v[NR][NR][NR]   =  $*; end{dump;emit @v,   "I"}' $indir/abixy

run_mlr --opprint head -n 2  then put -q '@v[NR][NR][NR]   =  $*; end{dump;emitp  @v,   "I", "J"}' $indir/abixy
run_mlr --opprint head -n 2  then put -q '@v[NR][NR][NR]   =  $*; end{dump;emit @v,   "I", "J"}' $indir/abixy

run_mlr --opprint head -n 2  then put -q '@v[NR][NR][NR]   =  $*; end{dump;emitp  @v,   "I", "J", "K"}' $indir/abixy
run_mlr --opprint head -n 2  then put -q '@v[NR][NR][NR]   =  $*; end{dump;emit @v,   "I", "J", "K"}' $indir/abixy
