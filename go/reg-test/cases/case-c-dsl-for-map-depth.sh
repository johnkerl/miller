mention 'for full oosvar'
run_mlr --from $indir/abixy put '@o[1][2] = 7; for(k1,v in @*) {$x+=10;$y+=100}'
run_mlr --from $indir/abixy put '@o[1][2] = 7; for((k1,k2),v in @*) {$x+=10;$y+=100}'
run_mlr --from $indir/abixy put '@o[1][2] = 7; for((k1,k2,k3),v in @*) {$x+=10;$y+=100}'

mention 'for oosvar submap'
run_mlr --from $indir/abixy put '@o[1][2][3] = 7; for(k1,v in @o[1][2]) {$x+=10;$y+=100}'
run_mlr --from $indir/abixy put '@o[1][2][3] = 7; for((k1,k2),v in @o[1][2]) {$x+=10;$y+=100}'
run_mlr --from $indir/abixy put '@o[1][2][3] = 7; for((k1,k2,k3),v in @o[1][2]) {$x+=10;$y+=100}'

mention 'for local'
run_mlr --from $indir/abixy put 'o[1][2] = 7; for(k1,v in o) {$x+=10;$y+=100}'
run_mlr --from $indir/abixy put 'o[1][2] = 7; for((k1,k2),v in o) {$x+=10;$y+=100}'
run_mlr --from $indir/abixy put 'o[1][2] = 7; for((k1,k2,k3),v in o) {$x+=10;$y+=100}'

mention 'for map-literal'
run_mlr --from $indir/abixy put 'for(k1,v in {1:{2:7}}) {$x+=10;$y+=100}'
run_mlr --from $indir/abixy put 'for((k1,k2),v in {1:{2:7}}) {$x+=10;$y+=100}'
run_mlr --from $indir/abixy put 'for((k1,k2,k3),v in {1:{2:7}}) {$x+=10;$y+=100}'
