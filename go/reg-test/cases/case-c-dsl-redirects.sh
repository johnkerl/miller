# ----------------------------------------------------------------
announce MAPPER TEE REDIRECTS

tee1=$reloutdir/tee1
mkdir -p $tee1

run_mlr --from $indir/abixy tee $tee1/out then nothing
run_cat $tee1/out

run_mlr --from $indir/abixy tee --no-fflush $tee1/out then nothing
run_cat $tee1/out

run_mlr --from $indir/abixy tee -a $tee1/out then nothing
run_cat $tee1/out

run_mlr --from $indir/abixy tee -o json $tee1/out then nothing
run_cat $tee1/out

# ----------------------------------------------------------------
announce DSL TEE REDIRECTS

tee2=$reloutdir/tee2
mkdir -p $tee2

run_mlr put -q 'tee > "'$tee2'/out.".$a, $*' $indir/abixy
run_cat $tee2/out.eks
run_cat $tee2/out.hat
run_cat $tee2/out.pan
run_cat $tee2/out.wye
run_cat $tee2/out.zee

run_mlr put -q --no-fflush 'tee > "'$tee2'/out.".$a, $*' $indir/abixy
run_cat $tee2/out.eks
run_cat $tee2/out.hat
run_cat $tee2/out.pan
run_cat $tee2/out.wye
run_cat $tee2/out.zee

run_mlr put -q 'tee >> "'$tee2'/out.".$a, $*' $indir/abixy
run_cat $tee2/out.eks
run_cat $tee2/out.hat
run_cat $tee2/out.pan
run_cat $tee2/out.wye
run_cat $tee2/out.zee

run_mlr put -q -o json 'tee > "'$tee2'/out.".$a, $*' $indir/abixy
run_cat $tee2/out.eks
run_cat $tee2/out.hat
run_cat $tee2/out.pan
run_cat $tee2/out.wye
run_cat $tee2/out.zee

run_mlr put -q 'tee | "tr \[a-z\] \[A-Z\]", $*' $indir/abixy

run_mlr put -q -o json 'tee | "tr \[a-z\] \[A-Z\]", $*' $indir/abixy

touch $tee2/err1
run_mlr put -q 'tee > stdout, $*' $indir/abixy 2> $tee2/err1
run_cat $tee2/err1

touch $tee2/err2
run_mlr put -q 'tee > stderr, $*' $indir/abixy 2> $tee2/err2
run_cat $tee2/err2

# ----------------------------------------------------------------
announce DSL PRINT REDIRECTS

print1=$reloutdir/print1
mkdir -p $print1

run_mlr put -q 'print > "'$print1'/out.".$a, "abi:".$a.$b.$i' $indir/abixy
run_cat $print1/out.eks
run_cat $print1/out.hat
run_cat $print1/out.pan
run_cat $print1/out.wye
run_cat $print1/out.zee

run_mlr put -q 'print > "'$print1'/out.".$a, "abi:".$a.$b.$i' $indir/abixy
run_cat $print1/out.eks
run_cat $print1/out.hat
run_cat $print1/out.pan
run_cat $print1/out.wye
run_cat $print1/out.zee

run_mlr put -q 'print >> "'$print1'/out.".$a, "abi:".$a.$b.$i' $indir/abixy
run_cat $print1/out.eks
run_cat $print1/out.hat
run_cat $print1/out.pan
run_cat $print1/out.wye
run_cat $print1/out.zee

run_mlr put -q 'print | "tr \[a-z\] \[A-Z\]",  "abi:".$a.$b.$i' $indir/abixy

touch $print1/err1
run_mlr put -q 'print > stdout, "abi:".$a.$b.$i' $indir/abixy 2> $print1/err1
run_cat $print1/err1

touch $print1/err2
run_mlr put -q 'print > stderr, "abi:".$a.$b.$i' $indir/abixy 2> $print1/err2
run_cat $print1/err2

# ----------------------------------------------------------------
announce DSL PRINTN REDIRECTS

printn1=$reloutdir/printn1
mkdir -p $printn1

run_mlr put -q 'printn > "'$printn1'/out.".$a, "abi:".$a.$b.$i' $indir/abixy
run_cat $printn1/out.eks
run_cat $printn1/out.hat
run_cat $printn1/out.pan
run_cat $printn1/out.wye
run_cat $printn1/out.zee

