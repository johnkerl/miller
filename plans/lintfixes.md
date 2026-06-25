# golangci-lint findings

golangci-lint v2.12.2, run against `./cmd/mlr ./pkg/...` (same as CI).

- **134 total issues**
- **64 unique source files**

## By linter (summary line from golangci-lint)

| Count | Linter |
|------:|--------|
| 50 | errcheck |
| 50 | staticcheck |
| 30 | ineffassign |
| 3 | govet |
| 1 | unused |

## By semantic type

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

## Unique source files with issues (64)

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
