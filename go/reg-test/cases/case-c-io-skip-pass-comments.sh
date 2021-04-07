# TODO: pending --ofmt

run_mlr --skip-comments --idkvp --oxtab cat $indir/comments/comments1.dkvp
run_mlr --pass-comments --idkvp --oxtab cat $indir/comments/comments1.dkvp
run_mlr --skip-comments --idkvp --oxtab cat $indir/comments/comments2.dkvp
run_mlr --pass-comments --idkvp --oxtab cat $indir/comments/comments2.dkvp
run_mlr --skip-comments --idkvp --oxtab cat $indir/comments/comments3.dkvp
run_mlr --pass-comments --idkvp --oxtab cat $indir/comments/comments3.dkvp
# It's annoying trying to check in text files (especially CSV) with CRLF
# to Git, given that it likes to 'fix' line endings for multi-platform use.
# It's easy to simply create CRLF on the fly.
run_mlr_no_output termcvt --lf2crlf < $indir/comments/comments1.dkvp > $outdir/comments1-crlf.dkvp
run_mlr --skip-comments --idkvp --oxtab cat $outdir/comments1-crlf.dkvp
run_mlr --pass-comments --idkvp --oxtab cat $outdir/comments1-crlf.dkvp

run_mlr --skip-comments-with @@ --idkvp --oxtab cat $indir/comments/comments1-atat.dkvp
run_mlr --pass-comments-with @@ --idkvp --oxtab cat $indir/comments/comments1-atat.dkvp
run_mlr --skip-comments-with @@ --idkvp --oxtab cat $indir/comments/comments2-atat.dkvp
run_mlr --pass-comments-with @@ --idkvp --oxtab cat $indir/comments/comments2-atat.dkvp
run_mlr --skip-comments-with @@ --idkvp --oxtab cat $indir/comments/comments3-atat.dkvp
run_mlr --pass-comments-with @@ --idkvp --oxtab cat $indir/comments/comments3-atat.dkvp

run_mlr --skip-comments --inidx --oxtab cat $indir/comments/comments1.nidx
run_mlr --pass-comments --inidx --oxtab cat $indir/comments/comments1.nidx
run_mlr --skip-comments --inidx --oxtab cat $indir/comments/comments2.nidx
run_mlr --pass-comments --inidx --oxtab cat $indir/comments/comments2.nidx
run_mlr --skip-comments --inidx --oxtab cat $indir/comments/comments3.nidx
run_mlr --pass-comments --inidx --oxtab cat $indir/comments/comments3.nidx
# It's annoying trying to check in text files (especially CSV) with CRLF
# to Git, given that it likes to 'fix' line endings for multi-platform use.
# It's easy to simply create CRLF on the fly.
run_mlr_no_output termcvt --lf2crlf < $indir/comments/comments1.nidx > $outdir/comments1-crlf.nidx
run_mlr --skip-comments --inidx --oxtab cat $outdir/comments1-crlf.nidx
run_mlr --pass-comments --inidx --oxtab cat $outdir/comments1-crlf.nidx

run_mlr --skip-comments --ijson --odkvp cat $indir/comments/comments1.json
run_mlr --pass-comments --ijson --odkvp cat $indir/comments/comments1.json
run_mlr --skip-comments --ijson --odkvp cat $indir/comments/comments2.json
run_mlr --pass-comments --ijson --odkvp cat $indir/comments/comments2.json
run_mlr --skip-comments --ijson --odkvp cat $indir/comments/comments3.json
run_mlr --pass-comments --ijson --odkvp cat $indir/comments/comments3.json
# It's annoying trying to check in text files (especially CSV) with CRLF
# to Git, given that it likes to 'fix' line endings for multi-platform use.
# It's easy to simply create CRLF on the fly.
run_mlr_no_output termcvt --lf2crlf < $indir/comments/comments1.json > $outdir/comments1-crlf.json
run_mlr --skip-comments --ijson --odkvp cat $outdir/comments1-crlf.json
run_mlr --pass-comments --ijson --odkvp cat $outdir/comments1-crlf.json