run_mlr put -q 'printn > "'$printn1'/out.".$a, "abi:".$a.$b.$i' $indir/abixy
run_cat $printn1/out.eks
run_cat $printn1/out.hat
run_cat $printn1/out.pan
run_cat $printn1/out.wye
run_cat $printn1/out.zee

run_mlr put -q 'printn >> "'$printn1'/out.".$a, "abi:".$a.$b.$i' $indir/abixy
run_cat $printn1/out.eks
run_cat $printn1/out.hat
run_cat $printn1/out.pan
run_cat $printn1/out.wye
run_cat $printn1/out.zee

run_mlr put -q 'printn | "tr \[a-z\] \[A-Z\]",  "abi:".$a.$b.$i' $indir/abixy

touch $printn1/err1
run_mlr put -q 'printn > stdout, "abi:".$a.$b.$i' $indir/abixy 2> $printn1/err1
run_cat $printn1/err1

touch $printn1/err2
run_mlr put -q 'printn > stderr, "abi:".$a.$b.$i' $indir/abixy 2> $printn1/err2
run_cat $printn1/err2

# ----------------------------------------------------------------
announce DSL DUMP REDIRECTS

dump1=$reloutdir/dump1
mkdir -p $dump1

run_mlr put -q '@v=$*; dump > "'$dump1'/out.".$a' $indir/abixy
run_cat $dump1/out.eks
run_cat $dump1/out.hat
run_cat $dump1/out.pan
run_cat $dump1/out.wye
run_cat $dump1/out.zee

run_mlr put -q '@v=$*; dump > "'$dump1'/out.".$a' $indir/abixy
run_cat $dump1/out.eks
run_cat $dump1/out.hat
run_cat $dump1/out.pan
run_cat $dump1/out.wye
run_cat $dump1/out.zee

run_mlr put -q '@v=$*; dump >> "'$dump1'/out.".$a' $indir/abixy
run_cat $dump1/out.eks
run_cat $dump1/out.hat
run_cat $dump1/out.pan
run_cat $dump1/out.wye
run_cat $dump1/out.zee

run_mlr put -q '@v=$*; dump | "tr \[a-z\] \[A-Z\]"' $indir/abixy

touch $dump1/err1
run_mlr put -q '@v[NR] = $*; NR == 2 { dump > stdout }' $indir/abixy 2> $dump1/err1
run_cat $dump1/err1

touch $dump1/err2
run_mlr put -q '@v[NR] = $*; NR == 2 { dump > stderr }' $indir/abixy 2> $dump1/err2
run_cat $dump1/err2

# ----------------------------------------------------------------
announce DSL EMITF REDIRECTS

emitf1=$reloutdir/emitf1
mkdir -p $emitf1

run_mlr put -q '@a=$a; @b=$b; emitf > "'$emitf1'/out.".$a.$b, @a, @b' $indir/abixy
run_cat $emitf1/out.ekspan
run_cat $emitf1/out.ekswye
run_cat $emitf1/out.ekszee
run_cat $emitf1/out.hatwye
run_cat $emitf1/out.panpan
run_cat $emitf1/out.panwye
run_cat $emitf1/out.wyepan
run_cat $emitf1/out.wyewye
run_cat $emitf1/out.zeepan
run_cat $emitf1/out.zeewye

run_mlr put -q '@a=$a; @b=$b; emitf > "'$emitf1'/out.".$a.$b, @a, @b' $indir/abixy
run_cat $emitf1/out.ekspan
run_cat $emitf1/out.ekswye
run_cat $emitf1/out.ekszee
run_cat $emitf1/out.hatwye
run_cat $emitf1/out.panpan
run_cat $emitf1/out.panwye
run_cat $emitf1/out.wyepan
run_cat $emitf1/out.wyewye
run_cat $emitf1/out.zeepan
run_cat $emitf1/out.zeewye

