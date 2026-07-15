# Plan: add GNU recutils (`.rec`) as a Miller I/O format

## Context

[Issue #378](https://github.com/johnkerl/miller/issues/378) requests support for
[GNU recutils](https://www.gnu.org/software/recutils/manual/index.html), a
human-editable, text-based record format: records are `FieldName: Value` lines,
records are separated by blank lines, and Miller's own docs already describe XTAB
as "perhaps most useful for looking at very wide... data" in a similar spirit. The
issue reporter noted the format's kinship with Miller's philosophy and with XTAB,
but explicitly left the integration decision to the maintainer.

Investigation found that Miller already has an even closer structural precedent
than XTAB: **DCF (Debian control file)** format
(`pkg/input/record_reader_dcf.go`, `pkg/output/record_writer_dcf.go`), which is
also `Key: value` lines with blank-line-separated paragraphs and continuation-line
folding. DCF's registration footprint (factory switches, `separators.go` maps,
`option_parse.go` flags, docs, regression tests) is the direct template for
wiring up recutils. DCF's *parsing*, however, delegates to an external
Debian-specific library and uses RFC822-style (leading-space) continuation —
wrong syntax for recutils (which uses `+`-prefixed and backslash-newline
continuation) — so the recutils reader needs to be hand-rolled, using **XTAB's**
blank-line stanza-scanning goroutine (`channelizedStanzaScanner` in
`pkg/input/record_reader_xtab.go`) as the structural template instead.

Scope decision: recutils' "record descriptor" (`%rec:`-prefixed) mechanism gives
the format a schema/constraint/type-checking layer (mandatory fields, types,
auto-increment keys, foreign-key-style references) approaching a lightweight
database. That layer is **out of scope** — Miller is a schema-less stream
processor. `%rec`-prefixed records are parsed with the exact same generic
logic as any other record (their fields, e.g. `%rec`, `%mandatory`, `%type`,
just become ordinary Miller fields with no special interpretation). This
mirrors how the existing DCF reader implements Debian control-file *syntax*
without any package-relationship *semantics*.

Two naming/behavior questions were resolved with the user:
- **Format name / flags**: `recutils`, not the shorter `rec` — flags are
  `--irecutils`, `--orecutils`, `--recutils`, matching the issue's own
  terminology for discoverability. (Not added to the `--c2j`-style
  format-conversion keystroke-saver matrix, following DCF's own precedent of
  being excluded from that 10x10 table.)
- **Malformed-line handling**: hard error (no leniency flag) for a data line
  missing the `": "` separator, or a `+`-continuation line with no preceding
  field in the stanza. This matches recutils' own spec (the separator is
  exactly colon-space) and keeps v1 scope small.

## Format semantics to implement

- Records = stanzas separated by one or more blank lines (leading/trailing
  blank lines ignored) — identical boundary logic to XTAB's.
- Each field is one line, `FieldName: Value`; split on the **first**
  occurrence of `": "` only (a value may legitimately contain `: ` later in
  the string).
- Two continuation mechanisms, both processed *before* colon-splitting, in
  this order:
  1. **Backslash-newline**: a line whose last character is `\` is joined
     directly (no separator) with the next physical line.
  2. **`+`-continuation**: a line starting with `+` continues the
     *previous field's* value with an embedded `\n`; a single leading space
     after the `+` is stripped (`+ text` → continuation text `text`); a bare
     `+` embeds an empty line.
- Comment lines (`#`-prefixed) reuse Miller's existing generic
  `readerOptions.CommentHandling`/`CommentString` machinery — same
  before-stanza-logic prefix check XTAB already does. No new option fields
  needed.
- `%rec`-prefixed descriptor records: no special handling (see Scope decision
  above) — parsed like any other record.
- PS is fixed at `": "`, RS is fixed at blank-line, FS is not applicable —
  none of these are configurable via `--ifs`/`--ips`/etc., matching DCF's
  "N/A" treatment. **No new fields needed in `TReaderOptions`/`TWriterOptions`**
  (`pkg/cli/option_types.go`) — this format needs zero format-specific option
  state, same as DCF.
