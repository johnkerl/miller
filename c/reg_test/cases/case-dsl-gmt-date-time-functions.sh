
run_mlr --csvlite put '$gmt = sec2gmt($sec)' $indir/sec2gmt
run_mlr --csvlite put '$gmt = sec2gmt($sec,1)' $indir/sec2gmt
run_mlr --csvlite put '$gmt = sec2gmt($sec,3)' $indir/sec2gmt
run_mlr --csvlite put '$gmt = sec2gmt($sec,6)' $indir/sec2gmt
run_mlr --csvlite put '$sec = gmt2sec($gmt)' $indir/gmt2sec
run_mlr --csvlite put '$gmtdate = sec2gmtdate($sec)' $indir/sec2gmt

run_mlr --icsv --opprint put '$gmt = strftime($sec, "%Y-%m-%dT%H:%M:%SZ")'  $indir/sec2gmt
run_mlr --icsv --opprint put '$gmt = strftime($sec, "%Y-%m-%dT%H:%M:%1SZ")' $indir/sec2gmt
run_mlr --icsv --opprint put '$gmt = strftime($sec, "%Y-%m-%dT%H:%M:%3SZ")' $indir/sec2gmt
run_mlr --icsv --opprint put '$gmt = strftime($sec, "%Y-%m-%dT%H:%M:%6SZ")' $indir/sec2gmt
run_mlr --icsv --opprint put '$sec = strptime($gmt, "%Y-%m-%dT%H:%M:%SZ")'  $indir/gmt2sec

run_mlr --csvlite sec2gmt sec $indir/sec2gmt

run_mlr --opprint put '$hms=sec2hms($sec);   $resec=hms2sec($hms);   $diff=$resec-$sec' $indir/sec2xhms
run_mlr --opprint put '$hms=fsec2hms($sec);  $resec=hms2fsec($hms);  $diff=$resec-$sec' $indir/fsec2xhms
run_mlr --opprint put '$hms=sec2dhms($sec);  $resec=dhms2sec($hms);  $diff=$resec-$sec' $indir/sec2xhms
run_mlr --opprint put '$hms=fsec2dhms($sec); $resec=dhms2fsec($hms); $diff=$resec-$sec' $indir/fsec2xhms

run_mlr --csvlite sec2gmt     sec $indir/sec2gmt
run_mlr --csvlite sec2gmtdate sec $indir/sec2gmt
