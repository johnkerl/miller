# Plan: auto-infer input format from file extension

Feature request: [issue #1188](https://github.com/johnkerl/miller/issues/1188) — given
`mlr ... mydata.csv`, default to `--icsv` without the user typing it; likewise `.tsv` →
`--itsv`, etc.

This is at heart a refactor of `pkg/stream/stream.go` and the record-reader layer, not a
CLI-flag tweak. Today a single record-reader is constructed once and handed the full list
of file names; with this feature, the reader (and its format-dependent option defaults)
must be chosen **per input file**.

Scope: **input side only.** Output-format inference is discussed under Open Questions but
is recommended out of scope for v1, with one carve-out for in-place mode (see complication
C9, which is a data-loss footgun if ignored).

## Current architecture (survey)

Where readers are constructed — one per process, not per file:

- `pkg/stream/stream.go:51` — `input.Create(&options.ReaderOptions, ...)` once, then
  `go recordReader.Read(fileNames, ...)` at `stream.go:88`.
- `pkg/input/record_reader_factory.go:9` — `input.Create` switches on
  `readerOptions.InputFileFormat` (csv, csvlite, dkvp, dkvpx, json, yaml, nidx,
  markdown, pprint, tsv, xtab, dcf, recutils, gen).
- Other `input.Create` call sites, each needing its own treatment:
  - `pkg/transformers/join.go:510` and `pkg/transformers/utils/join_bucket_keeper.go:163`
    — the join verb's left file, with its own `joinFlagOptions.ReaderOptions`.
  - `pkg/terminals/repl/session.go:49` — REPL.
  - `pkg/terminals/script/runner.go:24` — `mlr script` runner.

How readers consume files — each reader owns the whole file loop:

- `IRecordReader.Read(filenames []string, initialContext types.Context, ...)`
  (`pkg/input/record_reader.go:14`) takes the *list*; every concrete reader duplicates
  the same boilerplate: `filenames == nil` → no input (`mlr -n`); `len == 0` → stdin via
  `lib.OpenStdin`; else loop with `lib.OpenFileForRead`, calling a per-file
  `processHandle(handle, filename, &context, ...)` that resets per-file reader state and
  calls `context.UpdateForStartOfFile` (see e.g. `pkg/input/record_reader_csv.go:54-108`).
  Each `Read` sends exactly one end-of-stream marker after all files.
- AWK-ish variables: `types.Context.UpdateForStartOfFile` (`pkg/types/context.go:132`)
  increments FILENUM and resets FNR; NR accumulates across files. This continuity must
  survive the move to per-file readers.
- `downstreamDoneChannel` (mlr head fast-exit): readers poll it per batch
  (`record_reader_csv.go:171`, `line_reader.go:202`). The signal is a **one-shot buffered
  send**; whichever per-file scanner consumes it stops, and the file loop must remember
  done-ness so it doesn't open remaining files. (Today the loop just proceeds to the next
  file, whose scanner will never see the already-consumed signal — a latent inefficiency
  the refactor should fix, not replicate.)

Where format-dependent defaults are applied — once, at CLI-parse time:

- `cli.FinalizeReaderOptions` (`pkg/cli/option_parse.go:31`) fills in IFS/IPS/IRS and
  AllowRepeatIFS defaults keyed by `InputFileFormat` (with a special NIDX
  whitespace-regex case), and un-hexes/un-backslashes the separator strings **by mutating
  the one shared `TReaderOptions`**. Per-file formats mean per-format finalization; the
  un-escaping steps must not be re-applied to already-processed values.
- Reader constructors also validate format-specific constraints at construction time
  (e.g. `NewRecordReaderCSV` at `record_reader_csv.go:34-43`: single-char IFS, IRS
  restrictions, comment-string length) and cache derived state (`ifs0`).

File opening / compression:

- `lib.OpenFileForRead` (`pkg/lib/file_readers.go:49`): prepipe → popen; else encoding
  flag; else suffix-sniffed decompression for `.gz`, `.bz2`, `.zst`
  (`openEncodedHandleForRead`). `lib.PathToHandle` also supports `http://`, `https://`,
  `file://` URLs.
- In-place mode (`pkg/entrypoint/entrypoint.go:127-207`) already runs one full
  `Stream()` per file, re-parsing the command line each time, and infers input
  *compression* from the file name (`lib.FindInputEncoding`).

## Proposed design

### 1. Extension→format inference (pure, name-based — no content sniffing)

New function, e.g. `input.InferFormatFromFileName(path string) (format string, ok bool)`:

1. If the path is a URL (`http://`, `https://`, `file://`), strip scheme and any
   `?query`/`#fragment` before looking at the suffix.
