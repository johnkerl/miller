# Experiment (negative result): is `sort` bottlenecked on GC heap-scanning?

Status: **premise disproven** — no code change. Branch opened so the analysis
is reviewable, then closed.

## Premise

Accumulating verbs like `sort` hold the entire input in memory (~1M records for
`big.csv`). The hypothesis (from the earlier investigation's "remaining
frontiers" notes) was that they are bottlenecked on the **garbage collector
scanning that large, pointer-heavy retained heap**, and that a more compact
record representation would speed them up.

## Measurements (`big.csv`, 1M rows, `mlr --csv sort -nf quantity`, M1)

A CPU profile *looks* damning — `runtime.scanobject` ~30%, `madvise` ~12%,
GC-bitmap reads ~6% — about half the samples are GC/memory management. But CPU
samples are not wall-clock on a multicore box. Varying the GC knobs:

| Setting | wall-clock | RSS | GC cycles |
| --- | --- | --- | --- |
| `GOGC=500` (default) | 1.89s | 1233MB | 3 |
| `GOGC=1000` | 1.86s | 1226MB | 3 |
| `GOGC=2000` | 1.90s | 1213MB | 3 |
| `GOGC=off` | 1.87s | 1225MB | 3 |
| `GOMEMLIMIT=2500MiB` | 1.84s | 1202MB | 4 |

**Turning GC off entirely changes wall-clock by ~0%.** The scanning is real CPU
work, but only 3 cycles run and the concurrent collector overlaps them on
otherwise-idle cores, so they don't gate wall-clock.

Single-core, where GC *cannot* overlap:

| Setting | wall-clock |
| --- | --- |
| `GOMAXPROCS=1 GOGC=500` | 1.91s |
| `GOMAXPROCS=1 GOGC=off` | 2.08s |

Even with no overlap available, `GOGC=off` is *slower* (the unbounded heap costs
more in allocation/`madvise`/cache than the scans it avoids). So GC scanning
does not gate `sort`'s wall-clock on either multi- or single-core. (An earlier
ad-hoc "~35%" reading was machine-load noise.)

## Why `sort` costs what it does

The ~1.6s gap between `sort` (~1.85s) and `cat` (~0.22s) is mostly **structural,
not GC**: `sort` is non-streaming, so it cannot overlap reading and writing the
way `cat` does — it must read all input, then sort, then write. The remaining
critical-path costs are output I/O (`syscall` ~15%, an intrinsic floor), the
auto-flatten transformer's per-record `isFlattenable` guard (~8%, already
short-circuited and not cheaply removable), and the sort itself. None is a
large, easily-removable lever.

## The one real (but different) prize: memory, not speed

`sort` uses ~1.2GB RSS for a 44MB input — a ~27× blow-up, driven by the
pointer-heavy record model (`Mlrmap` → linked `MlrmapEntry` → `*Mlrval`, each
`Mlrval` 56 bytes with `string`/`interface{}`/`error` fields). A more compact
representation would lower that ceiling (letting users sort larger files before
OOM) — but it would **not** improve wall-clock, since GC scanning isn't the
gate. That is a large, invasive, correctness-sensitive data-model change for a
memory-only benefit, so it is out of scope here.

## Verdict

Premise disproven: GC scanning is not `sort`'s wall-clock bottleneck. No
speed-oriented change is warranted. Compact records remain a possibility *only*
if reducing `sort`'s memory ceiling (max sortable file size) becomes a goal in
its own right.
