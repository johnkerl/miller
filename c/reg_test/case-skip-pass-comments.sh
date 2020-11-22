
# ----------------------------------------------------------------
announce SKIP/PASS COMMENTS DKVP

mention input comments1.dkvp
run_cat $indir/comments/comments1.dkvp

mention skip comments1.dkvp
run_mlr --skip-comments --idkvp --oxtab cat < $indir/comments/comments1.dkvp
run_mlr --skip-comments --idkvp --oxtab cat   $indir/comments/comments1.dkvp

mention pass comments1.dkvp
run_mlr --pass-comments --idkvp --oxtab cat < $indir/comments/comments1.dkvp
run_mlr --pass-comments --idkvp --oxtab cat   $indir/comments/comments1.dkvp


mention input comments2.dkvp
run_cat $indir/comments/comments2.dkvp

mention skip comments2.dkvp
run_mlr --skip-comments --idkvp --oxtab cat < $indir/comments/comments2.dkvp
run_mlr --skip-comments --idkvp --oxtab cat   $indir/comments/comments2.dkvp

mention pass comments2.dkvp
run_mlr --pass-comments --idkvp --oxtab cat < $indir/comments/comments2.dkvp
run_mlr --pass-comments --idkvp --oxtab cat   $indir/comments/comments2.dkvp


mention input comments3.dkvp
run_cat $indir/comments/comments3.dkvp

mention skip comments3.dkvp
run_mlr --skip-comments --idkvp --oxtab cat < $indir/comments/comments3.dkvp
run_mlr --skip-comments --idkvp --oxtab cat   $indir/comments/comments3.dkvp

mention pass comments3.dkvp
run_mlr --pass-comments --idkvp --oxtab cat < $indir/comments/comments3.dkvp
run_mlr --pass-comments --idkvp --oxtab cat   $indir/comments/comments3.dkvp

# It's annoying trying to check in text files (especially CSV) with CRLF
# to Git, given that it likes to 'fix' line endings for multi-platform use.
# It's easy to simply create CRLF on the fly.
run_mlr_for_auxents_no_output termcvt --lf2crlf < $indir/comments/comments1.dkvp > $outdir/comments1-crlf.dkvp

mention input comments1-crlf.dkvp
run_cat $outdir/comments1-crlf.dkvp

mention skip comments1-crlf.dkvp
run_mlr --skip-comments --idkvp --oxtab cat < $outdir/comments1-crlf.dkvp
run_mlr --skip-comments --idkvp --oxtab cat   $outdir/comments1-crlf.dkvp

mention pass comments1-crlf.dkvp
run_mlr --pass-comments --idkvp --oxtab cat < $outdir/comments1-crlf.dkvp
run_mlr --pass-comments --idkvp --oxtab cat   $outdir/comments1-crlf.dkvp

# - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -
announce SKIP/PASS COMMENTS DKVP WITH ALTERNATE PREFIX

mention input comments1-atat.dkvp
run_cat $indir/comments/comments1-atat.dkvp

mention skip comments1-atat.dkvp
run_mlr --skip-comments-with @@ --idkvp --oxtab cat < $indir/comments/comments1-atat.dkvp
run_mlr --skip-comments-with @@ --idkvp --oxtab cat   $indir/comments/comments1-atat.dkvp

mention pass comments1-atat.dkvp
run_mlr --pass-comments-with @@ --idkvp --oxtab cat < $indir/comments/comments1-atat.dkvp
run_mlr --pass-comments-with @@ --idkvp --oxtab cat   $indir/comments/comments1-atat.dkvp


mention input comments2-atat.dkvp
run_cat $indir/comments/comments2-atat.dkvp

mention skip comments2-atat.dkvp
run_mlr --skip-comments-with @@ --idkvp --oxtab cat < $indir/comments/comments2-atat.dkvp
run_mlr --skip-comments-with @@ --idkvp --oxtab cat   $indir/comments/comments2-atat.dkvp

mention pass comments2-atat.dkvp
run_mlr --pass-comments-with @@ --idkvp --oxtab cat < $indir/comments/comments2-atat.dkvp
run_mlr --pass-comments-with @@ --idkvp --oxtab cat   $indir/comments/comments2-atat.dkvp


mention input comments3-atat.dkvp
run_cat $indir/comments/comments3-atat.dkvp

