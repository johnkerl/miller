# Plan: nested-field ("JSON") accessors for non-DSL verbs

Motivating issues:

- [issue #1815](https://github.com/johnkerl/miller/issues/1815) — `mlr -j rename
  Body.meta,Body.renamed_meta` silently matches nothing when `Body` is a nested map. The
  workarounds are DSL (`put '$Body.renamed_meta = $Body.meta; unset $Body.meta'`) or the
  flatten-sandwich (`flatten then rename Body.meta,Body.renamed_meta then unflatten`). Thread
  consensus: worth documenting; the deeper ask is that `rename` (and friends) understand paths.
- [issue #1534](https://github.com/johnkerl/miller/issues/1534) — nominally about CSV
  schema-change errors, but architecturally on point: `group-like` judged record schemas by
  *nested* key lists (`name,location_1,field_1`) while the CSV writer saw *flattened* keys
  (`name,location_1.1.lat,...`), so "same schema" groups came out heterogeneous. It's a live
  example of verbs seeing different field names than users see in non-JSON output — and of
  dotted flat names (`location_1.1.lat`) being real, user-visible identifiers in CSV-land.
- [Flatten/unflatten docs](https://miller.readthedocs.io/en/latest/flatten-unflatten/) — the
  existing contract this feature must not break.

The asymmetry, stated once:

- In the put/filter DSL, `$x.y.z` traverses nested structures: on `{"x":{"y":{"z":4}}}` it
  yields 4.
- In every other verb, a field name is a single flat string: `mlr cut -f x.y.z` looks for a
  field literally named `x.y.z` and finds nothing in that record.

Scope of this doc: plan and sharp edges only — no code. The recommendation (§ Proposed
design) is native path accessors at the Mlrmap level, adopted by a small verb set first,
behind an explicit per-verb opt-in.

## Current architecture (survey)

### Mlrmap accessors are flat; indexed accessors already exist

- Every verb-facing accessor takes a single `string` key and funnels through
  `findEntry(key)` (`pkg/mlrval/mlrmap_accessors.go:199`): `Has` :17, `Get` :21,
  `PutCopy` :95, `Remove` :505, `Rename` :724, `MoveToHead`/`MoveToTail` :514/:522, plus the
  bulk group-by helpers `GetSelectedValuesJoined` :583 / `GetSelectedValuesAndJoined` :610 /
  `HasSelectedKeys` :678, all looping flat lookups.
- Nesting-capable primitives already exist, keyed by `[]*Mlrval` index chains rather than
  dotted strings: `PutIndexed` (`mlrmap_accessors.go:538` → `putIndexedOnMap`,
  `pkg/mlrval/mlrval_collections.go:266`, with auto-deepening), `RemoveIndexed`
  (:542 → `removeIndexedOnMap`, `mlrval_collections.go:402` — removes only the leaf,
  leaving the parent map in place, possibly empty), and `getWithMlrvalArrayIndex`
  (`mlrmap_accessors.go:371`, the `$x[["a","b"]]` walker). These are battle-tested by the
  DSL and are the machinery a verb-side feature should reuse.

### How the DSL does `$x.y.z`

- The dot form is parsed as nested `DotOperator` nodes; `DotCallsiteNode.Evaluate`
  (`pkg/dsl/cst/builtin_functions.go:591-610`) does one-level `Get` per dot when the LHS is
  a map, else falls back to string concatenation. So `$x.y.z` is `(($x).y).z` — strictly
  leftmost, one key per segment, never "try `x.y` as a single key".
- The bracket form `$x["y"]["z"]` and lvalue assignments collect `[]*Mlrval` index chains
  and call `Get/Put/RemoveIndexed` (`pkg/dsl/cst/lvalues.go:115, 447, 532`).
- Crucially, the DSL already has a disambiguation *syntax*: `$x.y.z` means traversal,
  `${x.y.z}` means the literal flat name. The verbs have no equivalent syntax slot today —
  that's the heart of the design problem.

### How verbs consume field names

All flat, all inlined per verb; there is no shared field-name abstraction to hook:

- Parse: `cli.VerbGetStringArrayArg` turns `-f a,b,c` into `[]string`; ~40 verbs use it
  (cut, sort, having-fields, reorder, rename, stats1/2, top, uniq, count, fill-down,
  merge-fields, join, template, subs, case, ...).
- Lookup: cut tests `tr.fieldNameSet[pe.Key]` (`pkg/transformers/cut.go:197`) or
  `inrec.Get(name)` :221; having-fields iterates `inrec.Head` against a set; rename calls
  `inrec.Rename(pe.Key, ...)`; sort uses `GetSelectedValuesAndJoined`; reorder uses
  `MoveToHead`/`MoveToTail`.

### When verbs see nested vs flat records

Auto-flatten is appended *after* the whole verb chain, at write time
(`pkg/climain/mlrcli_parse.go:445-461`; decision logic and design rationale in
`pkg/cli/flatten_unflatten.go`):

- JSON→JSON: no flatten/unflatten inserted — verbs always see nested records. This is the
  case the feature targets.
- JSON→CSV: flatten runs after the last verb — verbs still see nested records. (#1534's
  group-like surprise lives here.)
- CSV→CSV / CSV→JSON: no unflatten on input — verbs see flat records, including literal
  dotted keys like `req.method`. The header comment at `flatten_unflatten.go:43-48`
  explicitly promises `mlr sort -f req.method` works on such data "with no surprises."
  That promise is the hard backward-compatibility constraint on this feature.

There is no per-verb "needs flattening" capability flag; the only verb-aware chain decision
is the `lastVerbName == "flatten"` check in `DecideFinalUnflatten`
(`flatten_unflatten.go:93`). `TransformerSetup` would be the natural home for such a flag if
one were needed.

### flatsep and prior art

- `TWriterOptions.FLATSEP`, default `"."` (`pkg/cli/option_types.go:94, :261`;
  `pkg/cli/separators.go:45`), set via `--flatsep`/`--jflatsep`. It's a writer-side option
  but is used by both the final flatten and the final unflatten, and as the default `-s`
  for the flatten/unflatten verbs.
- Unflatten's string→path splitter is `SplitAXHelper`
  (`pkg/mlrval/mlrmap_flatten_unflatten.go:279`) feeding `PutIndexed`; keys with empty
  segments (`.x`, `x..y`, trailing dot) are treated as literal, with a one-time stderr
  warning (:120-176). Flatten stringifies empty collections: `{}` → `"{}"`, `[]` → `"[]"`
  (`pkg/mlrval/mlrval_accessors.go:33-62`) — i.e. the flatten/unflatten round trip is
  *not* lossless.
- In-verb nesting precedent: `sort-within-records -r` recurses into submaps
  (`pkg/transformers/sort_within_records.go`); `flatten -f` flattens selected fields only.

## The core ambiguity

Given spec `x.y.z` and a record, the name can mean: a field literally named `x.y.z`; field
`x` holding map `{y: {z: ...}}`; field `x.y` holding `{z: ...}`; or field `x` holding
`{"y.z": ...}`. In general a spec with n dots has 2^n candidate splits — the user's
"triple-cased" is the n=2 case. Trying all splits is unpredictable and unimplementable in
any explainable way, so the plan is to never enumerate splits. Two deterministic rules:

1. **Exact flat key first.** If the record has a field literally named `x.y.z`, that's the
   match, full stop. This single rule preserves the CSV-world promise above: flat records
   with dotted headers behave exactly as today.
2. **Else strict per-segment traversal**, mirroring both the DSL dot operator and what
   flatten itself produces: split the spec on the separator; each segment is exactly one
   map key (or array index) at each level. `{"x.y": {"z": 4}}` is *not* reachable via
   `x.y.z` — accepted and documented; the escape hatches are the DSL and the
   flatten-sandwich. (Flatten produces the flat key `x.y.z` from that record, so the
   flatten-sandwich does reach it.)

This makes lookup two-cased, not exponential, and rule 1 means the flat interpretation
always wins when both exist in one record (sharp edge S3 below).

## Design options

- **Option A — flatten-sandwich sugar.** Automatically wrap the chain (or individual
  verbs) in `flatten ... unflatten`, mechanizing the known idiom. Pros: trivial to build;
  semantics are "the names you see in CSV output," which matches many users' mental model;
  regex verbs get path matching for free. Cons: the round trip is lossy (`{}`/`[]` become
  strings; the unflatten arrayify heuristic can turn maps with keys `"1","2",...` into
  arrays that weren't arrays; type inference re-runs on stringified values); whole-record
  flatten cost even when one field is touched; collisions when a record has both literal
  `x.y` and nested `x:{y:...}` (flatten produces duplicate keys); and it changes record
  shape for *other* verbs in the same chain unless scoped per-verb, which the chain
  architecture doesn't support today. Fine as a documented manual idiom; not recommended as
  the feature.
- **Option B — native path accessors** at the Mlrmap level, adopted verb by verb, reusing
  `Put/Remove/GetIndexed`. Pros: lossless, structure-preserving, precise per-verb
  semantics, no shape changes for neighboring verbs. Cons: real API surface; each verb
  needs its own semantic decisions (inventory below); long tail of verbs.
- **Option C — document-only.** What #1815 settled for. Zero risk, leaves the asymmetry.
  Worth doing regardless (the flatten-sandwich and DSL idioms belong in
  reference-verbs / flatten-unflatten docs), but it's not the feature.

**Recommendation: Option B, scoped to a small verb set, opt-in (Q1), with Option C's doc
work done in the same effort.**

## Opt-in surface

Never change default interpretation silently. Candidate surfaces:

- (a) Per-verb boolean flag, e.g. `mlr rename -p Body.meta,Body.renamed_meta` ("-p" for
  path; letter TBD per verb's free letters). Explicit, discoverable in each verb's help,
  adoptable verb-by-verb. Recommended.
- (b) Global main flag (`--nested-fields`) flipping interpretation for all supporting
  verbs. One switch, but action-at-a-distance, and a chain mixing verbs that do and don't
  support it becomes confusing.
- (c) In-name syntax, e.g. `-f '$.x.y.z'` (JSONPath-ish) or `x["y"]["z"]`. No flag needed,
  but invents a mini-language, collides in principle with literal names, and is unpleasant
  in shells.

Note that with the exact-key-first precedence rule, even a default-on behavior would be
almost backward compatible — the fallback only fires when the flat lookup misses, which
today yields "no match." But "almost" hides real changes: `cut -x -f x.y` and
`having-fields --none-defined x.y` would start *matching* where they matched nothing, and
per-record heterogeneity makes behavior data-dependent. Hence opt-in for v1; a default flip
can be revisited later (same posture as the `--iauto` plan).

## Proposed design

### 1. Path type and split

A parse-once type, e.g. in `pkg/mlrval` or `pkg/lib`:

```go
type FieldPath struct {
    original string     // the literal spec, for exact-match-first and for output naming
    indices  []*Mlrval  // split segments, ready for Get/Put/RemoveIndexed
}
```

Built once at verb-construction time from each `-f`/`-g` token (never per record). Splitting
uses the same separator and the same empty-segment rules as unflatten (`SplitAXHelper`
semantics): a spec with leading/trailing/doubled separators is treated as wholly literal —
no warning needed here since literal is always tried first anyway. Numeric segments become
int Mlrvals (1-based array indices, matching flatten output `x.1` and DSL indexing),
non-numeric become strings.

### 2. Mlrmap API additions

All additive, all delegating to existing indexed machinery, all honoring
exact-flat-key-first:

- `GetPath(path) *Mlrval` — `findEntry(original)` first, else walk `indices`.
- `HasPath(path) bool`.
- `RemovePath(path) bool` — leaf removal only; parent maps remain (matches DSL
  `unset $x.y`, per `removeIndexedOnMap`). Returns whether anything was removed.
- `PutPathCopy(path, value)` — via `PutIndexed` with auto-deepen (only needed by verbs
  that create fields; not needed for v1's cut/having-fields).
- `RenamePathLeaf(path, newLeafName) bool` — rename the last segment *within its parent
  map*, preserving position (`Mlrmap.Rename` semantics one level down). Cross-parent moves
  (`rename Body.meta,Other.meta`) are out of scope for v1 — that's a move, not a rename
  (S14).

Deliberately *not* added: any API that tries multiple splits of the original string.

### 3. Separator

Use `FLATSEP` (default `"."`), threaded from options to verb constructors, for consistency
with flatten/unflatten and with what users see in flattened output. Sharp edges S4 apply
(multi-char separators, writer-option provenance). Verbs that grow the opt-in flag could
also accept a per-verb `-s` override, mirroring the flatten/unflatten verbs.

### 4. v1 verb set and per-verb semantics

Driven by the issues and by expected demand; each needs its semantics pinned before code:

- **rename** (#1815): `rename -p old.path,new_leaf_name` — leaf rename in place. Decide:
  is the second element a full path (error unless it differs only in the leaf) or just the
  new leaf name? Recommend full path + validation, so the CLI shape matches non-p rename
  and the flatten-sandwich idiom (`rename Body.meta,Body.renamed_meta`).
- **cut**: `-f a.b` extracts preserving structure — output record `{"a": {"b": ...}}`, not
  `{"a.b": ...}` (Q6). Requires a "copy path into fresh record, creating parents" helper.
  `-x` removes the leaf, keeping siblings and the (possibly emptied) parent. `-o` ordering
  applies at top level of the reconstructed record. Interaction: two specs sharing a prefix
  (`-f a.b,a.c`) must merge into one `a`, preserving sub-order.
- **having-fields**: `--at-least`/`--all-defined` etc. gain path membership via `HasPath`.
  The regex variants stay flat (S10).
- **sort**: `-f a.b` sorts by the path value. Missing-path records need a defined ordering
  (today missing flat fields group at the end — reuse that). Path values that are maps or
  arrays: define as error-or-last (S12).
- **reorder**: plausible but semantically muddy (move leaf within its parent? hoist to top
  level?) — defer past v1 unless a crisp semantic emerges.

Follow-on tiers, each with its own sharp edges, explicitly out of v1:

- Group-by (`-g`) family: count, uniq, count-distinct, stats1, top, decimate,
  count-similar, fraction, histogram... Group-by keys go through joined-string map keys
  (`GetSelectedValuesJoined`); path lookups slot in, but map/array-valued results need a
  rule (S12), and the *output* field naming question (S7) hits every one of these.
- Value-field (`-f`) stats verbs: stats1/stats2/merge-fields/step — output names like
  `x.y_sum` are new *flat* names that will themselves auto-unflatten to `x: {y_sum: ...}`
  in CSV→JSON runs. Decide whether that's a feature or a bug before touching these.
- Leaf mutators: fill-down, fill-empty, sub/gsub/ssub, case, format-values — mechanical
  once the path API exists.
- join `-j/-l/-r` on nested keys — its own plan; join has a second reader and half/full
  streaming variants.

## Sharp edges inventory

- **S1 — split ambiguity.** 2^n candidate splits; resolved by exact-key-first + strict
  per-segment traversal, never split enumeration. `{"x.y": {"z": 4}}` unreachable via
  `x.y.z` — document with the flatten-sandwich as escape hatch.
- **S2 — literal dots are load-bearing in CSV-land.** `flatten_unflatten.go:43-48`
  explicitly promises flat `req.method` addressing; #1534 shows dotted flat names in real
  data. Exact-key-first preserves this even with the flag on; without the flag nothing
  changes at all.
- **S3 — both-present and per-record heterogeneity.** One record may carry literal `x.y`
  *and* nested `x:{y:...}`; different records in one stream may differ. Precedence makes
  each record deterministic, but users can still be surprised mid-stream — needs a docs
  callout, and `cut -x`/removal semantics must be verified against both-present records in
  regression tests (remove the flat one only? both? — recommend: flat only, since exact
  match won; test pins it).
- **S4 — separator provenance.** `FLATSEP` lives in *writer* options; a verb-side reader
  of it is a small layering smell (the `--iauto` plan hit the same issue from the other
  side). Multi-char separators (`--flatsep ::`) must work; a separator that also appears
  inside genuine nested keys (JSON keys may contain dots at any level) re-creates S1 one
  level down — strict segment matching means such keys are simply unreachable by path spec
  (reachable via DSL bracket syntax only). Document.
- **S5 — arrays.** Numeric segments as 1-based indices matches flatten output (`x.1`) and
  DSL aliasing (negative indices count from the end — decide whether to allow; recommend
  yes, free via `UnaliasArrayIndex`). Out-of-bounds reads → absent. Writes via auto-deepen
  create *maps* keyed `"1"`, not arrays (`NewMlrvalForAutoDeepen`) — v1 verbs don't
  auto-deepen, but any later Put-capable verb must decide (the unflatten `Arrayify` pass is
  what turns those into arrays today).
- **S6 — degenerate specs.** Leading/trailing/doubled separators: treat the whole spec as
  literal (unflatten precedent, `mlrmap_flatten_unflatten.go:120-176`); no path fallback.
- **S7 — derived-output naming.** Any verb that *creates* fields named after inputs
  (stats1 `x.y_sum`, merge-fields, step `x.y_delta`, count-distinct's `field` column)
  produces flat dotted names that downstream auto-unflatten will restructure. Deferring
  those verbs defers the problem, but the rule must exist before tier 2.
- **S8 — removal leftovers.** Removing the last child leaves `{}` (DSL-consistent). Under
  JSON output that's visible; under CSV output, flatten turns `{}` into the string `"{}"`.
  Consistent with `unset` today, but worth a regression case so it's chosen, not
  accidental.
- **S9 — performance.** Paths parse once at verb construction. Per-record cost when the
  flag is off: zero (existing code paths untouched). When on: exact `findEntry` first
  (hash-map hit for wide records), traversal only on miss. cut's hot loop currently does
  set-membership per record key — path mode inverts to per-spec probes; fine for typical
  spec counts, note in benchmarks (`make bench`).
- **S10 — regex forms.** `cut -r`, `having-fields --any-matching`, `rename -r` match flat
  key strings; a regex over nested structure is ill-defined. v1: regex + path flag is an
  error. A later option: match regexes against *flattened* names (S13's mental model), but
  that drags in flatten cost and S1 collisions — separate decision.
- **S11 — structure-preserving extraction (cut).** Building the output record requires
  copying partial subtrees with shared-prefix merging and stable ordering — new helper,
  needs its own unit tests (deep siblings, prefix overlap `a.b` + `a`, spec ordering vs
  `-o`).
- **S12 — non-scalar path results.** Group-by joining and sort comparison assume
  scalar-ish values. A path can resolve to a map/array. Options: error, json-encode for
  keying, or sort-last. Recommend json-encode for group-by keys (deterministic) and
  collections-sort-last for sort; decide before tier 2.
- **S13 — two mental models forever.** Users will hold both "flattened names" (CSV view)
  and "paths" (JSON view); this feature makes the second one real in verbs. The docs page
  (flatten-unflatten) must gain a section explaining that path specs and flattened names
  usually coincide (same separator, same segments) and exactly when they don't
  (S1's unreachable case, `{}`/`[]` lossiness, arrayify).
- **S14 — rename is not move.** `rename -p a.b,c.d` where the parent differs is a move
  with different ordering/overwrite semantics — reject in v1 with a clear error pointing
  at the DSL.
- **S15 — chain-position interactions.** The feature operates on whatever shape reaches
  the verb: on CSV input, paths mostly no-op (records are flat; exact-match rule handles
  it); after an explicit `flatten` verb, likewise. No new chain-insertion logic needed —
  and specifically, the existing `flatten then ... then unflatten` idiom must keep working
  unchanged (regression case).
- **S16 — REPL.** REPL verbs share transformer code and its own flatten/unflatten decision
  (`pkg/terminals/repl/verbs.go:600-615`); no divergence expected, but include a REPL
  smoke test.

## Phased implementation

Each phase independently mergeable, `make check` green throughout.

1. **Docs-first (Option C, immediate).** Document the DSL and flatten-sandwich idioms for
   nested rename/cut in reference-verbs and flatten-unflatten pages — closes the actual
   ask in #1815 regardless of the rest.
2. **Path core.** `FieldPath` + `GetPath`/`HasPath`/`RemovePath`/`RenamePathLeaf` in
   `pkg/mlrval`, unit tests covering S1/S3/S5/S6/S8 cases. Pure additive; no verb changes.
3. **v1 verbs.** rename, cut, having-fields, sort behind the per-verb opt-in flag; the
   structure-preserving extraction helper (S11); regression cases per verb including
   both-present records, arrays, CSV-input no-op, flatten-sandwich equivalence
   (`cut -p -f a.b` ≡ `flatten then cut -f a.b then unflatten` on lossless inputs).
4. **Tier 2.** Group-by family + S12 rule; then leaf mutators; stats value-fields last
   (S7 rule required first).
5. **Docs & help.** Verb help strings (feeds generated reference-verbs), flatten-unflatten
   page section (S13), man page via `make dev`.

## Open questions for the maintainer

- **Q1 — opt-in flag vs default-on.** Exact-key-first makes default-on *nearly* safe, but
  `cut -x`/having-fields matching semantics do change on nested data. Recommend per-verb
  opt-in flag now; revisit default at a major release.
- **Q2 — flag spelling.** One consistent letter across verbs (is `-p` free everywhere in
  the v1 set?) vs a long option `--paths` only. Long-option-only is safest against
  per-verb letter collisions.
- **Q3 — separator.** Reuse `FLATSEP` (recommended, one knob) vs hard `"."` vs per-verb
  `-s` only. If `FLATSEP`, note it's writer-scoped today (S4).
- **Q4 — rename second argument.** Full path (validated same-parent) vs bare new leaf
  name. Recommend full path for symmetry with flat rename and the sandwich idiom.
- **Q5 — negative array indices in specs.** DSL-consistent aliasing (recommend) vs
  positive-only (flatten never emits negatives, so specs-as-flattened-names don't need
  them).
- **Q6 — cut output shape.** Structure-preserving (recommended; JSON-native) vs
  flattened-key output (matches Option A's model). If anyone wants the latter they can
  say `flatten then cut`.
