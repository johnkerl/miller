#!/usr/bin/env python

import os
import sys
import time
import subprocess

# ================================================================
def main():
    mlrs       = [ "mlr" ]
    #
    kind       = "big"
    #kind      = "medium"
    #kind      = "small"
    in_csv     = f"$HOME/data/{kind}.csv"
    in_csvlite = f"$HOME/data/{kind}.csv"
    in_dkvp    = f"$HOME/data/{kind}.dkvp"
    in_nidx    = f"$HOME/data/{kind}.nidx"
    in_xtab    = f"$HOME/data/{kind}.xtab"
    in_json    = f"$HOME/data/{kind}.json"
    #
    cases = [
        ["check",  f"--csv --from {in_csv} check"],
        ["cat",    f"--csv --from {in_csv} cat"],
        ["tail",   f"--csv --from {in_csv} tail"],
        ["tac",    f"--csv --from {in_csv} tac"],
        ["sort-f", f"--csv --from {in_csv} sort -f shape"],
        ["sort-n", f"--csv --from {in_csv} sort -n quantity"],
        ["stats1", f"--csv --from {in_csv} stats1 -a min,mean,max -f quantity,rate -g shape"],
        None,
        ["chain-1", f"--csv --from {in_csv} put -f scripts/perf/chain-1.mlr"],
        ["chain-2", f"--csv --from {in_csv} put -f scripts/perf/chain-1.mlr then put -f scripts/perf/chain-1.mlr"],
        ["chain-3", f"--csv --from {in_csv} put -f scripts/perf/chain-1.mlr then put -f scripts/perf/chain-1.mlr then put -f scripts/perf/chain-1.mlr"],
        ["chain-4", f"--csv --from {in_csv} put -f scripts/perf/chain-1.mlr then put -f scripts/perf/chain-1.mlr then put -f scripts/perf/chain-1.mlr then put -f scripts/perf/chain-1.mlr"],
        None,
        ["cat-csv",     f"--csv     --from {in_csv}  cat"],
        ["cat-csvlite", f"--csvlite --from {in_csvlite}  cat"],
        ["cat-dkvp",    f"--dkvp    --from {in_dkvp} cat"],
        ["cat-nidx",    f"--nidx    --from {in_nidx} cat"],
        ["cat-xtab",    f"--xtab    --from {in_xtab} cat"],
        ["cat-json",    f"--json    --from {in_json} cat"],
    ]
    #
    nreps = 5

    if len(sys.argv) > 1:
        mlrs = sys.argv[1:]


    for case in cases:
        if case is None:
            print()
            continue
        for mlr in mlrs:
            desc, args = case
            avg = time_runs(desc, mlr, args, nreps)
            print(f"desc={desc},version={mlr},seconds={avg:.3f}")

# ----------------------------------------------------------------
def time_runs(desc, mlr, args, nreps):
    args = os.path.expandvars(args)
    cmd = f"{mlr} {args}"
    times = [time_run(mlr, args) for _ in range(nreps)]
    avg = sum(times) / len(times)
    return avg

# ----------------------------------------------------------------
def time_run(mlr, args):
    cmd = f"{mlr} {args} > /dev/null"
    start = time.perf_counter()
    subprocess.run(cmd, shell=True)
    elapsed = time.perf_counter() - start
    return elapsed

# ================================================================
if __name__ == "__main__":
    main()
