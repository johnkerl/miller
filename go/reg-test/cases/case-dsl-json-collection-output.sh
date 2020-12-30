
run_mlr --json cat <<EOF
{"x":1}
EOF

run_mlr --json cat <<EOF
{"x":[1,2,3]}
EOF

run_mlr --json cat <<EOF
{"x":[1,[2,3,4],5]}
EOF

run_mlr --json cat <<EOF
{"x":[1,[2,[3,4,5],6],7]}
EOF


run_mlr --json cat <<EOF
{"x":{}}
EOF

run_mlr --json cat <<EOF
{"x":{"a":1,"b":2,"c":3}}
EOF

run_mlr --json cat <<EOF
{"x":{"a":1,"b":{"c":3,"d":4,"e":5},"f":6}}
EOF


run_mlr --json cat <<EOF
{"x":{},"y":1}
EOF

run_mlr --json cat <<EOF
{"x":{"a":1,"b":2,"c":3},"y":4}
EOF

run_mlr --json cat <<EOF
{"x":{"a":1,"b":{"c":3,"d":4,"e":5},"f":6},"y":7}
EOF


$path_to_mlr --json cat <<EOF | run_mlr --json cat
{"x":1}
EOF

$path_to_mlr --json cat <<EOF | run_mlr --json cat
{"x":[1,2,3]}
EOF

$path_to_mlr --json cat <<EOF | run_mlr --json cat
{"x":[1,[2,3,4],5]}
EOF

$path_to_mlr --json cat <<EOF | run_mlr --json cat
{"x":[1,[2,[3,4,5],6],7]}
EOF


$path_to_mlr --json cat <<EOF | run_mlr --json cat
{"x":{}}
EOF

$path_to_mlr --json cat <<EOF | run_mlr --json cat
{"x":{"a":1,"b":2,"c":3}}
EOF

$path_to_mlr --json cat <<EOF | run_mlr --json cat
{"x":{"a":1,"b":{"c":3,"d":4,"e":5},"f":6}}
EOF


$path_to_mlr --json cat <<EOF | run_mlr --json cat
{"x":{},"y":1}
EOF

$path_to_mlr --json cat <<EOF | run_mlr --json cat
{"x":{"a":1,"b":2,"c":3},"y":4}
EOF

$path_to_mlr --json cat <<EOF | run_mlr --json cat
{"x":{"a":1,"b":{"c":3,"d":4,"e":5},"f":6},"y":7}
EOF



$path_to_mlr --json cat <<EOF | run_mlr --json --jvstack cat
{"x":[1,[2,[3,4,5],6],7]}
EOF
$path_to_mlr --json cat <<EOF | run_mlr --json --no-jvstack cat
{"x":[1,[2,[3,4,5],6],7]}
EOF

$path_to_mlr --json cat <<EOF | run_mlr --json --jvstack cat
{"x":{"a":1,"b":{"c":3,"d":4,"e":5},"f":6},"y":7}
EOF
$path_to_mlr --json cat <<EOF | run_mlr --json --no-jvstack cat
{"x":{"a":1,"b":{"c":3,"d":4,"e":5},"f":6},"y":7}
EOF