mention skip comments3-atat.dkvp
run_mlr --skip-comments-with @@ --idkvp --oxtab cat < $indir/comments/comments3-atat.dkvp
run_mlr --skip-comments-with @@ --idkvp --oxtab cat   $indir/comments/comments3-atat.dkvp

mention pass comments3-atat.dkvp
run_mlr --pass-comments-with @@ --idkvp --oxtab cat < $indir/comments/comments3-atat.dkvp
run_mlr --pass-comments-with @@ --idkvp --oxtab cat   $indir/comments/comments3-atat.dkvp

# ----------------------------------------------------------------
announce SKIP/PASS COMMENTS NIDX

mention input comments1.nidx
run_cat $indir/comments/comments1.nidx

mention skip comments1.nidx
run_mlr --skip-comments --inidx --oxtab cat < $indir/comments/comments1.nidx
run_mlr --skip-comments --inidx --oxtab cat   $indir/comments/comments1.nidx

mention pass comments1.nidx
run_mlr --pass-comments --inidx --oxtab cat < $indir/comments/comments1.nidx
run_mlr --pass-comments --inidx --oxtab cat   $indir/comments/comments1.nidx


mention input comments2.nidx
run_cat $indir/comments/comments2.nidx

mention skip comments2.nidx
run_mlr --skip-comments --inidx --oxtab cat < $indir/comments/comments2.nidx
run_mlr --skip-comments --inidx --oxtab cat   $indir/comments/comments2.nidx

mention pass comments2.nidx
run_mlr --pass-comments --inidx --oxtab cat < $indir/comments/comments2.nidx
run_mlr --pass-comments --inidx --oxtab cat   $indir/comments/comments2.nidx


mention input comments3.nidx
run_cat $indir/comments/comments3.nidx

mention skip comments3.nidx
run_mlr --skip-comments --inidx --oxtab cat < $indir/comments/comments3.nidx
run_mlr --skip-comments --inidx --oxtab cat   $indir/comments/comments3.nidx

mention pass comments3.nidx
run_mlr --pass-comments --inidx --oxtab cat < $indir/comments/comments3.nidx
run_mlr --pass-comments --inidx --oxtab cat   $indir/comments/comments3.nidx


# It's annoying trying to check in text files (especially CSV) with CRLF
# to Git, given that it likes to 'fix' line endings for multi-platform use.
# It's easy to simply create CRLF on the fly.
run_mlr_for_auxents_no_output termcvt --lf2crlf < $indir/comments/comments1.nidx > $outdir/comments1-crlf.nidx

mention input comments1-crlf.nidx
run_cat $outdir/comments1-crlf.nidx

mention skip comments1-crlf.nidx
run_mlr --skip-comments --inidx --oxtab cat < $outdir/comments1-crlf.nidx
run_mlr --skip-comments --inidx --oxtab cat   $outdir/comments1-crlf.nidx

mention pass comments1-crlf.nidx
run_mlr --pass-comments --inidx --oxtab cat < $outdir/comments1-crlf.nidx
run_mlr --pass-comments --inidx --oxtab cat   $outdir/comments1-crlf.nidx

# ----------------------------------------------------------------
announce SKIP/PASS COMMENTS JSON

mention input comments1.json
run_cat $indir/comments/comments1.json

mention skip comments1.json
run_mlr --skip-comments --ijson --odkvp cat < $indir/comments/comments1.json
run_mlr --skip-comments --ijson --odkvp cat   $indir/comments/comments1.json

mention pass comments1.json
run_mlr --pass-comments --ijson --odkvp cat < $indir/comments/comments1.json
run_mlr --pass-comments --ijson --odkvp cat   $indir/comments/comments1.json


mention input comments2.json
run_cat $indir/comments/comments2.json

mention skip comments2.json
run_mlr --skip-comments --ijson --odkvp cat < $indir/comments/comments2.json
run_mlr --skip-comments --ijson --odkvp cat   $indir/comments/comments2.json

mention pass comments2.json
run_mlr --pass-comments --ijson --odkvp cat < $indir/comments/comments2.json
run_mlr --pass-comments --ijson --odkvp cat   $indir/comments/comments2.json


mention input comments3.json
run_cat $indir/comments/comments3.json

mention skip comments3.json
run_mlr --skip-comments --ijson --odkvp cat < $indir/comments/comments3.json
run_mlr --skip-comments --ijson --odkvp cat   $indir/comments/comments3.json

