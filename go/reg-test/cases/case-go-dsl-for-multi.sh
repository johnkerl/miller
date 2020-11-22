run_mlr -n put -v 'for (k             in @*) {}'
run_mlr -n put -v 'for (k, v          in @*) {}'
run_mlr -n put -v 'for ((k1,k2),    v in @*) {}'
run_mlr -n put -v 'for ((k1,k2,k3), v in @*) {}'
