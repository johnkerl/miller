run_mlr put '$y=ssub($x, "e", "X")' <<EOF
x=eeee
EOF

run_mlr put '$y=sub($x, "e", "X")' <<EOF
x=eeee
EOF

run_mlr put '$y=gsub($x, "e", "X")' <<EOF
x=eeee
EOF