run_mlr put -q '@a=$a; @b=$b; emitf >> "'$emitf1'/out.".$a.$b, @a, @b' $indir/abixy
run_cat $emitf1/out.ekspan
run_cat $emitf1/out.ekswye
run_cat $emitf1/out.ekszee
run_cat $emitf1/out.hatwye
run_cat $emitf1/out.panpan
run_cat $emitf1/out.panwye
run_cat $emitf1/out.wyepan
run_cat $emitf1/out.wyewye
run_cat $emitf1/out.zeepan
run_cat $emitf1/out.zeewye

run_mlr put -q -o json '@a=$a; @b=$b; emitf > "'$emitf1'/out.".$a.$b, @a, @b' $indir/abixy
run_cat $emitf1/out.ekspan
run_cat $emitf1/out.ekswye
run_cat $emitf1/out.ekszee
run_cat $emitf1/out.hatwye
run_cat $emitf1/out.panpan
run_cat $emitf1/out.panwye
run_cat $emitf1/out.wyepan
run_cat $emitf1/out.wyewye
run_cat $emitf1/out.zeepan
run_cat $emitf1/out.zeewye

run_mlr put -q '@a=$a; @b=$b; emitf | "tr \[a-z\] \[A-Z\]", @a, @b' $indir/abixy

touch $emitf1/err1
run_mlr put -q '@a=$a; @b=$b; emitf > stdout, @a, @b' $indir/abixy 2> $emitf1/err1
run_cat $emitf1/err1

touch $emitf1/err2
run_mlr put -q '@a=$a; @b=$b; emitf > stderr, @a, @b' $indir/abixy 2> $emitf1/err2
run_cat $emitf1/err2

# ----------------------------------------------------------------
announce DSL EMITP REDIRECTS

emitp1=$reloutdir/emitp1
mkdir -p $emitp1

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp > "'$emitp1'/out.".$a.$b, @a' $indir/abixy
run_cat $emitp1/out.ekspan
run_cat $emitp1/out.ekswye
run_cat $emitp1/out.panpan
run_cat $emitp1/out.wyewye

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp > "'$emitp1'/out.".$a.$b, @a' $indir/abixy
run_cat $emitp1/out.ekspan
run_cat $emitp1/out.ekswye
run_cat $emitp1/out.panpan
run_cat $emitp1/out.wyewye

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp >> "'$emitp1'/out.".$a.$b, @a' $indir/abixy
run_cat $emitp1/out.ekspan
run_cat $emitp1/out.ekswye
run_cat $emitp1/out.panpan
run_cat $emitp1/out.wyewye

run_mlr head -n 4 then put -q -o json '@a[NR]=$a; @b[NR]=$b; emitp > "'$emitp1'/out.".$a.$b, @a' $indir/abixy
run_cat $emitp1/out.ekspan
run_cat $emitp1/out.ekswye
run_cat $emitp1/out.panpan
run_cat $emitp1/out.wyewye

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp | "tr \[a-z\] \[A-Z\]", @a' $indir/abixy

touch $emitp1/err1
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp > stdout, @a' $indir/abixy 2> $emitp1/err1
run_cat $emitp1/err1

touch $emitp1/err2
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp > stderr, @a' $indir/abixy 2> $emitp1/err2
run_cat $emitp1/err2


emitp2=$reloutdir/emitp2
mkdir -p $emitp2

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp > "'$emitp2'/out.".$a.$b, @a, "NR"' $indir/abixy
run_cat $emitp2/out.ekspan
run_cat $emitp2/out.ekswye
run_cat $emitp2/out.panpan
run_cat $emitp2/out.wyewye

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp > "'$emitp2'/out.".$a.$b, @a, "NR"' $indir/abixy
run_cat $emitp2/out.ekspan
run_cat $emitp2/out.ekswye
run_cat $emitp2/out.panpan
run_cat $emitp2/out.wyewye

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp >> "'$emitp2'/out.".$a.$b, @a, "NR"' $indir/abixy
run_cat $emitp2/out.ekspan
run_cat $emitp2/out.ekswye
run_cat $emitp2/out.panpan
run_cat $emitp2/out.wyewye

run_mlr head -n 4 then put -q -o json '@a[NR]=$a; @b[NR]=$b; emitp > "'$emitp2'/out.".$a.$b, @a, "NR"' $indir/abixy
run_cat $emitp2/out.ekspan
run_cat $emitp2/out.ekswye
run_cat $emitp2/out.panpan
run_cat $emitp2/out.wyewye

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp | "tr \[a-z\] \[A-Z\]", @a, "NR"' $indir/abixy

