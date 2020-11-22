run_mlr put -q '
  int n = strlen($x);
  print "input= <<".$x.">>";
  for (i = -n-2; i <= n+2; i += 1) {
    for (j = -n-2; j <= n+2; j += 1) {
      print "i: ".fmtnum(i,"%3lld")
        ."   j:".fmtnum(j,"%3lld")
        ."   substr(".$x.",".fmtnum(i,"%3lld").",".fmtnum(j,"%3lld")."): <<"
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
