run_mlr format-values $indir/abixy
run_mlr format-values -n $indir/abixy
run_mlr format-values -i %08llx -f %.6le -s X%sX $indir/abixy
run_mlr format-values -i %08llx -f %.6le -s X%sX -n $indir/abixy