run_mlr --skip-comments --ixtab --odkvp cat $indir/comments/comments1.xtab
run_mlr --pass-comments --ixtab --odkvp cat $indir/comments/comments1.xtab
run_mlr --skip-comments --ixtab --odkvp cat $indir/comments/comments2.xtab
run_mlr --pass-comments --ixtab --odkvp cat $indir/comments/comments2.xtab
# It's annoying trying to check in text files (especially CSV) with CRLF
# to Git, given that it likes to 'fix' line endings for multi-platform use.
# It's easy to simply create CRLF on the fly.
run_mlr_no_output termcvt --lf2crlf < $indir/comments/comments1.xtab > $outdir/comments1-crlf.xtab
run_mlr --skip-comments --ixtab --odkvp cat $outdir/comments1-crlf.xtab
run_mlr --pass-comments --ixtab --odkvp cat $outdir/comments1-crlf.xtab

run_mlr --skip-comments --icsvlite --odkvp cat $indir/comments/comments1.csv
run_mlr --pass-comments --icsvlite --odkvp cat $indir/comments/comments1.csv
run_mlr --skip-comments --icsvlite --odkvp cat $indir/comments/comments2.csv
run_mlr --pass-comments --icsvlite --odkvp cat $indir/comments/comments2.csv
# It's annoying trying to check in text files (especially CSV) with CRLF
# to Git, given that it likes to 'fix' line endings for multi-platform use.
# It's easy to simply create CRLF on the fly.
run_mlr_no_output termcvt --lf2crlf < $indir/comments/comments1.csv > $outdir/comments1-crlf.csv
run_mlr --skip-comments --icsvlite --odkvp cat $outdir/comments1-crlf.csv
run_mlr --pass-comments --icsvlite --odkvp cat $outdir/comments1-crlf.csv

run_mlr --skip-comments --icsv --odkvp cat $indir/comments/comments1.csv
run_mlr --pass-comments --icsv --odkvp cat $indir/comments/comments1.csv
run_mlr --skip-comments --icsv --odkvp cat $indir/comments/comments2.csv
run_mlr --pass-comments --icsv --odkvp cat $indir/comments/comments2.csv
# It's annoying trying to check in text files (especially CSV) with CRLF
# to Git, given that it likes to 'fix' line endings for multi-platform use.
# It's easy to simply create CRLF on the fly.
run_mlr_no_output termcvt --lf2crlf < $indir/comments/comments1.csv > $outdir/comments1-crlf.csv
run_mlr --skip-comments --icsv --odkvp cat $outdir/comments1-crlf.csv
run_mlr --pass-comments --icsv --odkvp cat $outdir/comments1-crlf.csv

run_mlr --oxtab put '$c=$a;$d=$b;$e=hexfmt($a);$f=hexfmt($b)' $indir/int64io.dkvp
# There is different rounding on i386 vs. x86_64 but the essential check is
# that overflow has been avoided. Hence the %g output format here.
run_mlr --oxtab --ofmt '%.8g' put '$p0=$p+0;$p1=$p+1;$p2=$p+2;$p3=$p+3' $indir/int64arith.dkvp
run_mlr --oxtab --ofmt '%.8g' put '$p0=$p-0;$p1=$p-1;$p2=$p-2;$p3=$p-3' $indir/int64arith.dkvp
run_mlr --oxtab --ofmt '%.8g' put '$p0=$p*0;$p1=$p*1;$p2=$p*2;$p3=$p*3' $indir/int64arith.dkvp
run_mlr --oxtab --ofmt '%.8g' put '$n0=$n+0;$n1=$n+1;$n2=$n+2;$n3=$n+3' $indir/int64arith.dkvp
run_mlr --oxtab --ofmt '%.8g' put '$n0=$n-0;$n1=$n-1;$n2=$n-2;$n3=$n-3' $indir/int64arith.dkvp
run_mlr --oxtab --ofmt '%.8g' put '$n0=$n*0;$n1=$n*1;$n2=$n*2;$n3=$n*3' $indir/int64arith.dkvp
