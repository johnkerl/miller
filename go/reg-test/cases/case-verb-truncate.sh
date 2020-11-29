
for n in 0 1 2 3 4 5; do
    run_mlr put '$y=truncate($x, '$n')' <<EOF
x=abcd
EOF
done

for n in 0 1 2 3 4 5; do
    run_mlr put '$y=truncate($x, '$n')' <<EOF
x=abcdefgh
EOF
done
