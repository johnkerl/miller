# ----------------------------------------------------------------
announce OOSVAR-FROM-SREC ASSIGNMENT

run_mlr put -v '@v     = $*' /dev/null
run_mlr put -v '@v[1]  = $*' /dev/null
run_mlr put -v '@v[$2] = $*' /dev/null
run_mlr put -v 'NR == 3 {@v     = $*}' /dev/null
run_mlr put -v 'NR == 3 {@v[1]  = $*}' /dev/null
run_mlr put -v 'NR == 3 {@v[$2] = $*}' /dev/null

run_mlr --oxtab put -q '@v = $*; end {emitp @v }' $indir/abixy-het

run_mlr --oxtab put -q '@v[$a] = $*; end {emitp @v      }' $indir/abixy-het
run_mlr --oxtab put -q '@v[$a] = $*; end {emitp @v, "a" }' $indir/abixy-het

run_mlr --oxtab put -q '@v[$a][$b] = $*; end {emitp @v          }' $indir/abixy-het
run_mlr --oxtab put -q '@v[$a][$b] = $*; end {emitp @v, "a"     }' $indir/abixy-het
run_mlr --oxtab put -q '@v[$a][$b] = $*; end {emitp @v, "a", "b"}' $indir/abixy-het

# ----------------------------------------------------------------
announce SREC-FROM-OOSVAR ASSIGNMENT

run_mlr put -v '$* = @v    ' /dev/null
run_mlr put -v '$* = @v[1] ' /dev/null
run_mlr put -v '$* = @v[$2]' /dev/null
run_mlr put -v 'NR == 3 {$* = @v    }' /dev/null
run_mlr put -v 'NR == 3 {$* = @v[1] }' /dev/null
run_mlr put -v 'NR == 3 {$* = @v[$2]}' /dev/null

run_mlr put '@v[NR] = $a; NR == 7 { @v = $*} ; $* = @v' $indir/abixy-het

# ----------------------------------------------------------------
announce OOSVAR-FROM-OOSVAR ASSIGNMENT

run_mlr put -v '@u    = @v'    /dev/null
run_mlr put -v '@u    = @v[1]' /dev/null
run_mlr put -v '@u[2] = @v'    /dev/null
run_mlr put -v '@u[2] = @v[1]' /dev/null

run_mlr put -v 'begin { @u    = @v }'    /dev/null
run_mlr put -v 'begin { @u    = @v[1] }' /dev/null
run_mlr put -v 'begin { @u[2] = @v }'    /dev/null
run_mlr put -v 'begin { @u[2] = @v[1] }' /dev/null

run_mlr put -v 'NR == 3 { @u    = @v }'    /dev/null
run_mlr put -v 'NR == 3 { @u    = @v[1] }' /dev/null
run_mlr put -v 'NR == 3 { @u[2] = @v }'    /dev/null
run_mlr put -v 'NR == 3 { @u[2] = @v[1] }' /dev/null

run_mlr put -v 'end { @u    = @v }'    /dev/null
run_mlr put -v 'end { @u    = @v[1] }' /dev/null
run_mlr put -v 'end { @u[2] = @v }'    /dev/null
run_mlr put -v 'end { @u[2] = @v[1] }' /dev/null


run_mlr put -q '@s    += $i; @t=@s;             end{dump; emitp@s; emitp @t}' $indir/abixy

run_mlr put -q '@s[1] += $i; @t=@s;             end{dump; emitp@s; emitp @t}' $indir/abixy
run_mlr put -q '@s[1] += $i; @t=@s[1];          end{dump; emitp@s; emitp @t}' $indir/abixy

run_mlr put -q '@s[1] += $i; @t[3]=@s;          end{dump; emitp@s; emitp @t}' $indir/abixy
run_mlr put -q '@s[1] += $i; @t[3]=@s[1];       end{dump; emitp@s; emitp @t}' $indir/abixy

run_mlr put -q '@s[1][2] += $i; @t=@s;             end{dump; emitp@s; emitp @t}' $indir/abixy
run_mlr put -q '@s[1][2] += $i; @t=@s[1];          end{dump; emitp@s; emitp @t}' $indir/abixy
run_mlr put -q '@s[1][2] += $i; @t=@s[1][2];       end{dump; emitp@s; emitp @t}' $indir/abixy

run_mlr put -q '@s[1][2] += $i; @t[3]=@s;          end{dump; emitp@s; emitp @t}' $indir/abixy
run_mlr put -q '@s[1][2] += $i; @t[3]=@s[1];       end{dump; emitp@s; emitp @t}' $indir/abixy
run_mlr put -q '@s[1][2] += $i; @t[3]=@s[1][2];    end{dump; emitp@s; emitp @t}' $indir/abixy

run_mlr put -q '@s[1][2] += $i; @t[3][4]=@s;       end{dump; emitp@s; emitp @t}' $indir/abixy
run_mlr put -q '@s[1][2] += $i; @t[3][4]=@s[1];    end{dump; emitp@s; emitp @t}' $indir/abixy
run_mlr put -q '@s[1][2] += $i; @t[3][4]=@s[1][2]; end{dump; emitp@s; emitp @t}' $indir/abixy

run_mlr --opprint put -q '@s[NR][NR] = $i/100; @t[NR*10]=@s; end{emitp@s,"A","B"; emitp @t,"C","D","E"}' $indir/abixy
