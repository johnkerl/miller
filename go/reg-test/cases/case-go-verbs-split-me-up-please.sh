
run_mlr --opprint cat           $indir/s.dkvp
run_mlr --opprint cat -n        $indir/s.dkvp
run_mlr --opprint cat -n -g a   $indir/s.dkvp
run_mlr --opprint cat -n -g a,b $indir/s.dkvp

run_mlr --opprint cut    -f x,a $indir/s.dkvp
run_mlr --opprint cut -o -f x,a $indir/s.dkvp
run_mlr --opprint cut -x -f x,a $indir/s.dkvp


run_mlr --opprint head -n 1 $indir/s.dkvp
run_mlr --opprint head -n 1 -g a $indir/s.dkvp
run_mlr --opprint head -n 1 -g a,b $indir/s.dkvp

run_mlr --opprint tail -n 1 $indir/medium.dkvp
run_mlr --opprint tail -n 1 -g a $indir/medium.dkvp
run_mlr --opprint tail -n 1 -g a,b $indir/medium.dkvp

run_mlr --opprint group-like $indir/het.dkvp
run_mlr --opprint group-by a $indir/medium.dkvp
run_mlr --opprint group-by a,b $indir/medium.dkvp

run_mlr --opprint rename a,AAA $indir/s.dkvp
run_mlr --opprint rename a,AAA,x,XXX $indir/s.dkvp
run_mlr --opprint rename none,such $indir/s.dkvp
run_mlr --opprint rename a,b $indir/s.dkvp
run_mlr --opprint rename i,j,a,b $indir/s.dkvp
run_mlr --opprint rename x,y,a,b $indir/s.dkvp

run_mlr --opprint label A,B,I $indir/s.dkvp
run_mlr --opprint label A,B,I,X,Y $indir/s.dkvp
run_mlr --opprint label A,B,I,X,Y,Z $indir/s.dkvp
run_mlr --opprint label b,i,x $indir/s.dkvp
run_mlr --opprint label x,i,b $indir/s.dkvp

run_mlr --opprint --from $indir/s.dkvp sort -f nonesuch
run_mlr --opprint --from $indir/s.dkvp sort -f a
run_mlr --opprint --from $indir/s.dkvp sort -f a,b
run_mlr --opprint --from $indir/s.dkvp sort -r a
run_mlr --opprint --from $indir/s.dkvp sort -r a,b
run_mlr --opprint --from $indir/s.dkvp sort -f a -r b
run_mlr --opprint --from $indir/s.dkvp sort -f b -n i
run_mlr --opprint --from $indir/s.dkvp sort -f b -nr i

run_mlr --json --from $indir/needs-sorting.json sort-within-records
run_mlr --json --from $indir/needs-regularize.json regularize

run_mlr --opprint --from $indir/needs-unsparsify.dkvp unsparsify
run_mlr --opprint --from $indir/needs-unsparsify.dkvp unsparsify --fill-with X
run_mlr --opprint --from $indir/abixy-het unsparsify
run_mlr --opprint --from $indir/abixy-het unsparsify -f a,b,i,x,y
run_mlr --opprint --from $indir/abixy-het unsparsify -f a,b,i,x,y then regularize
run_mlr --opprint --from $indir/abixy-het unsparsify -f aaa,bbb,iii,xxx,yyy
run_mlr --opprint --from $indir/abixy-het unsparsify -f aaa,bbb,iii,xxx,yyy then regularize
run_mlr --opprint --from $indir/abixy-het unsparsify -f a,b,i,x,y,aaa,bbb,iii,xxx,yyy
run_mlr --opprint --from $indir/abixy-het unsparsify -f a,b,i,x,y,aaa,bbb,iii,xxx,yyy then regularize

run_mlr --opprint --from $indir/medium.dkvp count
run_mlr --opprint --from $indir/medium.dkvp count -g a
run_mlr --opprint --from $indir/medium.dkvp count -g a,b
run_mlr --opprint --from $indir/medium.dkvp count -g a -n
run_mlr --opprint --from $indir/medium.dkvp count -g a,b -n
run_mlr --opprint --from $indir/medium.dkvp count -o NAME
run_mlr --opprint --from $indir/medium.dkvp count -g a -o NAME

run_mlr --opprint --from $indir/s.dkvp grep pan
run_mlr --opprint --from $indir/s.dkvp grep -v pan
run_mlr --opprint --from $indir/s.dkvp grep PAN
run_mlr --opprint --from $indir/s.dkvp grep -i PAN
run_mlr --opprint --from $indir/s.dkvp grep -i -v PAN

run_mlr --from $indir/s.dkvp skip-trivial-records
run_mlr --from $indir/skip-trivial-records.dkvp skip-trivial-records

echo 'a,b,c,d,e,f'   | run_mlr --inidx --ifs comma altkv
echo 'a,b,c,d,e,f,g' | run_mlr --inidx --ifs comma altkv

run_mlr --from $indir/s.csv                    --icsv --opprint remove-empty-columns
run_mlr --from $indir/remove-empty-columns.csv --icsv --opprint cat
run_mlr --from $indir/remove-empty-columns.csv --icsv --opprint remove-empty-columns

run_mlr --icsv --opprint fill-down -f z       $indir/remove-empty-columns.csv
run_mlr --icsv --opprint fill-down -f a,b,c,d $indir/remove-empty-columns.csv

run_mlr --icsv --opprint reorder -f x,i    $indir/s.dkvp
run_mlr --icsv --opprint reorder -f x,i -e $indir/s.dkvp

run_mlr decimate       -n 2 $indir/s.dkvp
run_mlr decimate -b    -n 2 $indir/s.dkvp
run_mlr decimate    -e -n 2 $indir/s.dkvp
run_mlr decimate -b -e -n 2 $indir/s.dkvp

run_mlr --opprint count-similar -g a $indir/s.dkvp
run_mlr --opprint count-similar -g b $indir/s.dkvp
run_mlr --opprint count-similar -g a,b $indir/s.dkvp
run_mlr --opprint count-similar -g a -o altnamehere $indir/s.dkvp

run_mlr --from $indir/s.dkvp --opprint put '$t = $i + 0.123456789' then sec2gmt a,t
run_mlr --from $indir/s.dkvp --opprint put '$t = $i + 0.123456789' then sec2gmt -1 a,t
run_mlr --from $indir/s.dkvp --opprint put '$t = $i + 0.123456789' then sec2gmt -2 a,t
run_mlr --from $indir/s.dkvp --opprint put '$t = $i + 0.123456789' then sec2gmt -3 a,t
run_mlr --from $indir/s.dkvp --opprint put '$t = $i + 0.123456789' then sec2gmt -4 a,t
run_mlr --from $indir/s.dkvp --opprint put '$t = $i + 0.123456789' then sec2gmt -5 a,t
run_mlr --from $indir/s.dkvp --opprint put '$t = $i + 0.123456789' then sec2gmt -6 a,t
run_mlr --from $indir/s.dkvp --opprint put '$t = $i + 0.123456789' then sec2gmt -7 a,t
run_mlr --from $indir/s.dkvp --opprint put '$t = $i + 0.123456789' then sec2gmt -8 a,t
run_mlr --from $indir/s.dkvp --opprint put '$t = $i + 0.123456789' then sec2gmt -9 a,t

mlr_expect_fail --from $indir/ten.dkvp gap
run_mlr --from $indir/ten.dkvp gap -n 4
run_mlr --from $indir/ten.dkvp gap -g a
run_mlr --from $indir/ten.dkvp sort -f a then gap -g a
run_mlr --from $indir/ten.dkvp sort -f a,b then gap -g a,b