touch $emitp2/err1
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp > stdout, @a, "NR"' $indir/abixy 2> $emitp2/err1
run_cat $emitp2/err1

touch $emitp2/err2
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp > stderr, @a, "NR"' $indir/abixy 2> $emitp2/err2
run_cat $emitp2/err2


emitp3=$reloutdir/emitp3
mkdir -p $emitp3

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp > "'$emitp3'/out.".$a.$b, (@a, @b)' $indir/abixy
run_cat $emitp3/out.ekspan
run_cat $emitp3/out.ekswye
run_cat $emitp3/out.panpan
run_cat $emitp3/out.wyewye

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp > "'$emitp3'/out.".$a.$b, (@a, @b)' $indir/abixy
run_cat $emitp3/out.ekspan
run_cat $emitp3/out.ekswye
run_cat $emitp3/out.panpan
run_cat $emitp3/out.wyewye

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp >> "'$emitp3'/out.".$a.$b, (@a, @b)' $indir/abixy
run_cat $emitp3/out.ekspan
run_cat $emitp3/out.ekswye
run_cat $emitp3/out.panpan
run_cat $emitp3/out.wyewye

run_mlr head -n 4 then put -q -o json '@a[NR]=$a; @b[NR]=$b; emitp > "'$emitp3'/out.".$a.$b, (@a, @b)' $indir/abixy
run_cat $emitp3/out.ekspan
run_cat $emitp3/out.ekswye
run_cat $emitp3/out.panpan
run_cat $emitp3/out.wyewye

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp | "tr \[a-z\] \[A-Z\]", (@a, @b)' $indir/abixy

touch $emitp3/err1
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp > stdout, (@a, @b)' $indir/abixy 2> $emitp3/err1
run_cat $emitp3/err1

touch $emitp3/err2
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp > stderr, (@a, @b)' $indir/abixy 2> $emitp3/err2
run_cat $emitp3/err2


emitp4=$reloutdir/emitp4
mkdir -p $emitp4

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp > "'$emitp4'/out.".$a.$b, (@a, @b), "NR"' $indir/abixy
run_cat $emitp4/out.ekspan
run_cat $emitp4/out.ekswye
run_cat $emitp4/out.panpan
run_cat $emitp4/out.wyewye

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp > "'$emitp4'/out.".$a.$b, (@a, @b), "NR"' $indir/abixy
run_cat $emitp4/out.ekspan
run_cat $emitp4/out.ekswye
run_cat $emitp4/out.panpan
run_cat $emitp4/out.wyewye

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp >> "'$emitp4'/out.".$a.$b, (@a, @b), "NR"' $indir/abixy
run_cat $emitp4/out.ekspan
run_cat $emitp4/out.ekswye
run_cat $emitp4/out.panpan
run_cat $emitp4/out.wyewye

run_mlr head -n 4 then put -q -o json '@a[NR]=$a; @b[NR]=$b; emitp > "'$emitp4'/out.".$a.$b, (@a, @b), "NR"' $indir/abixy
run_cat $emitp4/out.ekspan
run_cat $emitp4/out.ekswye
run_cat $emitp4/out.panpan
run_cat $emitp4/out.wyewye

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp | "tr \[a-z\] \[A-Z\]", (@a, @b), "NR"' $indir/abixy

touch $emitp4/err1
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp > stdout, (@a, @b), "NR"' $indir/abixy 2> $emitp4/err1
run_cat $emitp4/err1

touch $emitp4/err2
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp > stderr, (@a, @b), "NR"' $indir/abixy 2> $emitp4/err2
run_cat $emitp4/err2


emitp5=$reloutdir/emitp5
mkdir -p $emitp5
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp > "'$emitp5'/out.".$a.$b, @*' $indir/abixy
run_cat $emitp5/out.ekspan
run_cat $emitp5/out.ekswye
run_cat $emitp5/out.panpan
run_cat $emitp5/out.wyewye
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp | "tr \[a-z\] \[A-Z\]", @*' $indir/abixy

touch $emitp5/err1
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp > stdout, @*' $indir/abixy 2> $emitp5/err1
run_cat $emitp5/err1

