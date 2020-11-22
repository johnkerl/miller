run_mlr --oxtab --idkvp --irs lf   --ifs '\x2c'  --ips '\075'  cut -o -f x,a,i $indir/multi-sep.dkvp-crlf
run_mlr --oxtab --idkvp --irs lf   --ifs /, --ips '\x3d\x3a' cut -o -f x,a,i $indir/multi-sep.dkvp-crlf
