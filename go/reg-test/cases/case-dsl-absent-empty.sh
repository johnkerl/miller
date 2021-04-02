run_mlr put -q '@sum     += $x; end{emitp @sum}'      $indir/abixy
run_mlr put -q '@sum[$a] += $x; end{emitp @sum, "a"}' $indir/abixy
run_mlr put    '$nonesuch = @nonesuch' $indir/abixy

run_mlr put -q '@sum     += $x; end{emitp @sum}'      $indir/abixy-het
run_mlr put -q '@sum[$a] += $x; end{emitp @sum, "a"}' $indir/abixy-het
run_mlr put    '$nonesuch = @nonesuch' $indir/abixy-het

run_mlr put -q '@sum += $x; @sumtype = typeof(@sum); @xtype = typeof($x); emitf @sumtype, @xtype, @sum; end{emitp @sum}' $indir/abixy
run_mlr put -q '@sum += $x; @sumtype = typeof(@sum); @xtype = typeof($x); emitf @sumtype, @xtype, @sum; end{emitp @sum}' $indir/abixy-het

run_mlr put '$z = $x + $y' $indir/typeof.dkvp
run_mlr put '$z = $x + $u' $indir/typeof.dkvp

run_mlr put '@s = @s + $y; emitp @s' $indir/typeof.dkvp
run_mlr put '@s = @s + $u; emitp @s' $indir/typeof.dkvp

run_mlr put '$z = $x + $y; $x=typeof($x);$y=typeof($y);$z=typeof($z)' $indir/typeof.dkvp
run_mlr put '$z = $x + $u; $x=typeof($x);$y=typeof($y);$z=typeof($z)' $indir/typeof.dkvp

run_mlr put '@s = @s + $y; $x=typeof($x);$y=typeof($y);$z=typeof($z);$s=typeof(@s)' $indir/typeof.dkvp
run_mlr put '@s = @s + $u; $x=typeof($x);$y=typeof($y);$z=typeof($z);$s=typeof(@s)' $indir/typeof.dkvp

run_mlr cat <<EOF
x=1
x=
x=7
EOF

run_mlr --ofs tab put '$osum=@sum; $ostype=typeof(@sum);$xtype=typeof($x);@sum+=$x; $nstype=typeof(@sum);$nsum=@sum; end { emit @sum }' <<EOF
x=1
x=
x=7
EOF

run_mlr --ofs tab put '$osum=@sum; $ostype=typeof(@sum);$xtype=typeof($x);is_present($x){@sum+=$x}; $nstype=typeof(@sum);$nsum=@sum; end { emit @sum }' <<EOF
x=1
x=
x=7
EOF

run_mlr cat <<EOF
x=1
xxx=
x=7
EOF

run_mlr --ofs tab put '$osum=@sum; $ostype=typeof(@sum);$xtype=typeof($x);@sum+=$x; $nstype=typeof(@sum);$nsum=@sum; end { emit @sum }' <<EOF
x=1
xxx=
x=7
EOF

run_mlr --ofs tab put '$osum=@sum; $ostype=typeof(@sum);$xtype=typeof($x);is_present($x){@sum+=$x}; $nstype=typeof(@sum);$nsum=@sum; end { emit @sum }' <<EOF
x=1
xxx=
x=7
EOF

run_mlr cat <<EOF
x=1
x=
y=
x=7
EOF

run_mlr --ofs tab put '$xtype=typeof($x);$sum = $x + 10; $stype=typeof($sum)' <<EOF
x=1
x=
y=
x=7
EOF

run_mlr --ofs tab put '$xtype=typeof($x);$sum = is_present($x) ? $x + 10 : 999; $stype=typeof($sum)' <<EOF
x=1
x=
y=
x=7
EOF