touch $emitp5/err2
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp > stderr, @*' $indir/abixy 2> $emitp5/err2
run_cat $emitp5/err2


emitp6=$reloutdir/emitp6
mkdir -p $emitp6
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp > "'$emitp6'/out.".$a.$b, all' $indir/abixy
run_cat $emitp6/out.ekspan
run_cat $emitp6/out.ekswye
run_cat $emitp6/out.panpan
run_cat $emitp6/out.wyewye
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp | "tr \[a-z\] \[A-Z\]", all' $indir/abixy

touch $emitp6/err1
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp > stdout, all' $indir/abixy 2> $emitp6/err1
run_cat $emitp6/err1

touch $emitp6/err2
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp > stderr, all' $indir/abixy 2> $emitp6/err2
run_cat $emitp6/err2


emitp7=$reloutdir/emitp7
mkdir -p $emitp7
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp > "'$emitp7'/out.".$a.$b, @*, "NR"' $indir/abixy
run_cat $emitp7/out.ekspan
run_cat $emitp7/out.ekswye
run_cat $emitp7/out.panpan
run_cat $emitp7/out.wyewye
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp | "tr \[a-z\] \[A-Z\]", @*, "NR"' $indir/abixy

touch $emitp7/err1
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp > stdout, @*, "NR"' $indir/abixy 2> $emitp7/err1
run_cat $emitp7/err1

touch $emitp7/err2
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp > stderr, @*, "NR"' $indir/abixy 2> $emitp7/err2
run_cat $emitp7/err2


emitp8=$reloutdir/emitp8
mkdir -p $emitp8
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp > "'$emitp8'/out.".$a.$b, all, "NR"' $indir/abixy
run_cat $emitp8/out.ekspan
run_cat $emitp8/out.ekswye
run_cat $emitp8/out.panpan
run_cat $emitp8/out.wyewye
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp | "tr \[a-z\] \[A-Z\]", all, "NR"' $indir/abixy

touch $emitp8/err1
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp > stdout, all, "NR"' $indir/abixy 2> $emitp8/err1
run_cat $emitp8/err1

touch $emitp8/err2
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emitp > stderr, all, "NR"' $indir/abixy 2> $emitp8/err2
run_cat $emitp8/err2


# ----------------------------------------------------------------
announce DSL EMIT REDIRECTS

emit1=$reloutdir/emit1
mkdir -p $emit1

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > "'$emit1'/out.".$a.$b, @a' $indir/abixy
run_cat $emit1/out.ekspan
run_cat $emit1/out.ekswye
run_cat $emit1/out.panpan
run_cat $emit1/out.wyewye

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > "'$emit1'/out.".$a.$b, @a' $indir/abixy
run_cat $emit1/out.ekspan
run_cat $emit1/out.ekswye
run_cat $emit1/out.panpan
run_cat $emit1/out.wyewye

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit >> "'$emit1'/out.".$a.$b, @a' $indir/abixy
run_cat $emit1/out.ekspan
run_cat $emit1/out.ekswye
run_cat $emit1/out.panpan
run_cat $emit1/out.wyewye

run_mlr head -n 4 then put -q -o json '@a[NR]=$a; @b[NR]=$b; emit > "'$emit1'/out.".$a.$b, @a' $indir/abixy
run_cat $emit1/out.ekspan
run_cat $emit1/out.ekswye
run_cat $emit1/out.panpan
run_cat $emit1/out.wyewye

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit | "tr \[a-z\] \[A-Z\]", @a' $indir/abixy

touch $emit1/err1
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > stdout, @a' $indir/abixy 2> $emit1/err1
run_cat $emit1/err1

touch $emit1/err2
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > stderr, @a' $indir/abixy 2> $emit1/err2
run_cat $emit1/err2


emit2=$reloutdir/emit2
mkdir -p $emit2

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > "'$emit2'/out.".$a.$b, @a, "NR"' $indir/abixy
run_cat $emit2/out.ekspan
run_cat $emit2/out.ekswye
run_cat $emit2/out.panpan
run_cat $emit2/out.wyewye

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > "'$emit2'/out.".$a.$b, @a, "NR"' $indir/abixy
run_cat $emit2/out.ekspan
run_cat $emit2/out.ekswye
run_cat $emit2/out.panpan
run_cat $emit2/out.wyewye

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit >> "'$emit2'/out.".$a.$b, @a, "NR"' $indir/abixy
run_cat $emit2/out.ekspan
run_cat $emit2/out.ekswye
run_cat $emit2/out.panpan
run_cat $emit2/out.wyewye

