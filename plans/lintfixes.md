# golangci-lint fix tracking

Tracking issue: https://github.com/johnkerl/miller/issues/2109

golangci-lint v2.12.2, run against `./cmd/mlr ./pkg/...` (same as CI).

## Baseline (first run, branch `johnkerl/lintfixes` before any fixes)

- **134 total issues**
- **64 unique source files**

### By linter

| Count | Linter |
|------:|--------|
| 50 | errcheck |
| 50 | staticcheck |
| 30 | ineffassign |
| 3 | govet |
| 1 | unused |

### By semantic type

| Count | Category |
|------:|---------|
| 30 | `ineffassign`: variable assigned but never used after |
| 29 | `staticcheck QF1003`: use typed switch instead of if-else chain |
| 9 | `errcheck`: ignore return from `fmt.Fprintf/Fprint/Fprintln` |
| 8 | `errcheck`: ignore return from map `RemoveIndexed`/`PutIndexed` |
| 8 | `staticcheck ST1023/QF1011`: omit explicit type (inferred from RHS) |
| 6 | `errcheck`: ignore return from `Write()`/`WriteString()` |
| 6 | `errcheck`: ignore return from DSL emit/execute/redirect funcs |
| 5 | `errcheck`: ignore return from `cli.Finalize{Reader,Writer}Options` |
| 4 | `errcheck`: ignore return from `Close()` |
| 4 | `errcheck`: ignore return from `Flush()` |
| 3 | `errcheck`: generic unchecked return (anonymous funcs) |
| 3 | `errcheck`: ignore return from other stream/process ops |
| 3 | `staticcheck S1009`: omit nil check before `len()` (always safe for nil slices) |
| 3 | `staticcheck SA9003`: empty `else if` branch |
| 3 | `staticcheck QF1007`: merge conditional assignment into declaration |
| 2 | `errcheck`: ignore return from `os.Setenv`/`os.Unsetenv` |
| 2 | `govet stdmethods`: `MarshalJSON` has wrong signature |
| 2 | `staticcheck QF1001`: apply De Morgan's law to simplify negation |
| 1 | `govet`: unreachable code |
| 1 | `staticcheck QF1006`: lift check into loop condition |
| 1 | `staticcheck S1031`: unnecessary nil check before `range` |
| 1 | `unused`: unused struct field |
| **134** | **TOTAL** |

## PRs merged against issue #2109

