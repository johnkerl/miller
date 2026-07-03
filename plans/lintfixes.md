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

## Possible bug noticed in passing (not lint scope)

In `pkg/mlrval/mlrval_collections.go`, `removeIndexedOnArray` with a single in-bounds index removes
the element but then falls through to `return errors.New("array index out of bounds for unset")` —
the success path never returns nil. Callers currently ignore this error (that's several of the
errcheck findings), so fixing errcheck here will surface this. Worth deciding intended behavior first.

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