run_mlr head -n 4 then put -q -o pprint '@a[NR]=$a; @b[NR]=$b; emit > "'$emit2'/out.".$a.$b, @a, "NR"' $indir/abixy
run_cat $emit2/out.ekspan
run_cat $emit2/out.ekswye
run_cat $emit2/out.panpan
run_cat $emit2/out.wyewye

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit | "tr \[a-z\] \[A-Z\]", @a, "NR"' $indir/abixy

touch $emit2/err1
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > stdout, @a, "NR"' $indir/abixy 2> $emit2/err1
run_cat $emit2/err1

touch $emit2/err2
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > stderr, @a, "NR"' $indir/abixy 2> $emit2/err2
run_cat $emit2/err2


emit3=$reloutdir/emit3
mkdir -p $emit3

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > "'$emit3'/out.".$a.$b, (@a, @b)' $indir/abixy
run_cat $emit3/out.ekspan
run_cat $emit3/out.ekswye
run_cat $emit3/out.panpan
run_cat $emit3/out.wyewye

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > "'$emit3'/out.".$a.$b, (@a, @b)' $indir/abixy
run_cat $emit3/out.ekspan
run_cat $emit3/out.ekswye
run_cat $emit3/out.panpan
run_cat $emit3/out.wyewye

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit >> "'$emit3'/out.".$a.$b, (@a, @b)' $indir/abixy
run_cat $emit3/out.ekspan
run_cat $emit3/out.ekswye
run_cat $emit3/out.panpan
run_cat $emit3/out.wyewye

run_mlr head -n 4 then put -q --oxtab '@a[NR]=$a; @b[NR]=$b; emit > "'$emit3'/out.".$a.$b, (@a, @b)' $indir/abixy
run_cat $emit3/out.ekspan
run_cat $emit3/out.ekswye
run_cat $emit3/out.panpan
run_cat $emit3/out.wyewye

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit | "tr \[a-z\] \[A-Z\]", (@a, @b)' $indir/abixy

touch $emit3/err1
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > stdout, (@a, @b)' $indir/abixy 2> $emit3/err1
run_cat $emit3/err1

touch $emit3/err2
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > stderr, (@a, @b)' $indir/abixy 2> $emit3/err2
run_cat $emit3/err2


emit4=$reloutdir/emit4
mkdir -p $emit4

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > "'$emit4'/out.".$a.$b, (@a, @b), "NR"' $indir/abixy
run_cat $emit4/out.ekspan
run_cat $emit4/out.ekswye
run_cat $emit4/out.panpan
run_cat $emit4/out.wyewye

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > "'$emit4'/out.".$a.$b, (@a, @b), "NR"' $indir/abixy
run_cat $emit4/out.ekspan
run_cat $emit4/out.ekswye
run_cat $emit4/out.panpan
run_cat $emit4/out.wyewye

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit >> "'$emit4'/out.".$a.$b, (@a, @b), "NR"' $indir/abixy
run_cat $emit4/out.ekspan
run_cat $emit4/out.ekswye
run_cat $emit4/out.panpan
run_cat $emit4/out.wyewye

run_mlr head -n 4 then put -q --ojson '@a[NR]=$a; @b[NR]=$b; emit > "'$emit4'/out.".$a.$b, (@a, @b), "NR"' $indir/abixy
run_cat $emit4/out.ekspan
run_cat $emit4/out.ekswye
run_cat $emit4/out.panpan
run_cat $emit4/out.wyewye

run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit | "tr \[a-z\] \[A-Z\]", (@a, @b), "NR"' $indir/abixy

touch $emit4/err1
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > stdout, (@a, @b), "NR"' $indir/abixy 2> $emit4/err1
run_cat $emit4/err1

touch $emit4/err2
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > stderr, (@a, @b), "NR"' $indir/abixy 2> $emit4/err2
run_cat $emit4/err2