2. Strip one trailing compression suffix (`.gz`, `.bz2`, `.zst`) — mirrors
   `openEncodedHandleForRead` — so `data.csv.gz` infers csv.
3. Map the remaining extension, case-insensitively:

   | extension | format |
   |---|---|
   | `.csv` | `csv` |
   | `.tsv` | `tsv` |
   | `.json`, `.jsonl`, `.ndjson` | `json` (the JSON reader already handles both) |
   | `.yaml`, `.yml` | `yaml` |
   | `.md`, `.markdown` | `markdown` |
   | `.dkvp` | `dkvp` |
   | `.nidx` | `nidx` |
   | `.xtab` | `xtab` |
   | `.pprint` | `pprint` |
   | `.rec` | `recutils` |
   | `.dcf` | `dcf` |

   Deliberately unmapped: `.txt`, `.dat`, `.log`, and anything else ambiguous — those
   return `ok=false` and take the fallback (below). `.csv` maps to full `csv`, not
   `csvlite` (users who want csvlite say so explicitly).

Because inference is name-based only, **all per-file formats can be resolved up front**,
before any goroutine starts — so option-validation errors surface before any output is
produced, not mid-stream after file 3 of 7.

### 2. Policy: stdin and fallback

- **stdin: no inference** (per maintainer). `len(filenames)==0` uses the fallback format.
- **Fallback** for stdin and unmapped extensions: the default format, `dkvp`. (Open
  question Q3 discusses erroring instead.)
- `mlr -n` (nil filenames): unaffected.

### 3. CLI surface: opt-in flag, not a default flip

- Accept `auto` as an `InputFileFormat` value: `-i auto`, plus a dedicated `--iauto` flag
  in `FileFormatFlagSection` (`pkg/cli/option_parse.go:844`). Last-one-wins with other
  format flags, per existing CLI semantics: `--icsv --iauto` means auto.
- `.mlrrc` gives users a "make it my default" path (`iauto` on a line by itself), since
  .mlrrc lines are flags without leading dashes. This is the adoption path in lieu of
  changing the built-in default, which would silently break existing scripts that rely on
  `mlr cat foo.csv` parsing as DKVP (see Q1).
- `-o auto` / `--io auto` are **errors** in v1 (output can't be name-inferred; it usually
  goes to stdout). Revisit under Q2.

### 4. The stream.go / record-reader refactor

Two-layer split — this is the bulk of the work, and is a worthwhile cleanup even
independent of the feature (it deletes ~13 copies of the same stdin/loop/open
boilerplate):

- **Per-file readers**: each concrete reader keeps its constructor and its
  `processHandle`-shaped method; the `Read(filenames, ...)` file loop is deleted from
  all of them. New narrower interface, roughly:

  ```go
  type IFileRecordReader interface {
      ProcessHandle(handle io.Reader, filename string, context *types.Context,
          readerChannel chan<- []*types.RecordAndContext,
          errorChannel chan error,
          downstreamDoneChannel <-chan bool)
  }
  ```

  (Most readers already have exactly this method; the change is mostly mechanical.
  Signature detail to settle during implementation: some readers need to report
  "downstream done" back to the caller so the driver stops opening further files —
  probably a `bool` return.)

- **Driver**: one new `FileStreamReader` in `pkg/input` implementing the existing
  `IRecordReader` interface, owning: the nil/stdin/file-list branching, per-file
  `OpenFileForRead`/`OpenStdin` + close, `context.UpdateForStartOfFile`, the
  downstream-done latch across files, and the single end-of-stream marker.

  It is constructed with a pre-resolved plan: `[]struct{ fileName string; reader
  IFileRecordReader }` (plus a stdin entry when applicable). In non-auto mode that's the
  same reader for every file — behavior identical to today. In auto mode it's built by
  running inference per file at setup time.

- **`gen` pseudo-reader** (`pseudo_reader_gen.go`) reads no files; it keeps implementing
  `IRecordReader` directly and bypasses the driver.

- **Factory**: `input.Create` grows a companion, e.g.
  `input.CreateForFileNames(readerOptions, recordsPerBatch, fileNames) (IRecordReader, error)`,
  which handles the auto/non-auto/gen dispatch and becomes what `stream.go`, join, REPL,
  and script-runner call.

### 5. Per-format reader options

Refactor `FinalizeReaderOptions` so format-dependent defaulting is a pure derivation
rather than a one-shot mutation:

