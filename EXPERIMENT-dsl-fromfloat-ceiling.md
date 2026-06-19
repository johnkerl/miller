# Experiment (negative result): eliminating DSL numeric-temporary allocations

Status: **not pursued** — bounded ceiling too low to justify the cost/risk.
No code change beyond this note; branch opened so the analysis is reviewable,
then closed.

## Background

This continues the allocation-focused performance work (see PRs #2081, #2082,
#2083, #2086). After pooling DSL stack frames (#2086), the dominant remaining
allocator on the `put`/`filter` path is **`FromFloat`** (~6.8M objects for
`put -f chain-1.mlr` over 1M-record CSV), with `FromInt` close behind. These
are the interior temporaries of arithmetic expression evaluation: every binary
op (`+`, `/`, `**`, …) and math function (`log10`, …) returns a freshly
heap-allocated `*mlrval.Mlrval`.

The proposed fix was the classic tree-walking-interpreter optimization:
node-owned result slots + in-place bif variants, so an expression node computes
into reusable storage instead of allocating per evaluation.

## Why it was not pursued

### 1. The performance ceiling is ~10%

Before building anything, the ceiling was measured cheaply by hacking
`FromFloat` and `FromInt` to reuse package-global `Mlrval`s — i.e. making those
allocations *completely free*. (Output is garbage, but `chain-1` has no
value-dependent branching, so wall-clock timing is representative.)

| workload | baseline (post #2086) | FromFloat+FromInt allocation-free |
|---|---|---|
| `put` chain-1 | 0.71s | **0.64s (~10%)** |
| `put` chain-4 | 0.86s | ~0.79s (~8%) |

So eliminating *every* DSL numeric-temporary allocation — the entire point of
the refactor — yields ~10% on the put-heavy workload and ~0 on `cat`/`sort`.
A CPU profile explains why: `put` is ~27% I/O-bound (output writing), and on a
multicore box the allocation/GC work largely overlaps idle cores, so the huge
`FromFloat` *count* is not the wall-clock gate.

### 2. The cost and risk are large

- **Surface area:** `bifs.BinaryFunc` / `UnaryFunc` are `func(...) *Mlrval`,
  dispatched through `[MT_DIM]`/`[MT_DIM][MT_DIM]` disposition matrices whose
  cells each allocate via `FromFloat`/`FromInt`/etc. Making results in-place
  means changing the function-type signatures and **hundreds** of cells.
- **Aliasing / correctness:** reusing result storage is only safe if nothing
  retains a result past its immediate consumer. Scalar retention points *do*
  copy (verified: field `PutCopy`, oosvar direct `PutCopy`, local
  `value.Copy()`), but indexed/collection stores (e.g. `Oosvars.PutIndexed`)
  appear to store **by reference**, relying on the current always-fresh
  invariant. Breaking that would cause silent cross-record data corruption —
  the worst kind of bug — and would require a full audit + copy-on-store fixes.

## Verdict

~10% bounded upside on one workload, for a large refactor with real
data-corruption risk. Not worth it. The clean, safe DSL win (stack-frame
pooling, #2086) is kept; this avenue is closed.
