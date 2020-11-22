run_mlr --prepipe cat --odkvp join -j a -f $indir/join-het.dkvp $indir/abixy-het
run_mlr --odkvp join --prepipe cat -j a -f $indir/join-het.dkvp $indir/abixy-het
run_mlr --prepipe cat --odkvp join --prepipe cat -j a -f $indir/join-het.dkvp $indir/abixy-het
