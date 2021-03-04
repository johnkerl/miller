run_mlr put '$y = substr($x, 0, 0)' <<EOF
x=abcdefg
EOF

run_mlr put '$y = substr($x, 0, 7)' <<EOF
x=abcdefg
EOF

run_mlr put '$y = substr($x, 1, 7)' <<EOF
x=abcdefg
EOF

run_mlr put '$y = substr($x, 1, 6)' <<EOF
x=abcdefg
EOF

run_mlr put '$y = substr($x, 2, 5)' <<EOF
x=abcdefg
EOF

run_mlr put '$y = substr($x, 2, 3)' <<EOF
x=abcdefg
EOF

run_mlr put '$y = substr($x, 3, 3)' <<EOF
x=abcdefg
EOF

run_mlr put '$y = substr($x, 4, 3)' <<EOF
x=abcdefg
EOF

run_mlr put '$y = substr($x, 2, 5)' <<EOF
x=1234567
EOF

run_mlr put -q '
  int n = strlen($x);
  print "input= <<".$x.">>";
  for (i = -n-2; i <= n+2; i += 1) {
    for (j = -n-2; j <= n+2; j += 1) {
      print "i: ".fmtnum(i,"%3d")
        ."   j:".fmtnum(j,"%3d")
        ."   substr(".$x.",".fmtnum(i,"%3d").",".fmtnum(j,"%3d")."): <<"
        .substr($x, i, j) .">>";
    }
    print;
}
' << EOF
x=
x=o
x=o1
x=o123456789
EOF
