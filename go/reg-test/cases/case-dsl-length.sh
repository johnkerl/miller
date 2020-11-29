run_mlr put '$n = length($x)' <<EOF
x=1,y=abcdefg,z=3
EOF

run_mlr put '$n = length($y)' <<EOF
x=1,y=abcdefg,z=3
EOF

run_mlr put '$n = length($nonesuch)' <<EOF
x=1,y=abcdefg,z=3
EOF

run_mlr put '$n = length($*)' <<EOF
x=1,y=abcdefg,z=3
EOF

run_mlr put '$n = length([])' <<EOF
x=1,y=abcdefg,z=3
EOF

run_mlr put '$n = length([5,6,7])' <<EOF
x=1,y=abcdefg,z=3
EOF

run_mlr put '$n = length({})' <<EOF
x=1,y=abcdefg,z=3
EOF

run_mlr put '$n = length({"a":5,"b":6,"c":7})' <<EOF
x=1,y=abcdefg,z=3
EOF
