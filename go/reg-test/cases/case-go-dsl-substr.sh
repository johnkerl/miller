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
