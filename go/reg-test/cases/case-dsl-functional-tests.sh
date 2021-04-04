
run_mlr filter '$x>.3'    $indir/abixy
run_mlr filter '$x>.3;'   $indir/abixy
run_mlr filter '$x>0.3'   $indir/abixy
run_mlr filter '$x>0.3 && $y>0.3'   $indir/abixy
run_mlr filter '$x>0.3 || $y>0.3'   $indir/abixy
run_mlr filter 'NR>=4 && NR <= 7'   $indir/abixy

run_mlr filter -x '$x>.3'    $indir/abixy
run_mlr filter -x '$x>0.3'   $indir/abixy
run_mlr filter -x '$x>0.3 && $y>0.3'   $indir/abixy
run_mlr filter -x '$x>0.3 || $y>0.3'   $indir/abixy
run_mlr filter -x 'NR>=4 && NR <= 7'   $indir/abixy

run_mlr filter '$nosuchfield>.3'    $indir/abixy

run_mlr put '$x2 = $x**2'  $indir/abixy
run_mlr put '$x2 = $x**2;' $indir/abixy
run_mlr put '$z = -0.024*$x+0.13' $indir/abixy
run_mlr put '$c = $a . $b' $indir/abixy
run_mlr put '$ii = $i + $i' $indir/abixy
run_mlr put '$emptytest = $i + $nosuch' $indir/abixy

run_mlr --opprint put '$nr=NR;$fnr=FNR;$nf=NF;$filenum=FILENUM' $indir/abixy $indir/abixy

run_mlr --opprint put '$y=madd($x,10,37)' then put '$z=msub($x,10,37)' $indir/modarith.dat
run_mlr --opprint put '$y=mexp($x,35,37)' then put '$z=mmul($x,$y,37)' $indir/modarith.dat

run_mlr put '$z=min($x, $y)' $indir/minmax.dkvp
run_mlr put '$z=max($x, $y)' $indir/minmax.dkvp

run_mlr put '$o=min()' <<EOF
x=1,y=2,z=3
EOF

run_mlr put '$o=max()' <<EOF
x=1,y=2,z=3
EOF

run_mlr put '$o=min($x)' <<EOF
x=1,y=2,z=3
EOF

run_mlr put '$o=max($x)' <<EOF
x=1,y=2,z=3
EOF

run_mlr put '$o=min($x,$y)' <<EOF
x=1,y=2,z=3
EOF

run_mlr put '$o=max($x,$y)' <<EOF
x=1,y=2,z=3
EOF

run_mlr put '$o=min($x,$y,$z)' <<EOF
x=1,y=2,z=3
EOF

run_mlr put '$o=max($x,$y,$z)' <<EOF
x=1,y=2,z=3
EOF


run_mlr put '$u=min($x,$y);$v=max($x,$y)' <<EOF
x=1,y=b
EOF

run_mlr put '$u=min($x,$y);$v=max($x,$y)' <<EOF
x=a,y=2
EOF

run_mlr put '$u=min($x,$y);$v=max($x,$y)' <<EOF
x=a,y=b
EOF


run_mlr --icsvlite --oxtab put '${x+y} = ${name.x} + ${name.y}; ${x*y} = ${name.x} * ${name.y}' $indir/braced.csv
run_mlr --icsvlite --oxtab filter '${name.y} < ${z}' $indir/braced.csv

run_mlr --opprint put '$z = $x < 0.5 ? 0 : 1' $indir/abixy

run_mlr --csvlite filter 'true  && true'  $indir/b.csv
run_mlr --csvlite filter 'true  && false' $indir/b.csv
run_mlr --csvlite filter 'false && true'  $indir/b.csv
run_mlr --csvlite filter 'false && false' $indir/b.csv

run_mlr --csvlite filter 'true  || true'  $indir/b.csv
run_mlr --csvlite filter 'true  || false' $indir/b.csv
run_mlr --csvlite filter 'false || true'  $indir/b.csv
run_mlr --csvlite filter 'false || false' $indir/b.csv

run_mlr --csvlite filter 'true  ^^ true'  $indir/b.csv
run_mlr --csvlite filter 'true  ^^ false' $indir/b.csv
run_mlr --csvlite filter 'false ^^ true'  $indir/b.csv
run_mlr --csvlite filter 'false ^^ false' $indir/b.csv

# This tests boolean short-circuiting
run_mlr put '$x==2 && $a =~ "....." { $y=4 }'  $indir/short-circuit.dkvp

export X=97
export Y=98
run_mlr put '$x = ENV["X"]; $y = ENV[$name]' $indir/env-var.dkvp
export X=
export Y=
run_mlr put '$x = ENV["X"]; $y = ENV[$name]' $indir/env-var.dkvp

run_mlr -n put 'begin{ENV["HOME"]="foobar"} end{print ENV["HOME"]}'

run_mlr put '$y = toupper($x)' <<EOF
x=hello
EOF

run_mlr put '$y = toupper($x)' <<EOF
x=HELLO
EOF

run_mlr put '$y = toupper($x)' <<EOF
x=
EOF

run_mlr put '$y = toupper($z)' <<EOF
x=hello
EOF


run_mlr put '$y = tolower($x)' <<EOF
x=hello
EOF

run_mlr put '$y = tolower($x)' <<EOF
x=HELLO
EOF

run_mlr put '$y = tolower($x)' <<EOF
x=
EOF

run_mlr put '$y = tolower($z)' <<EOF
x=hello
EOF


run_mlr put '$y = capitalize($x)' <<EOF
x=hello
EOF

run_mlr put '$y = capitalize($x)' <<EOF
x=HELLO
EOF

run_mlr put '$y = capitalize($x)' <<EOF
x=
EOF

run_mlr put '$y = capitalize($z)' <<EOF
x=hello
EOF

# mention LHS value on first record should result in ZYX for process creation
export indir; run_mlr --from $indir/abixy put -q 'ENV["ZYX"]="CBA".NR; print | ENV["indir"]."/env-assign.sh" , "a is " . $a'