- Known, documented limitation: a value whose last physical line ends in a
  literal `\` is ambiguous on round-trip (indistinguishable from an
  in-progress backslash-continuation) — this is an inherent ambiguity in the
  recutils format itself (no in-value backslash-escaping exists), not a gap
  specific to this implementation. Document it in `file-formats.md.in` and
  pin the (lossy) behavior with a regression test rather than inventing
  non-standard escaping.

## Implementation

### Reader: `pkg/input/record_reader_rec.go` (new)

Copy the outer `Read`/`processHandle` boilerplate from
`record_reader_dcf.go` (open stdin/files, batch onto `readerChannel`, respect
`downstreamDoneChannel`). For stanza-scanning, copy-and-adapt
`channelizedStanzaScanner` from `record_reader_xtab.go` — same
`tStanza{dataLines, commentLines}` blank-line-boundary/comment-handling logic,
but call `NewLineReader(handle, "\n")` with a fixed `"\n"` (not
`readerOptions.IFS`, since RS isn't configurable here unlike XTAB).

Replace XTAB's whitespace-run `pairSplitter` with a new
`recordFromRECLines(stanza []string) (*mlrval.Mlrmap, error)` implementing the
three-pass continuation/split algorithm above:
1. Backslash-join pass over raw lines.
2. `+`-continuation fold pass (error if a `+` line has no preceding field yet
   in the stanza).
3. First-`": "` split pass per resulting logical line (error if no `": "` is
   found).

Both error cases propagate via `errorChannel`, consistent with how other
readers fail on structurally malformed input (e.g. CSV's ragged-row default).

### Writer: `pkg/output/record_writer_rec.go` (new)

Copy `record_writer_dcf.go` near-verbatim with two changes:
1. Drop DCF's hardcoded array-field comma-joining (`dcfValueString`) — use
   plain `mv.String()` instead, like XTAB's writer. recutils has no special
   list-valued fields; nested/array values go through Miller's standard
   auto-flatten mechanism, since recutils is a non-nestable format (no change
   needed to `pkg/cli/flatten_unflatten.go`'s `isNestable()` — the existing
   default already does the right thing, same as XTAB).
2. Fold embedded `\n` in a value using `"+ "`-prefixed continuation lines
   (recutils convention), not DCF's single-leading-space convention.

Blank line written after every record (no `onFirst` guard needed, matching
DCF's writer).

### Registration (mirror every DCF touch point)

- `pkg/input/record_reader_factory.go`: add `case "recutils":` →
  `NewRecordReaderREC`.
- `pkg/output/record_writer_factory.go`: add `case "recutils":` →
  `NewRecordWriterREC`.
- `pkg/cli/separators.go`: add `"recutils"` entries to all four maps
  (`defaultFSes`, `defaultPSes`, `defaultRSes`, `defaultAllowRepeatIFSes`),
  identical shape to the existing `"dcf"` entries (`"N/A"` / `"N/A"` / `"N/A"`
  / `false`). This is also the master legal-format-name list
  (`GetFileFormatNames()`), so this step alone makes `-i recutils`, `-o
  recutils`, `--io recutils` work.
- `pkg/cli/option_parse.go` (`FileFormatFlagSection`): add `--irecutils`,
  `--orecutils`, `--recutils` flags, copied from the `--idcf`/`--odcf`/`--dcf`
  blocks (~lines 888, 1142, 1340). Update
  `FormatConversionKeystrokeSaverPrintInfo`'s info string to add a sentence
  analogous to the existing DCF one ("recutils is also supported (use
  --recutils for recutils in and out)"); do **not** add recutils to the
  10x10 keystroke-saver matrix.
- No changes needed to `pkg/cli/option_types.go` or
  `pkg/cli/flatten_unflatten.go` (see above).

### Docs

- `docs/src/file-formats.md.in`: new `## recutils` section, modeled on the
  DCF section, placed just before "## Data-conversion keystroke-savers".
  Include a `cat data/sample.rec` block and an
  `mlr -i recutils -o json cat data/sample.rec` block. Note the `+`
  continuation and the trailing-backslash round-trip limitation.
- `docs/src/data/sample.rec` (new): a small, illustrative example (e.g.
  address-book style, echoing GNU recutils' own canonical manual example),
  exercising at least one `+`-continuation.
- `docs/src/reference-main-separators.md.in`: new table row, shaped like the
  DCF row but with the PS cell reading `Always ": "; not alterable`
  (colon-*space*, not bare colon — malformed otherwise per the hard-error
  decision above) and RS reading `N/A; records separated by blank lines`.
- `docs/src/glossary.md.in` (~line 467): add "recutils" to the bracketed list
  of non-line-oriented formats (`CSV, TSV, JSON, YAML, DCF, recutils, and
  others`).
- `docs/src/reference-main-flag-list.md.in` is auto-generated from
  `option_parse.go` — no manual edit, just rebuild via
  `make -C docs/src forcebuild`.

### Tests

- `pkg/input/record_reader_rec_test.go` (new): unit tests for
  `recordFromRECLines` covering backslash-join, `+`-continuation (including
  bare `+`), colon-in-value-after-first-occurrence, `+` with no preceding
  field → error, missing `": "` → error, and comments interleaved with
  continuation lines (comments are stripped by the scanner before the
  continuation passes ever see them).
- `test/cases/io-recutils/` (new), modeled on `test/cases/io-dcf/`:
  - `0001`: `csv → recutils` conversion smoke test.
  - `0002`: `recutils → recutils` round trip.
  - `0003`: `recutils → json` conversion.
  - `0004`: `+`-continuation and comment handling exercised at the CLI level
    (the one genuinely new behavior vs. DCF, which has no continuation
    syntax).
- `test/input/test.rec` (new fixture): multi-record file including a
  `+`-continued field and a comment line, reused across the cases above (like
  `test/input/test.dcf` is reused across DCF's cases).

## Verification

1. `make build` then manual smoke tests:
   - `mlr --icsv --orecutils cat somefile.csv` — confirm blank-line-separated
     `Key: value` output.
   - `mlr --irecutils --ojson cat docs/src/data/sample.rec` — confirm nested
     JSON output round-trips field values correctly, including the
     `+`-continued one.
   - `mlr --irecutils --orecutils cat docs/src/data/sample.rec` — round trip
     should reproduce equivalent records (modulo the documented
     trailing-backslash limitation).
   - Feed a malformed file (line with no `": "`, or a leading `+` with no
     prior field) and confirm a clear error, not a crash or silent
     misparse.
2. `mlr regtest -p test/cases/io-recutils` to populate `expout`/`experr`,
   then review the diffs by eye before committing them.
3. `make dev` (fmt, build, unit + regression tests, docs rebuild) and
   `make lint` — both required before pushing per this repo's CLAUDE.md.