- Split into (a) a once-only un-escaping step on user-supplied separator strings
  (unhex/unbackslash — must not run twice), done at CLI-parse time as now; and (b)
  `deriveReaderOptionsForFormat(base *TReaderOptions, format string) (*TReaderOptions, error)`
  which returns a **copy** with IFS/IPS/IRS/AllowRepeatIFS defaults (and the NIDX
  whitespace-regex special case) applied for that format, honoring the
  `*WasSpecified` booleans so explicit `--ifs` etc. still win for every inferred format.
- In non-auto mode, (b) is called once with the single format — same net behavior as
  today. In auto mode, called once per distinct inferred format (cache in a
  `map[string]*TReaderOptions`), then one concrete reader constructed per distinct
  format (readers already reset per-file state in `processHandle`, so sharing one reader
  across same-format files is safe and matches current cross-file behavior).
- When `InputFileFormat == "auto"` reaches finalize-time, skip the format-keyed lookups
  (they'd fail on the `defaultFSes` map) — defaults are applied per derived format
  instead.

## Complications inventory

Ones already flagged by the maintainer:

- **C1 — one reader for all files** → per-file (per-format) construction; addressed by §4.
- **C2 — stdin** → no inference, fallback format; addressed by §2.

Additional ones surfaced by this survey:

- **C3 — format-dependent option defaults are baked in at CLI-parse time.**
  IFS/IPS/IRS/AllowRepeatIFS defaults differ per format and are applied by mutating the
  single shared `TReaderOptions` (`option_parse.go:31-84`). `mlr --iauto cat a.csv b.nidx`
  needs comma-IFS for one file and whitespace-regex-IFS for the other. Needs the
  derive-per-format refactor (§5), including not double-applying unhex/unbackslash.
- **C4 — constructor-time validation and cached state.** Reader constructors validate
  and cache options (`record_reader_csv.go:34-51`). Mitigation: resolve formats and
  construct all readers eagerly at stream setup (possible because inference is
  name-based), so `mlr --iauto --ifs ';;' cat a.csv` fails before any records flow.
- **C5 — AWK-variable continuity across heterogeneous readers.** FILENUM/NR must keep
  accumulating when consecutive files use different readers. Solved by the driver owning
  one `types.Context` and passing it into each per-file read; also exactly one
  end-of-stream marker, sent by the driver.
- **C6 — `downstreamDoneChannel` is a one-shot signal.** Once a per-file scanner consumes
  it, later files can't see it. The driver must latch done-ness (via the per-file return,
  C4's interface note) and stop opening subsequent files. Today's per-reader loops appear
  to keep reading subsequent files after `head` is satisfied — the refactor should fix
  this, and it's worth a regression test (`mlr head -n 1 big1.csv big2.csv` should not
  read big2 to completion).
- **C7 — compressed and URL inputs.** `data.csv.gz` must infer csv (strip compression
  suffix, mirroring `openEncodedHandleForRead`); URLs need scheme/query stripping before
  suffix inspection. `--prepipe 'unzip -qc' foo.zip` → `.zip` unmapped → fallback (fine;
  prepipe users can state the format explicitly).
- **C8 — the join verb reads its own file.** `join -f left.csv` via
  `ingestLeftFile` (`join.go:506`) and the half-streaming
  `join_bucket_keeper.go:163` use `joinFlagOptions.ReaderOptions`, which inherit main
  reader options unless overridden by join's own `-i`. If the inherited/derived format is
  `auto`, these paths must run the same resolve-then-create helper (trivial: single known
  file name, resolvable up front). Without this, `-i auto` at main level would hit
  `input.Create`'s `default:` error ("input file format \"auto\" not found") inside join.
- **C9 — in-place mode (`mlr -I`) is a data-loss footgun with auto.** `-I` writes output
  back over the input file using the *output* format. `mlr -I --iauto put ... foo.csv`
  with default DKVP output would silently rewrite a CSV file as DKVP. v1 must do one of:
  (a) error on `-I` + auto unless an explicit output format is given, or (b) per file, set
  the writer format to the inferred reader format when no explicit `-o` was given —
  natural since `processFileInPlace` already re-parses options per file
  (`entrypoint.go:161`). Recommend (b); (a) is an acceptable stopgap. Either way this
  must not ship as "whatever falls out".
- **C10 — REPL and `mlr script`.** `repl/session.go:49` constructs its reader at session
  start, before any `:open file` — with auto, resolution has to happen per opened file.
  Simplest v1: REPL/script reject or ignore `auto` with a clear message, or resolve at
  `:open` time using the same helper. Decide during implementation; don't leave it
  crashing on the factory `default:` case.
