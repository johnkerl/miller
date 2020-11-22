# ----------------------------------------------------------------
announce JOIN

run_mlr --opprint join -s                -f $indir/joina.dkvp -l l -r r -j o $indir/joinb.dkvp
run_mlr --opprint join                   -f $indir/joina.dkvp -l l -r r -j o $indir/joinb.dkvp

run_mlr --opprint join -s      --ul      -f $indir/joina.dkvp -l l -r r -j o $indir/joinb.dkvp
run_mlr --opprint join         --ul      -f $indir/joina.dkvp -l l -r r -j o $indir/joinb.dkvp

run_mlr --opprint join -s           --ur -f $indir/joina.dkvp -l l -r r -j o $indir/joinb.dkvp
run_mlr --opprint join              --ur -f $indir/joina.dkvp -l l -r r -j o $indir/joinb.dkvp

run_mlr --opprint join -s --ul      --ur -f $indir/joina.dkvp -l l -r r -j o $indir/joinb.dkvp
run_mlr --opprint join         --ul --ur -f $indir/joina.dkvp -l l -r r -j o $indir/joinb.dkvp

run_mlr --opprint join -s --np --ul      -f $indir/joina.dkvp -l l -r r -j o $indir/joinb.dkvp
run_mlr --opprint join    --np --ul      -f $indir/joina.dkvp -l l -r r -j o $indir/joinb.dkvp

run_mlr --opprint join -s --np      --ur -f $indir/joina.dkvp -l l -r r -j o $indir/joinb.dkvp
run_mlr --opprint join    --np      --ur -f $indir/joina.dkvp -l l -r r -j o $indir/joinb.dkvp

run_mlr --opprint join -s --np --ul --ur -f $indir/joina.dkvp -l l -r r -j o $indir/joinb.dkvp
run_mlr --opprint join    --np --ul --ur -f $indir/joina.dkvp -l l -r r -j o $indir/joinb.dkvp

run_mlr join -l l -r r -j j -f $indir/joina.dkvp $indir/joinb.dkvp
run_mlr join -l l      -j r -f $indir/joina.dkvp $indir/joinb.dkvp
run_mlr join      -r r -j l -f $indir/joina.dkvp $indir/joinb.dkvp

run_mlr --opprint join -s                -f /dev/null -l l -r r -j o $indir/joinb.dkvp
run_mlr --opprint join                   -f /dev/null -l l -r r -j o $indir/joinb.dkvp

run_mlr --opprint join -s      --ul      -f /dev/null -l l -r r -j o $indir/joinb.dkvp
run_mlr --opprint join         --ul      -f /dev/null -l l -r r -j o $indir/joinb.dkvp

run_mlr --opprint join -s           --ur -f /dev/null -l l -r r -j o $indir/joinb.dkvp
run_mlr --opprint join              --ur -f /dev/null -l l -r r -j o $indir/joinb.dkvp

run_mlr --opprint join -s      --ul --ur -f /dev/null -l l -r r -j o $indir/joinb.dkvp
run_mlr --opprint join         --ul --ur -f /dev/null -l l -r r -j o $indir/joinb.dkvp

run_mlr --opprint join -s --np --ul      -f /dev/null -l l -r r -j o $indir/joinb.dkvp
run_mlr --opprint join    --np --ul      -f /dev/null -l l -r r -j o $indir/joinb.dkvp

run_mlr --opprint join -s --np      --ur -f /dev/null -l l -r r -j o $indir/joinb.dkvp
run_mlr --opprint join    --np      --ur -f /dev/null -l l -r r -j o $indir/joinb.dkvp

run_mlr --opprint join -s --np --ul --ur -f /dev/null -l l -r r -j o $indir/joinb.dkvp
run_mlr --opprint join    --np --ul --ur -f /dev/null -l l -r r -j o $indir/joinb.dkvp


run_mlr --opprint join -s                -f $indir/joina.dkvp -l l -r r -j o /dev/null
run_mlr --opprint join                   -f $indir/joina.dkvp -l l -r r -j o /dev/null

run_mlr --opprint join -s      --ul      -f $indir/joina.dkvp -l l -r r -j o /dev/null
run_mlr --opprint join         --ul      -f $indir/joina.dkvp -l l -r r -j o /dev/null

run_mlr --opprint join -s           --ur -f $indir/joina.dkvp -l l -r r -j o /dev/null
run_mlr --opprint join              --ur -f $indir/joina.dkvp -l l -r r -j o /dev/null

run_mlr --opprint join -s      --ul --ur -f $indir/joina.dkvp -l l -r r -j o /dev/null
run_mlr --opprint join         --ul --ur -f $indir/joina.dkvp -l l -r r -j o /dev/null

run_mlr --opprint join -s --np --ul      -f $indir/joina.dkvp -l l -r r -j o /dev/null
run_mlr --opprint join    --np --ul      -f $indir/joina.dkvp -l l -r r -j o /dev/null

run_mlr --opprint join -s --np      --ur -f $indir/joina.dkvp -l l -r r -j o /dev/null
run_mlr --opprint join    --np      --ur -f $indir/joina.dkvp -l l -r r -j o /dev/null

run_mlr --opprint join -s --np --ul --ur -f $indir/joina.dkvp -l l -r r -j o /dev/null
run_mlr --opprint join    --np --ul --ur -f $indir/joina.dkvp -l l -r r -j o /dev/null

run_mlr --odkvp join -j a -f $indir/join-het.dkvp $indir/abixy-het
run_mlr --odkvp join -j a -f $indir/abixy-het     $indir/join-het.dkvp
run_mlr --odkvp join --np --ul --ur -j a -f $indir/join-het.dkvp $indir/abixy-het
run_mlr --odkvp join --np --ul --ur -j a -f $indir/abixy-het     $indir/join-het.dkvp

run_mlr --idkvp --oxtab join --lp left_ --rp right_ -j i -f $indir/abixy-het $indir/abixy-het

for sorted_flag in "-s" ""; do
  for pairing_flags in "" "--np --ul" "--np --ur"; do
    for i in 1 2 3 4 5 6; do
      run_mlr join $sorted_flag $pairing_flags -l l -r r -j j -f $indir/het-join-left $indir/het-join-right-r$i
      for j in 1 2 3 4 5 6; do
        if [ "$i" -le "$j" ]; then
          run_mlr join $sorted_flag $pairing_flags -l l -r r -j j -f $indir/het-join-left $indir/het-join-right-r$i$j
        fi
      done
    done
  done
done
