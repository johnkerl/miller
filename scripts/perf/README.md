# Miller Performance Benchmarks

Scripts for reproducing the performance graphs at
https://miller.readthedocs.io/en/latest/performance/

## Quick start

```
# Once — if ~/data/big.csv doesn't already exist (run from repo root):
bash scripts/perf/prep-perf-data.sh

# Collect timings and render graphs (run from scripts/perf/):
cd scripts/perf
bash run-perf.sh ~/bin/mlr-6.18.1 ~/bin/mlr-6.19.0
```

That runs `time-verbs.py` with both executables (5 reps each, ~18 cases),
saves `timings-YYYY-MM-DD.dat`, then calls `plot-timings.sh` to produce three
PNGs: `-verbs.png`, `-chains.png`, `-cats.png`.

## Prerequisites

- Python 3
- [`pgr`](https://github.com/johnkerl/pgr) (used by `plot-timings.sh` to render PNG graphs)
- One or more Miller executables to compare

## Step-by-step

### Step 1 — Generate test data (one-time)

Run from the Miller repo root:

```
bash scripts/perf/prep-perf-data.sh
```

This creates `~/data/big.csv` (~1 million rows) and derived files in DKVP,
NIDX, XTAB, and JSON formats, plus `small.csv` and `medium.csv` subsets.

### Step 2 — Collect timings

Run from `scripts/perf/`, passing one or more executables to compare:

```
python time-verbs.py ~/bin/mlr-6.18.1 ~/bin/mlr-6.19.0 \
  > timings-$(date +%Y-%m-%d).dat
```

Output is DKVP, one record per (case, executable): `desc=...,version=...,seconds=...`

Five reps are averaged per case. Cases cover:

- **Verbs**: `check`, `cat`, `tail`, `tac`, `sort -f`, `sort -n`, `stats1`
- **Then-chains**: one to four chained `put -f chain-1.mlr` steps
- **Formats**: `cat` across CSV, CSVLITE, DKVP, NIDX, XTAB, JSON

### Step 3 — Plot

```
bash plot-timings.sh timings-YYYY-MM-DD.dat
```

Produces three PNGs alongside the dat file:
- `*-verbs.png` — verb timings
- `*-chains.png` — then-chain depth
- `*-cats.png` — format comparison

## Files

| File | Purpose |
|------|---------|
| `prep-perf-data.sh` | Generate `~/data/big.csv` and format variants (one-time) |
| `time-verbs.py` | Time each verb/format/chain case; outputs DKVP |
| `chain-1.mlr` | DSL script used for then-chain benchmarks |
| `run-perf.sh` | Wrapper: collect timings + plot |
| `plot-timings.sh` | Read a `.dat` file and render PNGs via `pgr` |
| `timings-2026-02-22.sh` | Original one-off plot script (hardcoded filename) |
| `timings-2026-02-22.dat` | Timing data from the February 2026 run |