| PR | Title | Merged |
|----|-------|--------|
| [#2076](https://github.com/johnkerl/miller/pull/2076) | ci: add golangci-lint workflow | 2026-06-25 |
| [#2108](https://github.com/johnkerl/miller/pull/2108) | Next batch of lint fixes | 2026-06-28 |
| [#2110](https://github.com/johnkerl/miller/pull/2110) | Next batch of lint fixes | 2026-06-28 |
| [#2112](https://github.com/johnkerl/miller/pull/2112) | Convert if/else-if chains to typed switch (staticcheck QF1003) | 2026-06-28 |

PR #2112 alone fixed ~29 QF1003 findings across 72 files. PRs #2108 and #2110 addressed other batches.

**Note:** The `johnkerl/lintfixes` branch listed commits `0f8b931bc` (MarshalJSON rename) and
`f57edc037` (unreachable code) in an earlier version of this file, but those were never pushed —
`origin/johnkerl/lintfixes` ends at the plan-file commit. Those fixes were redone on `johnkerl/lint4`.

## Work done on branch `johnkerl/lint4` (2026-07-03)

- Rename `MarshalJSON` → `FormatAsJSON` on `Mlrval` and `Mlrmap` (fixes `govet stdmethods`, 2 findings;
  13 references across 6 files). The old signature shadowed `json.Marshaler` with an incompatible one.
- Remove unreachable `return nil` after exhaustive if-else in `pkg/mlrval/mlrval_collections.go`
  (fixes `govet unreachable`, 1 finding).
- The `unused` finding (`exeName` field in `pkg/terminals/script/types.go`) was already resolved on
  `main` by an unrelated change (field renamed to `name` and now referenced).

## Remaining after `johnkerl/lint4` govet batch (2026-07-03)

**84 issues: 50 errcheck, 34 staticcheck.** All `govet`, `unused`, and `ineffassign` findings are resolved.

Staticcheck breakdown: 16 ST1023 (omit explicit type), 3 SA9003 (empty branch), 3 S1009 (redundant
nil check before len), 3 QF1007 (merge conditional assignment), 3 QF1006 (lift into loop condition),
3 QF1001 (De Morgan), 2 QF1011 (omit explicit type), 1 S1031 (nil check before range).

## Priority order for remaining fixes

1. `errcheck` — 50 findings. Varies in seriousness:
   - DSL emit/execute return values, `RemoveIndexed`/`PutIndexed`, `Flush`/`Close` — worth checking individually
   - `fmt.Fprintf`/`fmt.Fprint` to stderr — genuinely safe to ignore; consider `//nolint:errcheck`
   - `os.Setenv`/`os.Unsetenv` — low-risk but cheap to fix
2. `staticcheck` style — 34 findings, mostly mechanical (type inference, nil-check cleanup, etc.)

## Proposed round 5 (branch `johnkerl/lint5`, 2026-07-03): staticcheck batch

Recommendation: do all 34 staticcheck findings as one mechanical PR, and defer errcheck to
rounds 6+. This inverts the priority order above, deliberately: the staticcheck fixes are
low-risk and retire an entire linter, while the errcheck sites need per-site judgment and a
few of them (see rounds 6+ below) may be real user-facing bugs deserving focused PRs.

### 1. ST1023 + QF1011 — omit explicit type (18 findings, mechanical)

Drop the redundant type from `var x T = <expr-of-type-T>` declarations:

- `pkg/bifs/arithmetic.go:721,722,727,739,740,904` — RHS is already `float64(...)`/`int64(...)` conversions
- `pkg/bifs/bits.go:226,227`
- `pkg/dsl/cst/udf.go:566`
- `pkg/input/record_reader_dkvp_nidx.go:219`
- `pkg/transformers/put_or_filter.go:195`
- `pkg/transformers/split.go:91,92,94,96,98,99,101` — cluster of `var x bool = false` / `var x string = "..."`

### 2. S1009 + S1031 — redundant nil checks (4 findings, mechanical)

- `pkg/dsl/cst/blocks.go:66`, `pkg/dsl/cst/evaluable.go:20,100` — drop `x != nil &&` before `len(x)`
- `pkg/mlrval/mlrval_yaml.go:200` — drop nil check before `range`

### 3. QF1007 — merge conditional assignment into declaration (3 findings, mechanical)

All three are the `found := false; if cond { found = true }` pattern → `found := cond`:

- `pkg/terminals/help/entry.go:824` (`helpByExactSearchOne`), `:860` (`helpByApproximateSearchOne`)
- `pkg/terminals/repl/verbs.go:877` (`handleHelpFindSingle`)

### 4. QF1006 — lift break condition into loop condition (3 findings, mechanical)

- `pkg/bifs/relative_time.go:31,93` — `for { if remainingInput == "" { break } ... }` →
  `for remainingInput != "" { ... }`
- `pkg/mlrval/mlrmap_accessors.go:759` (`Mlrmap.Label`) — `for { if i >= numNewNames { break } ... }` →
  `for i < numNewNames { ... }` (the interior `pe == nil` break stays)

### 5. QF1001 — De Morgan (3 findings, small judgment calls)

- `pkg/dsl/cst/builtin_functions.go:878,960` — `!(btype == MT_ABSENT || btype == MT_BOOL)` →
  `btype != mlrval.MT_ABSENT && btype != mlrval.MT_BOOL`. Equally readable; apply.
- `pkg/terminals/help/entry_which.go:167` — the tokenizer's rune predicate. Mechanical De Morgan
  gives `(r < 'a' || r > 'z') && ...` which is arguably worse; prefer rewriting positively:
  `return !isWordRune(r)` with a small helper (or inline
  `!(('a' <= r && r <= 'z') || ('0' <= r && r <= '9') || r == '-' || r == '_')` restructured
  so the negation is outermost, which is what staticcheck accepts).

### 6. SA9003 — empty branches (3 findings, comment-carrying branches)

All three are deliberate no-op branches that exist only to hold an explanatory comment.
Delete the empty branch and keep the comment adjacent — no behavior change:

- `pkg/dsl/cst/for.go:185,375` — `} else if indexMlrval.IsAbsent() { // data-heterogeneity no-op }` →
  drop the branch, keep a `// else if absent: data-heterogeneity no-op` comment after the closing brace
- `pkg/terminals/repl/verbs.go:203` — `if len(args) == 0 { // zero file names is stdin, which is readable }` →
  drop the block, keep the comment above the following loop

After this round: **50 issues remain, all errcheck.**

## Proposed rounds 6+ : errcheck (50 findings, grouped by treatment)

### Round 6a — propagate: likely real bugs (7 findings)

- `cli.FinalizeReaderOptions` / `cli.FinalizeWriterOptions` (6): `pkg/terminals/repl/entry.go:150,151`,
  `pkg/terminals/script/entry.go:126,127`, `pkg/transformers/join.go:256`,
  `pkg/transformers/put_or_filter.go:343`. `FinalizeReaderOptions` returns
  `unrecognized input format %q` — ignoring it means e.g. `mlr join -i badformat` silently
  proceeds with wrong separators instead of erroring. `join.go`/`put_or_filter.go` are inside
  verb constructors that already return `(transformer, error)`, so propagation is easy; the
  repl/script entry points should print and exit.
- `pkg/stream/stream.go:104` — final `bufferedOutputStream.Flush()` of the main output stream.
  A failure here (full disk, closed pipe) currently loses output silently while exiting 0.
  Propagate into `retval`.

### Round 6b — inspect individually: post-#2129 meaningful errors (4 findings)

`PutIndexed`/`RemoveIndexed` at `pkg/mlrval/mlrmap_flatten_unflatten.go:182,229` and
`pkg/runtime/stack.go:457,488`. Now that #2129 fixed the masked error path, these returns are
meaningful. Decide per site whether to propagate or `_ =` with a comment saying why ignoring
is correct (e.g. index provably in bounds at the call site).

### Round 7 — explicit ignores / config (39 findings)

- **Usage printers, `fmt.Fprint*` (9):** `pkg/terminals/repl/entry.go:33,44`,
  `pkg/terminals/script/entry.go:24,25`, `pkg/transformers/altkv.go:26`,
  `pkg/transformers/count.go:27`, `pkg/transformers/fill_down.go:26,27,28`.
  Decision point: add a `.golangci.yml` (repo currently has none — CI runs defaults) with
  `errcheck.exclude-functions` for `fmt.Fprint/Fprintf/Fprintln`, which also future-proofs
  against new usage functions; or `_ =` each site. Config is less churn.
- **REPL/script interactive write/flush paths (11):** `pkg/terminals/repl/entry.go:189`,
  `pkg/terminals/repl/verbs.go:557,558,620,621,645,668,698`, `pkg/terminals/script/entry.go:146`,
  `pkg/terminals/script/runner.go:130,131`. Interactive-session output; on failure (EPIPE etc.)
  there is nothing useful to do. `_ =` with a brief comment.
- **Cleanup/teardown (14):** `cmd/mlr/main.go:58`, `pkg/cli/option_parse.go:3461`,
  `pkg/lib/readfiles.go:54,83`, `pkg/lib/halfpipe.go:43,90`,
  `pkg/terminals/regtest/invoker.go:75,76`, `pkg/terminals/regtest/regtester.go:152,155,161,331,464,545`.
  Read-side `Close()`, `os.Remove` of temp files, `os.Setenv` in the regtest harness — `_ =` throughout.
- **In-memory pipe for usage capture (3):** `pkg/transformers/aaa_transformer_json.go:43,48,50` —
  `io.Copy`/`Close` on an `os.Pipe` used to capture usage text into a buffer; cannot meaningfully
  fail. `_ =`.
- **`pkg/dsl/cst/dump.go:205,247` (2):** DSL `dump`-to-file output handlers. Check whether the
  surrounding output-handler code path has an error convention to join (it returns error in
  places); otherwise `_ =`.

## Bug noticed in passing (fixed on `johnkerl/lint4`)

In `pkg/mlrval/mlrval_collections.go`, `removeIndexedOnArray` with a single in-bounds index removed
the element but then fell through to `return errors.New("array index out of bounds for unset")` —
the success path never returned nil. Callers ignored the error (several of the errcheck findings),
which masked it. Fixed by returning nil on the success path, matching `removeIndexedOnMap`;
unit tests added in `pkg/mlrval/mlrval_collections_test.go`. No observable behavior change today,
but errcheck fixes for `RemoveIndexed` call sites can now propagate the error meaningfully.

## Unique source files with issues (baseline, 64 files)

```
pkg/auxents/hex.go
pkg/auxents/lecat.go
pkg/auxents/termcvt.go
pkg/auxents/unhex.go
pkg/bifs/random.go
pkg/cli/option_parse.go
pkg/climain/mlrcli_mlrrc.go
pkg/climain/mlrcli_shebang.go
pkg/dkvpx/dkvpx_reader.go
pkg/dsl/cst/blocks.go
pkg/dsl/cst/builtin_functions.go
pkg/dsl/cst/dump.go
pkg/dsl/cst/emit_emitp.go
pkg/dsl/cst/evaluable.go
pkg/dsl/cst/for.go
pkg/dsl/cst/if.go
pkg/dsl/cst/lvalues.go
pkg/dsl/cst/print.go
pkg/dsl/cst/statements.go
pkg/dsl/cst/udf.go
pkg/dsl/cst/uds.go
pkg/dsl/cst/validate.go
pkg/input/line_reader.go
pkg/input/record_reader_xtab.go
pkg/mlrval/mlrmap_accessors.go
pkg/mlrval/mlrmap_flatten_unflatten.go
pkg/mlrval/mlrmap_json.go
pkg/mlrval/mlrval_collections.go
pkg/mlrval/mlrval_copy.go
pkg/mlrval/mlrval_infer.go
pkg/mlrval/mlrval_json.go
pkg/mlrval/mlrval_new.go
pkg/mlrval/mlrval_yaml.go
pkg/stream/stream.go
pkg/terminals/help/entry.go
pkg/terminals/repl/dsl.go
pkg/terminals/repl/entry.go
pkg/terminals/repl/verbs.go
pkg/terminals/script/types.go
pkg/terminals/terminals.go
pkg/transformers/altkv.go
pkg/transformers/bar.go
pkg/transformers/bootstrap.go
pkg/transformers/case.go
pkg/transformers/cat.go
pkg/transformers/check.go
pkg/transformers/count.go
pkg/transformers/fill_down.go
pkg/transformers/flatten.go
pkg/transformers/fraction.go
pkg/transformers/group_by.go
pkg/transformers/having_fields.go
pkg/transformers/join.go
pkg/transformers/json_stringify.go
pkg/transformers/merge_fields.go
pkg/transformers/nest.go
pkg/transformers/put_or_filter.go
pkg/transformers/reshape.go
pkg/transformers/seqgen.go
pkg/transformers/sort.go
pkg/transformers/split.go
pkg/transformers/summary.go
pkg/transformers/tee.go
pkg/transformers/unspace.go
```
