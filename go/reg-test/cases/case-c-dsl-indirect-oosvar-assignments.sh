run_mlr --opprint put -v '@s = NR; $t = @s; $u=@["s"]; $v = $t - $u' $indir/abixy

run_mlr put -v '@t["u"] = NR; $tu = @["t"]["u"]; emitp all' $indir/abixy
run_mlr put -v '@t["u"] = NR; $tu = @["t"]["u"]; emitp @*' $indir/abixy

run_mlr put -v '@["s"] = $x; emitp all' $indir/abixy

run_mlr put -v '@["t"]["u"] = $y; emitp all' $indir/abixy

# xxx @* on the right
# xxx @* on the left
