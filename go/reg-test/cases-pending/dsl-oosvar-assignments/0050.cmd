mlr --opprint put -q '@s[NR][NR] = $i/100; @t[NR*10]=@s; end{emitp@s,"A","B"; emitp @t,"C","D","E"}' reg-test/input/abixy
