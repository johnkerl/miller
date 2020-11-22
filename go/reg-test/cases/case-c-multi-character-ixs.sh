run_mlr --oxtab --idkvp --irs lf   --ifs '\x2c'  --ips '\075'  cut -o -f x,a,i $indir/multi-sep.dkvp-crlf
run_mlr --oxtab --idkvp --irs lf   --ifs /, --ips '\x3d\x3a' cut -o -f x,a,i $indir/multi-sep.dkvp-crlf

# ----------------------------------------------------------------
announce MULTI-CHARACTER IRS/IFS FOR NIDX

run_mlr --oxtab --inidx --irs lf   --ifs ,  cut -o -f 4,1,3 $indir/multi-sep.dkvp-crlf
run_mlr --oxtab --inidx --irs lf   --ifs /, cut -o -f 4,1,3 $indir/multi-sep.dkvp-crlf
run_mlr --oxtab --inidx --irs crlf --ifs ,  cut -o -f 4,1,3 $indir/multi-sep.dkvp-crlf
run_mlr --oxtab --inidx --irs crlf --ifs /, cut -o -f 4,1,3 $indir/multi-sep.dkvp-crlf

# ----------------------------------------------------------------
announce MULTI-CHARACTER IRS/IFS FOR CSVLITE

run_mlr --oxtab --icsvlite --irs lf   --ifs ,  cut -o -f x/,a/,i/ $indir/multi-sep.csv-crlf
run_mlr --oxtab --icsvlite --irs lf   --ifs /, cut -o -f x,a,i    $indir/multi-sep.csv-crlf
run_mlr --oxtab --icsvlite --irs crlf --ifs ,  cut -o -f x/,a/,i/ $indir/multi-sep.csv-crlf
run_mlr --oxtab --icsvlite --irs crlf --ifs /, cut -o -f x,a,i    $indir/multi-sep.csv-crlf

# ----------------------------------------------------------------
announce MULTI-CHARACTER SEPARATORS FOR XTAB

run_mlr --xtab --ifs crlf --ofs Z cut -x -f b $indir/truncated.xtab-crlf
run_mlr --xtab --ips . --ops @ cut -x -f b $indir/dots.xtab
run_mlr --xtab --ips ": " --ops '@@@@' put '$sum=int($a+$b)' $indir/multi-ips.dkvp

# ----------------------------------------------------------------
announce EMBEDDED IPS FOR XTAB

run_mlr --xtab cat $indir/embedded-ips.xtab

# ----------------------------------------------------------------
announce MULTI-CHARACTER IRS FOR PPRINT

run_mlr --pprint --irs crlf --ifs / --ofs @ cut -x -f b $indir/dots.pprint-crlf

# ----------------------------------------------------------------
announce MULTI-CHARACTER IRS/IFS/IPS FOR DKVP

run_mlr --oxtab --idkvp --irs lf   --ifs ,  --ips =  cut -o -f x,a,i $indir/multi-sep.dkvp-crlf
run_mlr --oxtab --idkvp --irs lf   --ifs /, --ips =: cut -o -f x,a,i $indir/multi-sep.dkvp-crlf
run_mlr --oxtab --idkvp --irs crlf --ifs ,  --ips =  cut -o -f x,a,i $indir/multi-sep.dkvp-crlf
run_mlr --oxtab --idkvp --irs crlf --ifs /, --ips =: cut -o -f x,a,i $indir/multi-sep.dkvp-crlf

# ----------------------------------------------------------------
announce DOUBLE PS

run_mlr --opprint cat $indir/double-ps.dkvp