emit5=$reloutdir/emit5
mkdir -p $emit5
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > "'$emit5'/out.".$a.$b, @*' $indir/abixy
run_cat $emit5/out.ekspan
run_cat $emit5/out.ekswye
run_cat $emit5/out.panpan
run_cat $emit5/out.wyewye
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit | "tr \[a-z\] \[A-Z\]", @*' $indir/abixy

touch $emit5/err1
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > stdout, @*' $indir/abixy 2> $emit5/err1
run_cat $emit5/err1

touch $emit5/err2
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > stderr, @*' $indir/abixy 2> $emit5/err2
run_cat $emit5/err2


emit6=$reloutdir/emit6
mkdir -p $emit6
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > "'$emit6'/out.".$a.$b, all' $indir/abixy
run_cat $emit6/out.ekspan
run_cat $emit6/out.ekswye
run_cat $emit6/out.panpan
run_cat $emit6/out.wyewye
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit | "tr \[a-z\] \[A-Z\]", all' $indir/abixy

touch $emit6/err1
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > stdout, all' $indir/abixy 2> $emit6/err1
run_cat $emit6/err1

touch $emit6/err2
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > stderr, all' $indir/abixy 2> $emit6/err2
run_cat $emit6/err2


emit7=$reloutdir/emit7
mkdir -p $emit7
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > "'$emit7'/out.".$a.$b, @*, "NR"' $indir/abixy
run_cat $emit7/out.ekspan
run_cat $emit7/out.ekswye
run_cat $emit7/out.panpan
run_cat $emit7/out.wyewye
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit | "tr \[a-z\] \[A-Z\]", @*, "NR"' $indir/abixy

touch $emit7/err1
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > stdout, @*, "NR"' $indir/abixy 2> $emit7/err1
run_cat $emit7/err1

touch $emit7/err2
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > stderr, @*, "NR"' $indir/abixy 2> $emit7/err2
run_cat $emit7/err2


emit8=$reloutdir/emit8
mkdir -p $emit8
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > "'$emit8'/out.".$a.$b, all, "NR"' $indir/abixy
run_cat $emit8/out.ekspan
run_cat $emit8/out.ekswye
run_cat $emit8/out.panpan
run_cat $emit8/out.wyewye
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit | "tr \[a-z\] \[A-Z\]", all, "NR"' $indir/abixy

touch $emit8/err1
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > stdout, all, "NR"' $indir/abixy 2> $emit8/err1
run_cat $emit8/err1

touch $emit8/err2
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > stderr, all, "NR"' $indir/abixy 2> $emit8/err2
run_cat $emit8/err2


emit9=$reloutdir/emit9
mkdir -p $emit9
run_mlr head -n 4 then put -q 'emit > "'$emit9'/out.".$a.$b, $*' $indir/abixy
run_cat $emit9/out.ekspan
run_cat $emit9/out.ekswye
run_cat $emit9/out.panpan
run_cat $emit9/out.wyewye
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit | "tr \[a-z\] \[A-Z\]", $*, "NR"' $indir/abixy

touch $emit9/err1
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > stdout, $*, "NR"' $indir/abixy 2> $emit9/err1
run_cat $emit9/err1

touch $emit9/err2
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > stderr, $*, "NR"' $indir/abixy 2> $emit9/err2
run_cat $emit9/err2


emit10=$reloutdir/emit10
mkdir -p $emit10
run_mlr head -n 4 then put -q 'emit > "'$emit10'/out.".$a.$b, mapexcept($*, "a", "b")' $indir/abixy
run_cat $emit10/out.ekspan
run_cat $emit10/out.ekswye
run_cat $emit10/out.panpan
run_cat $emit10/out.wyewye
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit | "tr \[a-z\] \[A-Z\]", mapexcept($*, "a", "b"), "NR"' $indir/abixy

touch $emit10/err1
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > stdout, mapexcept($*, "a", "b"), "NR"' $indir/abixy 2> $emit10/err1
run_cat $emit10/err1

touch $emit10/err2
run_mlr head -n 4 then put -q '@a[NR]=$a; @b[NR]=$b; emit > stderr, mapexcept($*, "a", "b"), "NR"' $indir/abixy 2> $emit10/err2
run_cat $emit10/err2