mention pass comments3.json
run_mlr --pass-comments --ijson --odkvp cat < $indir/comments/comments3.json
run_mlr --pass-comments --ijson --odkvp cat   $indir/comments/comments3.json


# It's annoying trying to check in text files (especially CSV) with CRLF
# to Git, given that it likes to 'fix' line endings for multi-platform use.
# It's easy to simply create CRLF on the fly.
run_mlr_for_auxents_no_output termcvt --lf2crlf < $indir/comments/comments1.json > $outdir/comments1-crlf.json

mention input comments1-crlf.json
run_cat $outdir/comments1-crlf.json

mention skip comments1-crlf.json
run_mlr --skip-comments --ijson --odkvp cat < $outdir/comments1-crlf.json
run_mlr --skip-comments --ijson --odkvp cat   $outdir/comments1-crlf.json

mention pass comments1-crlf.json
run_mlr --pass-comments --ijson --odkvp cat < $outdir/comments1-crlf.json
run_mlr --pass-comments --ijson --odkvp cat   $outdir/comments1-crlf.json

# ----------------------------------------------------------------
announce SKIP/PASS COMMENTS XTAB

mention input comments1.xtab
run_cat $indir/comments/comments1.xtab

mention skip comments1.xtab
run_mlr --skip-comments --ixtab --odkvp cat < $indir/comments/comments1.xtab
run_mlr --skip-comments --ixtab --odkvp cat   $indir/comments/comments1.xtab

mention pass comments1.xtab
run_mlr --pass-comments --ixtab --odkvp cat < $indir/comments/comments1.xtab
run_mlr --pass-comments --ixtab --odkvp cat   $indir/comments/comments1.xtab


mention input comments2.xtab
run_cat $indir/comments/comments2.xtab

mention skip comments2.xtab
run_mlr --skip-comments --ixtab --odkvp cat < $indir/comments/comments2.xtab
run_mlr --skip-comments --ixtab --odkvp cat   $indir/comments/comments2.xtab

mention pass comments2.xtab
run_mlr --pass-comments --ixtab --odkvp cat < $indir/comments/comments2.xtab
run_mlr --pass-comments --ixtab --odkvp cat   $indir/comments/comments2.xtab


# It's annoying trying to check in text files (especially CSV) with CRLF
# to Git, given that it likes to 'fix' line endings for multi-platform use.
# It's easy to simply create CRLF on the fly.
run_mlr_for_auxents_no_output termcvt --lf2crlf < $indir/comments/comments1.xtab > $outdir/comments1-crlf.xtab

mention input comments1-crlf.xtab
run_cat $outdir/comments1-crlf.xtab

mention skip comments1-crlf.xtab
run_mlr --skip-comments --ixtab --odkvp cat < $outdir/comments1-crlf.xtab
run_mlr --skip-comments --ixtab --odkvp cat   $outdir/comments1-crlf.xtab

mention pass comments1-crlf.xtab
run_mlr --pass-comments --ixtab --odkvp cat < $outdir/comments1-crlf.xtab
run_mlr --pass-comments --ixtab --odkvp cat   $outdir/comments1-crlf.xtab

# ----------------------------------------------------------------
announce SKIP/PASS COMMENTS CSVLITE

mention input comments1.csv
run_cat $indir/comments/comments1.csv

mention skip comments1.csv
run_mlr --skip-comments --icsvlite --odkvp cat < $indir/comments/comments1.csv
run_mlr --skip-comments --icsvlite --odkvp cat   $indir/comments/comments1.csv

mention pass comments1.csv
run_mlr --pass-comments --icsvlite --odkvp cat < $indir/comments/comments1.csv
run_mlr --pass-comments --icsvlite --odkvp cat   $indir/comments/comments1.csv


mention input comments2.csv
run_cat $indir/comments/comments2.csv

mention skip comments2.csv
run_mlr --skip-comments --icsvlite --odkvp cat < $indir/comments/comments2.csv
run_mlr --skip-comments --icsvlite --odkvp cat   $indir/comments/comments2.csv

mention pass comments2.csv
run_mlr --pass-comments --icsvlite --odkvp cat < $indir/comments/comments2.csv
run_mlr --pass-comments --icsvlite --odkvp cat   $indir/comments/comments2.csv


