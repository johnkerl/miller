run_mlr                --opprint stats1 -a sum -f x $indir/ofmt.dat
run_mlr --ofmt '%.3lf' --opprint stats1 -a sum -f x $indir/ofmt.dat
run_mlr --opprint --ofmt '%.3lf' stats1 -a sum -f x $indir/ofmt.dat