- **C11 — mixed formats in one run are now easy to trigger.** Heterogeneous records are
  native to Miller, so `mlr --iauto cat a.csv b.json` "just works", but users will see
  format-specific side effects (e.g. CSV output emitting new header blocks on schema
  change). Docs should show a mixed-format example. Also note format-specific input flags
  (`--allow-ragged-csv-input`, `--csv-trim-leading-space`, implicit-header, etc.) apply
  whenever the *inferred* format is csv/tsv — harmless for other formats, worth a doc
  sentence.
- **C12 — surprise on the output side.** `mlr --iauto cat data.csv` prints DKVP to
  stdout. Harmless (visible immediately) but guaranteed to generate "it doesn't work"
  reports from the very users #1188 is for. See Q2.
- **C13 — factory/validation error paths.** `"auto"` must be handled everywhere
  `InputFileFormat` is switch/map-keyed: `input.Create` default case,
  `FinalizeReaderOptions` map lookups, `pkg/cli/flatten_unflatten.go:99` (auto-flatten
  decides based on input format — with auto, the per-format derived options must carry
  the resolved format so flatten/unflatten heuristics see `csv`, not `auto`; audit other
  `InputFileFormat` consumers with `grep -rn InputFileFormat pkg/`).
- **C14 — the `dkvpx`/fixed-width/barred-pprint variants.** Inference only ever selects
  canonical formats; variant selection (`BarredPprintInput`, `FixedWidthSpec`, `dkvpx`)
  stays explicit-flag-only. No `.pprint`-barred inference.
- **C15 — case/AV edge cases in names.** Uppercase extensions (`.CSV`), files with no
  extension, dotfiles (`.csv` as an entire filename — treat as extensionless), a literal
  `-` filename (Miller doesn't special-case it today; keep it that way), Windows path
  separators. All belong in the inference unit tests.

## Phased implementation

Phases are separately mergeable, each leaving `make check` green.

1. **Reader-loop extraction (pure refactor, no behavior change).**
   Introduce `IFileRecordReader` + `FileStreamReader` driver; delete the per-reader file
   loops; keep `input.Create` signature; fix the C6 done-latch as part of the driver.
   This is the risky/bulky phase — all 13 readers touched mechanically. Regression suite
   is the safety net; add the C6 test here.
2. **Inference function + options derivation.**
   `InferFormatFromFileName` with unit tests (C7, C15 cases);
   `deriveReaderOptionsForFormat` refactor of `FinalizeReaderOptions` with unit tests
   proving explicit `--ifs/--ips/--irs` still override per derived format.
3. **Wire up `-i auto` / `--iauto` for the main stream.**
   Flag-table entry, `"auto"` handling at all C13 sites, eager per-file resolution in
   stream setup, fallback policy (§2). Regression cases: per-extension inference,
   mixed formats, `.csv.gz`, unmapped extension → dkvp, stdin → dkvp, explicit separator
   overrides under auto.
4. **Secondary consumers: join, in-place, REPL/script.**
   C8 (join left file), C9 (in-place policy — implement (b) or (a)), C10 (REPL/script).
   Regression cases for each, including the C9 "must not rewrite csv as dkvp" case.
5. **Docs.**
   `docs/src/file-formats.md.in` section on auto-inference (extension table, stdin/
   fallback rules, .mlrrc adoption tip, mixed-format example); flag help text (feeds the
   auto-generated `reference-main-flag-list`); `mlr help` topics; man page regen via
   `make dev`.

## Open questions for the maintainer

- **Q1 — default-on vs opt-in.** #1188 asks for this as *default* behavior. That flips
  parsing of `mlr cat foo.csv` from DKVP to CSV — behavior-breaking for scripts (however
  few) that depend on it. Recommendation: ship opt-in (`--iauto`, .mlrrc-able) now;
  consider flipping the default in a major release after the machinery has soaked.
- **Q2 — output-side inference.** Options: (a) none (v1 recommendation, minus the C9
  in-place carve-out); (b) `--auto` convenience flag = infer input per file *and*, when
  no explicit output format was given and all inputs infer to a single common format, use
  that for output too (resolvable up front since inference is name-based; must define
  behavior for mixed inputs — probably fall back to dkvp or error). (b) is what
  #1188-style users likely actually want day-to-day; fine as a fast-follow.
- **Q3 — unmapped extension under auto: fallback to dkvp, or error?** Fallback is
  forgiving and matches the stdin story; erroring is more predictable ("you asked for
  auto and I can't tell what `.dat` is"). Recommendation: fallback + document; a strict
  variant can come later if requested.
- **Q4 — reuse one reader per distinct format vs one per file.** Plan assumes per-format
  reuse (matches today's cross-file behavior exactly, since readers reset state per
  file). Per-file construction is marginally simpler to reason about but re-runs
  validation redundantly. Low stakes either way; decide in phase 3.
