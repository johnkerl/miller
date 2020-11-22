run_mlr --from $indir/abixy-het put '$length = length($a)'
run_mlr --from $indir/abixy-het put '$length = length($*)'
run_mlr --from $indir/xyz2 put '$length= length({3:4, 5:{6:7}, 8:{9:{10:11}}})'
run_mlr --from $indir/xyz2 put 'o = {3:4, 5:{6:7}, 8:{9:{10:11}}}; $length = length(o)'
run_mlr --from $indir/xyz2 put '@o = {3:4, 5:{6:7}, 8:{9:{10:11}}}; $length = length(@o)'

run_mlr --from $indir/abixy-het put '$depth = depth($a)'
run_mlr --from $indir/abixy-het put '$depth = depth($*)'
run_mlr --from $indir/xyz2 put '$depth= depth({3:4, 5:{6:7}, 8:{9:{10:11}}})'
run_mlr --from $indir/xyz2 put 'o = {3:4, 5:{6:7}, 8:{9:{10:11}}}; $depth = depth(o)'
run_mlr --from $indir/xyz2 put '@o = {3:4, 5:{6:7}, 8:{9:{10:11}}}; $depth = depth(@o)'

run_mlr --from $indir/abixy-het put '$leafcount = leafcount($a)'
run_mlr --from $indir/abixy-het put '$leafcount = leafcount($*)'
run_mlr --from $indir/xyz2 put '$leafcount= leafcount({3:4, 5:{6:7}, 8:{9:{10:11}}})'
run_mlr --from $indir/xyz2 put 'o = {3:4, 5:{6:7}, 8:{9:{10:11}}}; $leafcount = leafcount(o)'
run_mlr --from $indir/xyz2 put '@o = {3:4, 5:{6:7}, 8:{9:{10:11}}}; $leafcount = leafcount(@o)'

# xxx ternary operator (and most RHS expressions in the grammar) doesn't handle maps.
# run_mlr --from $indir/abixy-het put -q 'o = haskey(NR==4 ? {"a": NF} : {"b": NF}, "a"); print "NR=".NR.",haskeya=".o'
