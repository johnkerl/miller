
run_mlr --xtab put    '$quot=$pf1/10.0;$iquot=$pf1//10.0;$mod=$pf1%10.0' $indir/mixed-types.xtab
run_mlr --xtab put    '$quot=$pi1/10  ;$iquot=$pi1//10  ;$mod=$pi1%10  ' $indir/mixed-types.xtab
run_mlr --xtab put    '$quot=$nf1/10.0;$iquot=$nf1//10.0;$mod=$nf1%10.0' $indir/mixed-types.xtab
run_mlr --xtab put    '$quot=$ni1/10  ;$iquot=$ni1//10  ;$mod=$ni1%10  ' $indir/mixed-types.xtab
run_mlr --xtab put -F '$quot=$pf1/10.0;$iquot=$pf1//10.0;$mod=$pf1%10.0' $indir/mixed-types.xtab
run_mlr --xtab put -F '$quot=$pi1/10  ;$iquot=$pi1//10  ;$mod=$pi1%10  ' $indir/mixed-types.xtab
run_mlr --xtab put -F '$quot=$nf1/10.0;$iquot=$nf1//10.0;$mod=$nf1%10.0' $indir/mixed-types.xtab
run_mlr --xtab put -F '$quot=$ni1/10  ;$iquot=$ni1//10  ;$mod=$ni1%10  ' $indir/mixed-types.xtab

run_mlr --xtab put    '$quot=$pf1/-10.0;$iquot=$pf1//-10.0;$mod=$pf1%-10.0' $indir/mixed-types.xtab
run_mlr --xtab put    '$quot=$pi1/-10  ;$iquot=$pi1//-10  ;$mod=$pi1%-10  ' $indir/mixed-types.xtab
run_mlr --xtab put    '$quot=$nf1/-10.0;$iquot=$nf1//-10.0;$mod=$nf1%-10.0' $indir/mixed-types.xtab
run_mlr --xtab put    '$quot=$ni1/-10  ;$iquot=$ni1//-10  ;$mod=$ni1%-10  ' $indir/mixed-types.xtab
run_mlr --xtab put -F '$quot=$pf1/-10.0;$iquot=$pf1//-10.0;$mod=$pf1%-10.0' $indir/mixed-types.xtab
run_mlr --xtab put -F '$quot=$pi1/-10  ;$iquot=$pi1//-10  ;$mod=$pi1%-10  ' $indir/mixed-types.xtab
run_mlr --xtab put -F '$quot=$nf1/-10.0;$iquot=$nf1//-10.0;$mod=$nf1%-10.0' $indir/mixed-types.xtab
run_mlr --xtab put -F '$quot=$ni1/-10  ;$iquot=$ni1//-10  ;$mod=$ni1%-10  ' $indir/mixed-types.xtab
