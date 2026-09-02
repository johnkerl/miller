# Plan: asof joins (new verb `join-asof`)

## Context

[Issue #173](https://github.com/johnkerl/miller/issues/173) opens with a
request for arbitrary join conditions (`join -c '$start_time <= $time &&
$end_time > $time'`), but the motivating use case — and every follow-up
comment on the thread — is specifically an **asof join**: for each record in
a stream (e.g. trades), attach fields from the reference-table record (e.g.
currency exchange rates) whose key is the closest one at-or-before the
stream record's key. The thread accumulates prior art as other tools grew
this exact feature:

- R data.table: *rolling joins* (`roll = TRUE`)
- pandas: [`merge_asof`](https://pandas.pydata.org/docs/reference/api/pandas.merge_asof.html)
  (with `direction=backward|forward|nearest`, `tolerance`,
  `allow_exact_matches`, `by=` equality keys)
- DuckDB 0.8.0+: SQL `ASOF JOIN`
- ClickHouse: `ASOF JOIN`
- Polars: `join_asof`

This plan proposes a dedicated verb, **`join-asof`**, with semantics modeled
on pandas `merge_asof` (the most fully-featured of the group), rather than
either (a) bolting asof flags onto the already-complex `join` verb
(`pkg/transformers/join.go` is 739 lines with two streaming modes whose
flag semantics — `-s`/`-u`, doubly-streaming — don't transfer cleanly to
asof matching), or (b) implementing the original generalized `-c`
DSL-condition join (see "Non-goals" below).

The canonical example from the issue:

```
# rates.csv          # trades.csv       # desired output of
time,rate            time               #   backward asof join
1,1.0                1                  time rate
4,1.1                2                  1    1.0
5,1.5                3                  2    1.0
7,2.0                4                  3    1.0
9,1.8                5                  4    1.1
10,1.7               6                  5    1.5
11,1.9               7                  6    1.5
                     8                  7    2.0
                     9                  8    2.0
                     10                 9    1.8
                                        10   1.7
```

## Naming

`join-asof` rather than `asof-join`:

- Sorts adjacent to `join` in the alphabetized verb list (`mlr help verbs`,
  the docs page, `aaa_transformer_table.go`), which aids discovery.
- Matches the operation-then-modifier order of polars `join_asof` and
  pandas `merge_asof`.

## Left/right terminology

Follow `mlr join`'s existing convention: the `-f` file is the **left**
(lookup/reference) table, ingested into memory; the main record stream on
the Miller argument list is the **right** (probe) side. Each right record
pairs with **at most one** left record.

Note this is the *opposite* of pandas' argument naming — in
`pd.merge_asof(trades, rates)` the probe side (trades) is "left" — but
consistency with `mlr join` (where `-f` = left, held in memory) matters
more here than consistency with pandas. Docs should include a sentence
calling this out for pandas users.

For the canonical example: `rates.csv` goes in `-f`; `trades.csv` is the
main stream.

## Proposed CLI signature

```
Usage: mlr join-asof [options]
Joins each record of the main input stream with the record from the left
file whose key is nearest, per the selected direction: the most recent
at-or-before key (backward, the default), the earliest at-or-after key
(forward), or the closest in either direction (nearest). This is the
"asof join" of pandas, polars, DuckDB, and ClickHouse, and the "rolling
join" of R data.table -- typically used to align timestamped data (e.g.
trades vs. sparse exchange-rate quotes).
```

Flags (OptionSpec table, following `joinOptions` in
`pkg/transformers/join.go`):

| Flag | Arg | Type | Description |
|---|---|---|---|
| `-f` | `{left file name}` | filename | Left/lookup file name for the join. Required. |
| `-j` | `{name}` | string | Asof-key field name for output; also the left and right key field name unless overridden by `-l`/`-r`. Required. Exactly one field name (not a comma-separated list). |
| `-l` | `{name}` | string | Asof-key field name in the left file; defaults to the `-j` value. |
| `-r` | `{name}` | string | Asof-key field name in the right (main-stream) records; defaults to the `-j` value. |
| `-g` | `{a,b,c}` | csv-list | Optional comma-separated equality-match field names (pandas `by=`): records pair only when these fields match exactly, with the asof rule applied within each group. |
| `--lg` | `{a,b,c}` | csv-list | Equality-match field names in the left file; defaults to `-g` values. |
| `--rg` | `{a,b,c}` | csv-list | Equality-match field names in the right records; defaults to `-g` values. |
| `--dir` | `{d}` | string | Match direction: `backward` (left key <= right key; default), `forward` (left key >= right key), or `nearest` (smallest absolute distance; ties resolve backward). |
| `--strict` | | bool | Exclude exact key matches (pandas `allow_exact_matches=False`): backward becomes strictly-less-than, forward strictly-greater-than, nearest excludes equal keys. |
| `--tolerance` | `{n}` | string | Maximum allowed absolute difference between right and left keys; pairs farther apart than this are treated as unmatched. Requires numeric keys (see semantics below). |
| `--asof-key-field` | `{name}` | string | If supplied, also emit the *matched left record's* key value into paired output records under this field name. Useful with `nearest`/`tolerance` to see what actually matched. |
| `--lk`, `--left-keep-field-names` | `{a,b,c}` | csv-list | Keep only these fields from the left file (asof-key and `-g` fields automatically included), as in `join --lk`. |
| `--lp` | `{text}` | string | Prefix for non-key output field names from the left file, as in `join --lp`. |
| `--rp` | `{text}` | string | Prefix for non-key output field names from the right records, as in `join --rp`. |
| `--np` | | bool | Do not emit paired records. |
| `--ul` | | bool | Emit unpaired records from the left file (those never matched by any right record), after end of stream. |
| `--ur` | | bool | Emit unpaired records from the right stream (those with no qualifying left match, e.g. before the first quote, beyond tolerance, or missing key fields). |
| `--prepipe`, `--prepipex` | `{command}` | string | As in `join`: shell command to prepipe the left file through. |

Plus the same left-file format-override passthrough as `join` (`-i`,
`--icsv`, `--ijson`, `--irs/--ifs/--ips`, `--implicit-csv-header`, etc.),
handled by the same `cli.FLAG_TABLE.Parse` fallthrough in the parse loop.

Deliberate differences from `join`:

- No `-s`/`-u` flags. There is only one mode (see "Streaming model"): the
  left file is ingested and sorted in memory; the right stream is fully
  streaming and **need not be sorted** (unlike pandas/polars, which require
  both sides sorted). A doubly-streaming sorted mode is deferred (see
  Non-goals).
- `-j`/`-l`/`-r` take a single field name, not a list — the asof key is one
  field by construction. Multi-field equality keys go in `-g`.

### Example invocations

```
# Canonical example from the issue:
mlr --csv join-asof -f rates.csv -j time trades.csv

# Different key names on each side; keep only the rate column from the left:
mlr --csv join-asof -f rates.csv -l quote_time -r trade_time -j time --lk rate trades.csv

# Multi-currency: exact match on currency, asof on time, within 60 seconds:
mlr --csv join-asof -f rates.csv -j time -g currency --tolerance 60 --ur trades.csv

# Forward-looking match, excluding same-timestamp quotes:
mlr --csv join-asof -f rates.csv -j time --dir forward --strict trades.csv
```

## Semantics

### Matching rule

For each right record:

1. Extract the `--rg` field values; if any is absent, the record is
   right-unpaired. Look up the left-record bucket with equal `--lg` values
   (whole-record-string equality on the joined group values, exactly as
   `join` buckets by `GetSelectedValuesJoined`). No bucket → right-unpaired.
2. Extract the `-r` asof-key value; absent → right-unpaired.
3. Binary-search the bucket's key-sorted left records for the candidate
   per `--dir`:
   - `backward`: largest left key `<=` right key (`<` under `--strict`).
   - `forward`: smallest left key `>=` right key (`>` under `--strict`).
   - `nearest`: whichever of the backward and forward candidates has the
     smaller absolute key distance; on an exact tie, the backward
     (earlier) candidate wins. Under `--strict`, equal-key candidates are
     excluded first.
4. If `--tolerance {n}` is given, reject the candidate when
   `abs(rightkey - leftkey) > n`. Both keys must be numeric for the
   difference to be computable; if either is not, the candidate is
   rejected (record becomes right-unpaired). `n` itself must parse as a
   non-negative number — CLI error otherwise.
5. No qualifying candidate → right-unpaired. Otherwise emit the merged
   record (unless `--np`) and mark that specific left record as paired
   (for `--ul` accounting).

### Key ordering / comparison

Keys are compared with Miller's natural-ordering comparator
(`mlrval.NaturalAscendingComparator`, `pkg/mlrval/mlrval_sort.go` — the
same collation as `sort -t`): numeric values compare numerically
(int/float mixed OK), strings compare lexically, and the ordering is a
total order across mixed types. This makes both integer/float epoch times
and ISO-8601 timestamp strings (which sort lexically = chronologically)
work without any type flags. `nearest`'s distance comparison and
`--tolerance` additionally need numeric subtraction, hence the numeric
requirement in step 4; `backward`/`forward` without tolerance work fine on
string keys.

### Duplicate left keys

If several left records in a bucket share the key selected by the
direction rule, the **last one in left-file input order** wins (pandas
behavior). Implementation: stable-sort the bucket by key, binary-search
for the rightmost record with the chosen backward key / leftmost with the
chosen forward key as appropriate — for backward the rightmost duplicate,
for forward also the *last* duplicate of the smallest qualifying key, per
pandas. Output is always at most one record per right record.

### Output record layout

Mirroring `join`'s `formAndEmitPairs`, a paired output record contains, in
order:

1. The asof-key field, named by `-j`, carrying the **right** record's key
   value (the probe time — matching pandas/DuckDB output and the issue's
   desired output).
2. The `-g` equality fields under their `-g` output names (values equal on
   both sides by construction).
3. If `--asof-key-field {name}` was given: the matched left record's key
   value under `{name}`.
4. Left-record fields not already added (i.e. not the left key, not `--lg`
   fields), each prefixed with `--lp` if given, filtered by `--lk` if
   given.
5. Right-record fields not already added, each prefixed with `--rp`.

Unpaired records (under `--ul`/`--ur`) get the same key-rename-to-output-
names and prefixing treatment as `join`'s `transformUnpairedRecord`, so
column names stay consistent for downstream `unsparsify` etc.

For the canonical example this yields exactly `time,rate` per the issue.

### Heterogeneity and edge cases

- Left records missing the `-l` key or any `--lg` field: never matchable;
  emitted only under `--ul`.
- Empty left file: every right record is unpaired.
- Left file need **not** be pre-sorted; it is sorted internally after
  ingest. Right stream order is irrelevant to correctness (each record is
  matched independently), and input order of right records is preserved
  in the output.
- `--ul` output ordering: as in `join`'s half-streaming mode, left
  unpairables are emitted at end of stream, in left-file order.
- All emit flags unset (`--np` without `--ul`/`--ur`): CLI error, same
  message pattern as `join`.

## Streaming model and performance

Single mode, analogous to `join`'s default half-streaming mode:

- Ingest the entire left file at construction/first-record time (reusing
  `join`'s left-file reader logic), bucketing by `--lg` values.
- Per bucket: stable-sort records by key once (`sort.SliceStable`,
  `O(L log L)`), storing a parallel slice of key mlrvals.
- Per right record: hash lookup of bucket + binary search (`sort.Search`),
  `O(log L)`. The right stream is never buffered.

Memory is proportional to the left file only — appropriate since the
lookup table (rates) is typically much smaller than the stream (trades).

## Implementation sketch

1. **New file `pkg/transformers/join_asof.go`** (~450 lines), containing:
   - `verbNameJoinAsof = "join-asof"`, `joinAsofOptions []OptionSpec`,
     `JoinAsofSetup TransformerSetup`, usage func, ParseCLI func —
     structured exactly like `join.go`'s.
   - `tJoinAsofOptions` struct (leftFileName, key/group field names,
     direction enum, strict bool, tolerance *mlrval, prefixes, keep-set,
     emit flags, prepipe, `joinFlagOptions cli.TOptions`).
   - `TransformerJoinAsof` with lazy first-record ingest (same
     `tr.ingested` pattern as `transformHalfStreaming`), bucket map
     `map[string]*asofBucket` where

     ```go
     type asofBucket struct {
         recordsAndContexts []*types.RecordAndContext // stable-sorted by key
         keys               []*mlrval.Mlrval          // parallel slice
         paired             []bool                    // per-record, for --ul
     }
     ```

   - Match/merge logic per the semantics above. The merge loop is a
     single-left-record specialization of `join.go`'s `formAndEmitPairs`.
2. **Shared left-file ingestion.** `join.go`'s `ingestLeftFile`
   (`join.go:505-587`) is coupled to `TransformerJoin`'s bucketing.
   Extract the reader-plumbing core (create reader, channels, drain to a
   `[]*types.RecordAndContext`, including the error-channel-before-EOS
   check) into a helper in `pkg/transformers/utils/` (e.g.
   `utils.IngestFileToSlice(fileName string, readerOpts
   *cli.TReaderOptions) ([]*types.RecordAndContext, error)`); both `join`
   and `join-asof` then do their own bucketing over the slice. If the
   refactor proves riskier than expected, fall back to a private copy in
   `join_asof.go` and file a follow-up to unify.
3. **Registration**: add `JoinAsofSetup` to `STANDARD_TRANSFORMERS` in
   `pkg/transformers/aaa_transformer_table.go` immediately after
   `JoinSetup` (line 43) — the table is alphabetical. `mlr help verbs`,
   `mlr --usage-all-verbs`, and manpage content are generated from this
   table and the usage func, so they pick the verb up automatically.
4. **CLI validation** in ParseCLI: require `-f` and `-j`; reject
   comma-containing `-j`/`-l`/`-r` values with a message pointing to `-g`
   for equality keys; validate `--dir` ∈ {backward, forward, nearest};
   validate `--tolerance` parses as a non-negative number; require at
   least one emit flag effective; same `-l`/`-r` (and `--lg`/`--rg`)
   defaulting rules as `join`'s.

## Test cases

### Test inputs

New files under `test/input/`:

- `asof-rates.csv` — the issue's 7-row rates table.
- `asof-trades.csv` — the issue's 10-row trades table.
- `asof-rates-multi.csv` — rates with a `currency` column (e.g. CAD and
  EUR interleaved, distinct gap structures per currency) for `-g` tests.
- `asof-trades-multi.csv` — trades with a `currency` column, including one
  currency absent from the rates file.
- `asof-rates-dup.csv` — rates with a repeated `time` value (two quotes at
  t=5) for the duplicate-key rule.
- `asof-rates-iso.csv` / `asof-trades-iso.csv` — ISO-8601 string
  timestamps.
- `asof-rates.json` — JSON twin of `asof-rates.csv` for the mixed-format
  case.

### Regression cases (`test/cases/verb-join-asof/00NN/{cmd,expout,experr}`)

Happy paths:

1. **Backward default** (the issue's canonical example):
   `mlr --icsv --opprint join-asof -f test/input/asof-rates.csv -j time
   test/input/asof-trades.csv` — expout is exactly the issue's desired
   table:

   ```
   time rate
   1    1.0
   2    1.0
   3    1.0
   4    1.1
   5    1.5
   6    1.5
   7    2.0
   8    2.0
   9    1.8
   10   1.7
   ```

2. **Forward**: `--dir forward` — expected rates
   `1.0, 1.1, 1.1, 1.1, 1.5, 2.0, 2.0, 1.8, 1.8, 1.7` for trades 1–10.
3. **Nearest with tie-break**: `--dir nearest` — trades 6 and 8 are
   equidistant between quotes and must take the backward value
   (`1.5, 2.0` respectively); full expected column
   `1.0, 1.0, 1.1, 1.1, 1.5, 1.5, 2.0, 2.0, 1.8, 1.7`.
4. **Strict backward**: `--strict --ur` — trade 1 has no strictly-earlier
   quote and comes out unpaired; trade 4 gets 1.0 (not 1.1), trade 5 gets
   1.1, etc.
5. **Tolerance**: `--tolerance 1 --ur` — trade 3 (distance 2 from quote
   at t=1) is unpaired; all others pair.
6. **Group-by**: `-g currency` on the multi-currency inputs — asof applies
   within each currency; trades in the currency absent from rates are
   unpaired (verify both with and without `--ur`).
7. **Duplicate left keys**: rates-dup input, trade at t=5 and t=6 — both
   must take the *last* t=5 quote in file order.
8. **String keys**: ISO-8601 inputs, backward join — verifies lexical
   collation path.
9. **Key-name overrides**: `-j time -l quote_time -r trade_time` (rename
   variants of the inputs, or reuse via `rename` in the cmd pipeline) —
   output key column is named `time`.
10. **Prefixes and collisions**: both sides carrying a same-named non-key
    field, with `--lp l_ --rp r_`, and again with no prefixes (last-put
    wins per `join` semantics — pin whichever behavior `formAndEmitPairs`
    ordering produces).
11. **`--lk`**: keep only a subset of left fields; asof-key auto-included.
12. **`--asof-key-field matched_time`**: paired records carry the left
    key; combined with `nearest` so the matched key differs from the probe
    key.
13. **`--np --ul --ur`**: only-unpaired output, both sides, checking
    end-of-stream `--ul` emission order.
14. **`--ul` marking granularity**: with the canonical inputs, quote at
    t=11 is never matched and must appear under `--ul`; all other quotes
    must not.
15. **Right records missing the key field**: heterogeneous dkvp input
    where some records lack `time` — unpaired, surfaced only with `--ur`.
16. **Empty left file**: all right records unpaired; no crash.
17. **Mixed formats**: `--ijson` after the verb for a JSON left file with
    CSV main input (mirroring `test/cases/verb-join-mixed-format`).
18. **Unsorted left file**: rates file shuffled on disk — output identical
    to case 1 (internal sort).
19. **In-chain use**: `mlr ... cat then join-asof ...` to confirm normal
    then-chain behavior.

Error cases (expout empty, experr pinned):

20. Missing `-f`; missing `-j`.
21. `-j a,b` (comma in asof key) — error mentioning `-g`.
22. `--dir sideways` — error listing valid directions.
23. `--tolerance banana` and `--tolerance -1` — parse/negativity errors.
24. `--np` alone (no emit flags effective).

### Unit tests

- `pkg/transformers/join_asof_test.go` (or in `utils` if the search is
  factored there): table-driven tests for the candidate-selection function
  — direction × strict × tolerance × duplicate-keys × empty-bucket ×
  string-keys, isolated from record plumbing. This is where the fiddly
  binary-search boundary conditions (first element, last element,
  all-equal keys) get exhaustive coverage.
- A case in `aaa_transformer_json_test.go`-style if applicable for CLI
  parse of the new verb (match how neighboring verbs are covered).

## Documentation

- `docs/src/reference-verbs.md.in`: new `## join-asof` section directly
  after `## join` (~line 616), with GENMD-RUN-COMMAND blocks for
  `mlr join-asof --help` and worked examples. Example data files
  `docs/src/data/asof-rates.csv` and `docs/src/data/asof-trades.csv`
  (docs GENMD commands run relative to `docs/src`). Include: the
  canonical backward example, `--dir nearest`, `-g currency`, and
  `--tolerance` with `--ur`.
- Cross-references: add a pointer from the `## join` docs section ("for
  nearest-key rather than equal-key matching, see join-asof"), and a
  sentence in `join`'s usage string is *not* needed — keep usage strings
  self-contained.
- A short paragraph for pandas/SQL users naming the equivalents
  (`merge_asof`, `ASOF JOIN`, rolling joins) — these are the search terms
  people arrive with, per the issue thread — and the left/right naming
  caveat from "Left/right terminology" above.
- Rebuild with `make -C docs/src forcebuild`; `make dev` regenerates the
  manpage and verb listings.

## Non-goals (this pass)

- **Generalized DSL join conditions** (`join -c '<expr>'`, the issue's
  original framing). Arbitrary conditions defeat both hashing and binary
  search, forcing an O(N·M) cross-product scan with DSL evaluation per
  pair, and require wiring a CST evaluator with a two-record binding
  convention (left vs. right field namespacing) into a verb. Every
  concrete use case in the thread is served by asof semantics. If demand
  remains after `join-asof` ships, that's a separate plan; the CLI surface
  proposed here doesn't foreclose it.
- **Doubly-streaming sorted mode** (`join -s` analog) for left files too
  large for memory. The bucket-keeper machinery doesn't transfer directly
  (asof needs lookback across bucket boundaries); defer until someone
  asks.
- **DSL-level asof lookup function** (e.g. a map-backed
  `asof_lookup(map, key)` BIF). Possibly useful someday; out of scope.
- **Interpolation** (e.g. linear interpolation between bracketing quotes,
  merge-ordered-style). data.table and pandas keep this separate from
  asof joins; so should Miller.

## Verification

1. `make build`; manual run of the canonical example and a spot-check of
   `--dir nearest` / `--tolerance` against pandas `merge_asof` output on
   the same data.
2. `make check` (unit + regression, including the new
   `test/cases/verb-join-asof/` cases).
3. `make dev` — confirms docs regeneration picks up the new verb and no
   existing goldens shift (in particular, `mlr help verbs` output and
   manpage goldens will change and need regenerating in the same commit).
4. `make lint` before pushing.