# It's annoying trying to check in text files (especially CSV) with CRLF
# to Git, given that it likes to 'fix' line endings for multi-platform use.
# It's easy to simply create CRLF on the fly.
run_mlr_for_auxents_no_output termcvt --lf2crlf < $indir/comments/comments1.csv > $outdir/comments1-crlf.csv

mention input comments1-crlf.csv
run_cat $outdir/comments1-crlf.csv

mention skip comments1-crlf.csv
run_mlr --skip-comments --icsvlite --odkvp cat < $outdir/comments1-crlf.csv
run_mlr --skip-comments --icsvlite --odkvp cat   $outdir/comments1-crlf.csv

mention pass comments1-crlf.csv
run_mlr --pass-comments --icsvlite --odkvp cat < $outdir/comments1-crlf.csv
run_mlr --pass-comments --icsvlite --odkvp cat   $outdir/comments1-crlf.csv

# ----------------------------------------------------------------
announce SKIP/PASS COMMENTS CSV

mention input comments1.csv
run_cat $indir/comments/comments1.csv

mention skip comments1.csv
run_mlr --skip-comments --icsv --odkvp cat < $indir/comments/comments1.csv
run_mlr --skip-comments --icsv --odkvp cat   $indir/comments/comments1.csv

mention pass comments1.csv
run_mlr --pass-comments --icsv --odkvp cat < $indir/comments/comments1.csv
run_mlr --pass-comments --icsv --odkvp cat   $indir/comments/comments1.csv


mention input comments2.csv
run_cat $indir/comments/comments2.csv

mention skip comments2.csv
run_mlr --skip-comments --icsv --odkvp cat < $indir/comments/comments2.csv
run_mlr --skip-comments --icsv --odkvp cat   $indir/comments/comments2.csv

mention pass comments2.csv
run_mlr --pass-comments --icsv --odkvp cat < $indir/comments/comments2.csv
run_mlr --pass-comments --icsv --odkvp cat   $indir/comments/comments2.csv


# It's annoying trying to check in text files (especially CSV) with CRLF
# to Git, given that it likes to 'fix' line endings for multi-platform use.
# It's easy to simply create CRLF on the fly.
run_mlr_for_auxents_no_output termcvt --lf2crlf < $indir/comments/comments1.csv > $outdir/comments1-crlf.csv

mention input comments1-crlf.csv
run_cat $outdir/comments1-crlf.csv

mention skip comments1-crlf.csv
run_mlr --skip-comments --icsv --odkvp cat < $outdir/comments1-crlf.csv
run_mlr --skip-comments --icsv --odkvp cat   $outdir/comments1-crlf.csv

mention pass comments1-crlf.csv
run_mlr --pass-comments --icsv --odkvp cat < $outdir/comments1-crlf.csv
run_mlr --pass-comments --icsv --odkvp cat   $outdir/comments1-crlf.csv

# ----------------------------------------------------------------
announce INT64 I/O

run_mlr --oxtab put '$c=$a;$d=$b;$e=hexfmt($a);$f=hexfmt($b)' $indir/int64io.dkvp

# There is different rounding on i386 vs. x86_64 but the essential check is
# that overflow has been avoided. Hence the %g output format here.
run_mlr --oxtab --ofmt '%.8g' put '$p0=$p+0;$p1=$p+1;$p2=$p+2;$p3=$p+3' $indir/int64arith.dkvp
run_mlr --oxtab --ofmt '%.8g' put '$p0=$p-0;$p1=$p-1;$p2=$p-2;$p3=$p-3' $indir/int64arith.dkvp
run_mlr --oxtab --ofmt '%.8g' put '$p0=$p*0;$p1=$p*1;$p2=$p*2;$p3=$p*3' $indir/int64arith.dkvp
run_mlr --oxtab --ofmt '%.8g' put '$n0=$n+0;$n1=$n+1;$n2=$n+2;$n3=$n+3' $indir/int64arith.dkvp
run_mlr --oxtab --ofmt '%.8g' put '$n0=$n-0;$n1=$n-1;$n2=$n-2;$n3=$n-3' $indir/int64arith.dkvp
run_mlr --oxtab --ofmt '%.8g' put '$n0=$n*0;$n1=$n*1;$n2=$n*2;$n3=$n*3' $indir/int64arith.dkvp
