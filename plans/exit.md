# Plan: survey & phased removal of scattered `os.Exit` calls

## Context

Miller has ~170 `os.Exit` call sites scattered through non-test Go code under `pkg/`
(plus a handful in `cmd/experiments/`, which are throwaway sandboxes and out of scope).
The north star: `os.Exit` should happen only in `main` or its immediate callee
(`entrypoint.Main`), with everything below returning `error`. Not every last site must
be eliminated — genuine internal-coding-error assertions can stay — but user-facing
error and exit paths should propagate.

See also [issue 2200](https://github.com/johnkerl/miller/issues/2200).

Motivation:

- **General hygiene** — exits deep in library code make the packages unusable as a
  library, skip deferred cleanup (e.g. the pprof defers in `cmd/mlr/main.go`), and
  create unreachable/untestable code paths.
- **Pre-work for [#341](https://github.com/johnkerl/miller/issues/341)** (DSL `exit`
  statement) — an `exit 0`/`exit N` from inside a `put` expression must shut the
  record stream down cleanly (flush writers, stop readers) and carry a *chosen* exit
  code up to `main`. Raw `os.Exit` mid-goroutine can't do that.
- **Pre-work for [#440](https://github.com/johnkerl/miller/issues/440)** (strict mode)
  — missing-field reads become *errors*, which need a propagation path from the DSL
  evaluator up through the transformer chain to the process exit code.

Prior art already in-tree (this effort is a continuation, not a fresh start):

- Commit `6b32e1f41` "Reduce number of os.Exit callsites, part 1 of n" (#1055).
- `pkg/cli/errors.go`: sentinels `cli.ErrHelpRequested` (→ exit 0) and
  `cli.ErrUsagePrinted` (→ exit 1).
- `pkg/cli/verb_utils.go`: error-returning verb-arg helpers (`VerbCheckArgCount`,
  `VerbGetStringArg`, …, `VerbErrorf`) — the stated replacement for exit-on-failure
  helpers.
- `TransformerParseCLIFunc` already returns `(RecordTransformer, error)`
  (`pkg/transformers/aaa_record_transformer.go`).
- `pkg/climain/errors_json.go`: typed `CLIError{Kind,Token,Msg}` + `--errors-json`
  structured emission.
- `mlrval.FromErrorString` error-valued Mlrvals + `state.NoExitOnFunctionNotFound`
  (used in `pkg/dsl/cst/udf.go`) — the DSL-layer error convention.
- REPL precedent: `pkg/terminals/repl/dsl.go` prints DSL errors and continues (#1976),
  while the same errors from `cst.ExecuteMainBlock` make `put_or_filter.go` exit.

## Survey: where the ~170 sites live

Counts are non-test `.go` under `pkg/`. Breakdown: 151 × `os.Exit(1)`, 15 ×
`os.Exit(0)`, plus `os.Exit(exitCode)` / `os.Exit(entry.main(args))` dispatchers.

| Cluster | Sites (approx) | Containing signature returns error today? | Difficulty |
|---|---|---|---|
| 1. `pkg/cli` main-flag parsing: `option_parse.go` (14), `flag_types.go`, `mlrcli_util.go` (`CheckArgCount`) | ~18 | No — `FlagParser` is void (`flag_types.go:58`); `FlagTable.Parse` returns bool | Signature change, then mechanical |
| 2. `pkg/climain`: `mlrcli_parse.go` (11), `mlrcli_mlrrc.go` (`loadMlrrcOrDie`) | 12 | Yes — `ParseCommandLine` and both passes return error | Mechanical |
| 3a. Transformer ParseCLI leftovers: constructor-error exits (e.g. `split.go`, `tee.go`, `subs.go`, `nest.go`, `join.go`, `surv.go`) | ~15 | Yes — `TransformerParseCLIFunc` returns error | Mechanical |
| 3b. Transformer **runtime** (`Transform` path): `tee.go`, `split.go`, `step.go`, `put_or_filter.go`, `utils/join_bucket_keeper.go` (5) | ~20 | No — `RecordTransformer.Transform` is void | Interface change |
| 4. DSL runtime: `pkg/dsl/cst/hofs.go` (14), `udf.go` (8), `root.go`, `evaluable.go` | ~25 | No — `IEvaluable.Evaluate` returns only `*mlrval.Mlrval` | Semantic; use error-Mlrval convention |
| 5. Writers/values: `record_writer_yaml.go` (4), `record_writer_json_jsonl.go` (2), `mlrval_output.go`, `mlrval_get.go`; `pkg/lib` (`WriteTempFileOrDie`, `CompileMillerRegexOrDie`, `mlrmath.go`) | ~18 | Mixed — many have `...OrError` siblings already | Mostly mechanical |
| 6. Sub-entrypoints: `pkg/auxents/*` (hex, unhex, lecat, termcvt), `pkg/terminals/*` (repl, script, regtest, completion, mcp) | ~40 | `main(args []string) int` per aux/terminal | Self-contained, mechanical |
| 7. Deliberate exit-0: `--version`, help, `--list-color-*`, `put -e`/`--explain` "valid" paths | ~15 | Varies | Convert to sentinels mapped at entrypoint |
| Keep as-is | — | `lib.InternalCodingErrorIf` family (~365 assertion call sites), `asserting_*` builtins (`pkg/bifs/types.go` — abort is their specified behavior), `os.Exit(entry.main(args))` dispatchers in `auxents.go`/`terminals.go` | — |

## Design decisions

### D1. Single exit point

`cmd/mlr/main.go` stays exit-free (so its pprof/trace defers run); `entrypoint.Main`
becomes the one place that calls `os.Exit`, mapping errors to codes:

```go
// pkg/entrypoint: nil → 0; lib.ExitRequest{Code} → Code;
// cli.ErrHelpRequested → 0; everything else → print (or EmitStructuredError) + 1
```

The `ErrHelpRequested`/`ErrUsagePrinted` → exit translation currently mid-stack in
`mlrcli_parse.go:406-409` moves up to this boundary.

### D2. `Transform` returns `error`; exit codes ride a typed sentinel

Change the streaming interface:

```go
type RecordTransformer interface {
    Transform(inrecAndContext *types.RecordAndContext,
        outputRecordsAndContexts *[]*types.RecordAndContext,
        inputDownstreamDoneChannel <-chan bool,
        outputDownstreamDoneChannel chan<- bool,
    ) error
}
```

`runSingleTransformer`/`ChainTransformer` forward a non-nil error to a new
`chan error` consumed by `stream.Stream`'s select loop (alongside the existing
`inputErrorChannel`); the anemic `dataProcessingErrorChannel chan bool` is subsumed.
This also fixes the standing `XXX ... exit 1 & goroutine cleanup` in
`pkg/output/channel_writer.go`, whose print-and-signal-bool becomes send-the-error.

**Exit(0) vs exit(1) disambiguation for #341 lives in the error's *type*, not the
transport** — the `io.EOF` pattern (sentinel error as control flow):

```go
// pkg/lib (or pkg/types): satisfies error
type ExitRequest struct{ Code int }
```

A DSL `exit N` statement extends the existing `BlockExitPayload` machinery
(`pkg/dsl/cst/types.go:101-105` — BREAK/CONTINUE/RETURN already propagate out of
nested blocks) with a `BLOCK_EXIT_EXIT` status carrying the code; at the block
boundary it becomes `ExitRequest{N}` returned as an error from `ExecuteMainBlock` →
`put_or_filter.Transform` → chain → `stream.Stream` → `entrypoint.Main`, which exits
with `Code` — after the normal flush path.

The synchronous error return is *more* amenable to #341 than a side-channel: a
free-standing error channel races with `doneWritingChannel` in the Stream select loop
and gives no ordering guarantee that records emitted *before* the `exit` statement get
flushed. A returned error lets `ChainTransformer` forward already-produced records,
signal downstream-done (the same drain mechanism `head -n` uses), and only then
surface the `ExitRequest`.

### D3. DSL runtime errors use the existing error-Mlrval convention

The `...OrDie` helpers in `hofs.go`/`udf.go` return `mlrval.FromErrorString(...)`
values instead of exiting (generalizing what `state.NoExitOnFunctionNotFound` already
does at `udf.go:120-132`). `cst.Execute*Blocks` already return `error`;
`put_or_filter.go` stops swallowing those with `os.Exit(1)` and returns them from
`Transform` (per D2). This is exactly the plumbing #440 strict mode needs.

### D4. What stays

- `lib.InternalCodingErrorIf` / `...WithMessageIf` — should-never-happen assertions,
  by their own doc comment. (`MLR_PANIC_ON_INTERNAL_ERROR` already gives stack traces.)
- `asserting_*` DSL builtins (`pkg/bifs/types.go` `assertingCommon`) — abort on
  assertion failure is their documented contract.
- Top-of-stack dispatchers `os.Exit(entry.main(args))` in `auxents.go`/`terminals.go`
  — these *are* the immediate-callee exits; inner exits within each aux/terminal
  convert to `return <code>` up to their `main(args) int`.

### Behavioral invariants

Stderr message text and process exit codes must not change (the regression suite
asserts on both; see also recent exit-code fixes #2171, #2146). `--errors-json`
categorization must keep working — and gains coverage, since stream-time errors that
today bypass `EmitStructuredError` will now arrive at the entrypoint as errors.

## Phases (one PR each, each green under `make dev` + `make lint`)

### Phase 1 — top level + all mechanical swaps (clusters 2, 3a, 5, 7 partial)
- `entrypoint.Main() (MainReturn, error)` (or int code); the two exits in
  `entrypoint.go` become the single mapping point; introduce `ExitRequest`.
- `mlrcli_parse.go` validation exits → returned `CLIError`s; `loadMlrrcOrDie` →
  `loadMlrrc(...) error`; move help/usage-sentinel exit translation to entrypoint.
- Transformer ParseCLI constructor exits → `return nil, err` (`split.go:234` et al.).
- Migrate callers to existing error siblings: `CompileMillerRegexOrDie` →
  `CompileMillerRegex`, mlrval `...OrDie` getters → `...OrError`, `WriteTempFileOrDie`
  → error-returning version; `mlrcli_util.go CheckArgCount` callers → `VerbCheckArgCount`
  style where signatures already allow.
- Deliberate exit-0 sites reachable from error-returning frames (`--version`, `put -e`
  "valid", `subs.go:180`) → sentinel returns.

### Phase 2 — `FlagParser` signature (cluster 1)
- `type FlagParser func(...) error`; `FlagTable.Parse` → `(bool, error)`; update the
  inline `parser:` closures in `option_parse.go` and the four `Parse` call sites
  (`mlrcli_parse.go`, `mlrcli_mlrrc.go`, `transformers/split.go`).
- Tedious but rote; also lets `pkg/terminals/completion` drop its
  "parsers call os.Exit" workaround (noted in its header comment).

### Phase 3 — streaming interface (cluster 3b) **← the load-bearing phase for #341/#440**
- `Transform` returns `error`; all ~55 verbs get `return nil` on happy paths.
- `ChainTransformer`/`runSingleTransformer` + `ChannelWriter` propagate errors via
  `chan error`; retire `dataProcessingErrorChannel chan bool`; `stream.Stream` selects
  and returns.
- Convert runtime exits in `tee.go`, `split.go`, `step.go`, `join_bucket_keeper.go`,
  and `put_or_filter.go` (the `ExecuteBeginBlocks`/`ExecuteMainBlock`/
  `ExecuteEndBlocks` err → exit sites) into returned errors.

### Phase 4 — sub-entrypoints (cluster 6)
- Per aux/terminal: usage funcs stop taking an `exitCode int` and exiting; inner I/O
  helpers return errors; each `main(args []string) int` returns codes to the
  dispatcher. `repl/entry.go`, `script/entry.go`, `regtest`, `hex/unhex/lecat/termcvt`.

### Phase 5 — DSL runtime (cluster 4)
- `hofs.go`/`udf.go` `...OrDie` helpers → error-Mlrvals per D3; delete
  `NoExitOnFunctionNotFound` special-casing once return-errors is the only behavior.
- Behavior note to document: HOF misuse now surfaces as `mlr: <msg>` via error
  propagation (same message, same exit 1 at top) rather than instant abort; REPL keeps
  its print-and-continue behavior for free.

## Verification

- `make dev` (fmt, build, unit + regression, docs) and `make lint` per phase.
- Exit-code checks: `mlr --version; echo $?` (0), unknown verb/flag (1), `mlr put
  'syntax error'` (1), `mlr join` with unreadable left file (1, per #2171), broken
  pipe `mlr cat big | head -1` (code unchanged), `--errors-json` output shape.
- `mlr repl` still print-and-continues on DSL errors.
- Grep gate per phase: `grep -rn os.Exit pkg/ | grep -v _test` shrinks to the
  documented keep-list; final state ≈ `entrypoint.go`, `logger.go`, `bifs/types.go`
  (asserting), `auxents.go`/`terminals.go` dispatchers.
