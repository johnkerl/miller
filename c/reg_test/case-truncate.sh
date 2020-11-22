
for n in 0 1 2 3 4 5; do
    echo "x=abcd" | run_mlr put '$y=truncate($x, '$n')'
done
for n in 0 1 2 3 4 5; do
    echo "x=abcdefgh" | run_mlr put '$y=truncate($x, '$n')'
done
