# Roadmap: making Miller more AI-friendly (issue #2098)

This is a living roadmap for making Miller drivable by an LLM agent, derived
from [issue #2098](https://github.com/johnkerl/miller/issues/2098) and
@aborruso's comment on it. Each PR section below is self-contained so that a future
PR can be opened against it. Update status as work lands.

## Context

Miller already has near-complete *introspection coverage* (`mlr help topics`:
verbs, functions, keywords, flags, exact/approximate search). The gap for agents
is **shape, not coverage**: nearly everything is emitted as human prose via
`fmt.Printf`, so an agent must scrape text and ends up hallucinating flags and
signatures — the highest-volume failure mode. The arc below moves Miller's
introspection surface from prose to a stable, parseable structure, then builds
operability (self-correction, validation, an MCP server) on top of it.

Two tracks, per the issue:
- **Discoverability** — how an agent learns what Miller can do (structured help
  catalog, capability index/router, worked-example corpus).
- **Operability** — how an agent runs Miller and self-corrects (structured
  errors, a DSL validate/dry-run, a `describe` schema verb, an MCP server).

### Grounding facts (verified in the codebase)

- **Help dispatch is name-based string matching with no flag parsing.**
  `HelpMain(args []string)` (`pkg/terminals/help/entry.go:232`) strips `help`,
  special-cases `find`, then matches `args[0]` against `handlerLookupTable`
  (`entry.go:254-276`); unmatched falls through to exact/approximate search
  (`entry.go:279-281`). Handlers are `zaryHandlerFunc`/`varArgHandlerFunc`
  (`entry.go:43-55`). An `--as-json` modifier must therefore be **extracted
  from args before dispatch**, not parsed by an existing flag layer.
- **All catalog structs use private (lowercase) fields and have no JSON tags:**
  - `BuiltinFunctionInfo` — `pkg/dsl/cst/builtin_function_manager.go:42` (name,
    class, help, examples, arity fields). Registry:
    `BuiltinFunctionManagerInstance`; accessors `LookUp`,
    `GetBuiltinFunctionNames`, `ListBuiltinFunctionsInClass`.
  - `Flag` / `FlagSection` / `FlagTable` — `pkg/cli/flag_types.go:66,78,86`
    (name, altNames, **arg** in curly-brace notation `{a,b,c}`, help, parser,
    suppressFlagEnumeration). Accessors `GetFlagNames`, `ListFlagsForSection`,
    `FlagTakesArg`.
  - `TransformerSetup` — `pkg/transformers/aaa_record_transformer.go:52` (Verb,
    UsageFunc, ParseCLIFunc, IgnoresInput). Registry `TRANSFORMER_LOOKUP_TABLE`
    (`aaa_transformer_table.go`); accessors `LookUp`, `GetVerbNames`,
    `ShowHelpForTransformer`.
  - Keywords — `KEYWORD_USAGE_TABLE` of `{name, usageFunc}`
    (`pkg/dsl/cst/keyword_usage.go:11-74`); help lives *inside* the func bodies.
  - **Consequence:** serialization needs exported DTO/"view" structs populated
    from these registries — we cannot just add JSON tags to private fields, and
    we should not export the internals (this keeps the wire shape decoupled and
    versionable).
- **Verb usage and keyword help write directly to the terminal.** Each verb
  hand-writes `UsageFunc(*os.File)` that `Printf`s its options (e.g.
  `pkg/transformers/cat.go:22`); keyword `usageFunc()` prints to stdout
  (`pkg/dsl/cst/keyword_usage.go`). **We refactor these sinks rather than
  hijacking the file descriptor:** change `TransformerUsageFunc` and the keyword
  usage funcs to take an `io.Writer`, with existing callers passing `os.Stdout`.
  A buffer then collects the same text cleanly, with no pipe/redirect tricks.
  Verb options remaining prose-only is the Tier-1/Tier-2 dividing line.
- **`FLAG_TABLE.NilCheck()`** (`pkg/cli/flag_types.go:310`) is the existing
  build-time completeness pattern (exercised via a `mlr help` entrypoint + a
  regression test). We mirror it to track verb-option migration in PR3.

---

## Cross-cutting design (applies to all PRs)

1. **DTO layer.** Add exported view structs in a new package (proposed
   `pkg/terminals/help/catalog/` or `pkg/help/catalog/`) — e.g. `Catalog`,
   `FunctionInfo`, `FlagInfo`, `VerbInfo`, `KeywordInfo`, `OptionSpec` — each
   with explicit `json:"..."` tags (snake_case). Populate them from the existing
   registries via the accessors above. Internal structs stay private; the DTO is
   the stable wire contract.
2. **Versioning / cache keys.** Every full/partial JSON document carries
   top-level `mlr_version` (from the same source as `mlr version`) and
   `catalog_schema_version` (an integer bumped on shape changes). Miller is a
   static binary, so the catalog changes only when the binary does — these make
   the dump a perfect cache key for an MCP server or any tool (re-fetch only on a
   binary/schema bump; no TTLs).
3. **Opt-in.** Two equivalent ways to ask for JSON, neither spelled `--json`
   (that top-level flag already means JSON I/O format):
   - **Per-call flag `--as-json`** — used inside the `help` namespace, where it
     is unambiguous (e.g. `mlr help --as-json`, `mlr help verb cat --as-json`).
   - **Env var `MLR_HELP_JSON` (truthy)** — a set-once global so an agent opts
     in once rather than per-call.
   `--as-json` and a truthy `MLR_HELP_JSON` are equivalent; the flag wins if
   both are present. Centralize the "should I emit JSON?" decision in one helper.
4. **Output discipline.** JSON goes to stdout, one document per invocation, no
   colorization, deterministic key/element ordering (sort by name) so diffs and
   agent parsing are stable.
5. **Examples never rot.** Worked examples surfaced in the catalog are
   CI-tested; aim for a runnable example on **every** verb (not just functions)
   — an agent pattern-matches off `mlr cat -n -g shape` faster than off prose.
   Hook into the existing regression-test / docs-build machinery.

---

## PR 1 — Tier 1: `mlr help --as-json` machine-readable catalog (foundation)

**Goal.** One call yields a structured, parseable model of Miller's entire
surface; per-item `--as-json` for targeted fetches. Plain (no-`--as-json`)
output is byte-for-byte unchanged. Everything downstream builds on this.

**Surface.**
- `mlr help --as-json` — full catalog as one JSON document.
- `mlr help verb cat --as-json` — one or more verbs.
- `mlr help function splitax --as-json` — one or more functions.
- `mlr help flag --ifs --as-json` — one or more flags.
- `mlr help keyword ENV --as-json` — one or more keywords.
- A truthy `MLR_HELP_JSON` makes all of the above emit JSON without the flag.

**Shape (Tier 1).**
- `mlr_version`, `catalog_schema_version` at top level.
- **Functions:** `name`, `class`, `help`, `examples[]`, arity info, **and a
  structured signature** `{params: [{name, type}], return: type}` — see the
  signature note below.
- **Flags:** `section`, `name`, `alt_names[]`, `arg`, `help`.
- **Verbs:** `name`, `summary` (one line), `ignores_input`, and `usage_text`
  (the verb's rendered `UsageFunc` output) as the Tier-1 fallback for
  not-yet-structured options.
- **Keywords:** `name`, `help` text.

**Implementation.**
- New DTO package (cross-cutting #1).
- **Render verb usage via an `io.Writer`, not a captured fd.** Change
  `TransformerUsageFunc` (and the keyword usage funcs) to take `io.Writer`;
  existing callers pass `os.Stdout`, and the catalog builder passes a
  `bytes.Buffer` to collect `usage_text` / keyword help. This is the "right
  place" refactor — no pipe/`os.File` hijacking. Touch the
  `TransformerUsageFunc` typedef (`aaa_record_transformer.go`), the dispatch in
  `aaa_transformer_table.go:85`, every verb's `UsageFunc`, and the keyword
  usage funcs (`keyword_usage.go`). Mechanical but broad.
- **Structured function signatures (go deeper, don't parse prose).** Rather than
  scraping the human first-line, derive `{params, return}` from the
  function-info API in `builtin_function_manager.go`: the arity fields
  (`hasMultipleArities`, `minimum/maximumVariadicArity`) plus the typed func
  pointers (`unaryFunc`, `binaryFunc`, …) already encode arity/shape. Add
  accessor(s) on `BuiltinFunctionInfo` that expose this as structured data and
  feed the DTO. Keep the human first-line in `help` too.
- **`--as-json` extraction:** in `HelpMain` (`entry.go:232`), scan/strip
  `--as-json` (and consult `MLR_HELP_JSON`) before the name-based dispatch
  (`entry.go:254`); thread a `wantJSON bool` into the per-topic handlers. Add a
  builder that walks all four registries for the no-arg full-dump case.
- Reuse the registry accessors listed in Grounding facts; no registry refactor.

**Tests.** Golden-JSON regression cases under the existing regression harness; a
schema-completeness test (every function/flag/verb/keyword appears; required
fields non-empty) in the spirit of `NilCheck`.

---

## PR 2 — Discovery: JSON index + capability router

**Goal.** Cheap first calls so an agent can *choose* before drilling in.

- **`mlr help --as-json --index`** → `[{kind, name, summary}]` across verbs,
  functions, flags, keywords — names + one-line summaries only, no
  bodies/examples/usage_text. (Delta over existing `list-verbs`/`list-functions`,
  which are names-only.) Reuse the summary extraction from PR1. This is the cheap
  first call that lets an agent pick a verb before fetching its full entry.
- **`mlr which "join two files on a key"`** → ranked JSON
  `[{verb, score, summary}]`. Build on Miller's existing exact/approximate help
  search (`helpByApproximateSearchOne` and the `*Approximate*` accessors in
  `entry.go`). **Signal confidence via exit code** (e.g. `0` confident match,
  `2` no confident match) so the agent branches on status, not prose. `mlr which`
  is the reverse of `--index` (intent → verb vs. browse-all), short-circuiting
  the common "which verb?" round-trip.

**Tests.** Index covers every catalog item; `which` returns the expected top
verb + exit code for a handful of canonical intents.

---

## PR 3 — Tier 2: structured verb options (+ enum value-sets)

**Goal.** Replace each verb's `usage_text` blob with a structured option list;
verbs upgrade independently.

**Model.**
- Add optional `Options []OptionSpec` to `TransformerSetup`
  (`aaa_record_transformer.go:52`), default `nil`.
- `OptionSpec`: `{Flag, Arg, Type, Desc string; Repeatable bool; Values []string}`.
  `Type` is a small enum: `bool | string | int | float | csv-list | regex |
  filename | format | enum`.
- **Finite domains emit their value set:** where an option has a fixed domain
  (e.g. output format), set `Type:"enum"` and populate `Values`
  (e.g. `["csv","tsv","json","jsonl","pprint","markdown","dkvp","nidx","xtab"]`).
  Agents hallucinate *values*, not just flags — emitting the actual enum attacks
  value-hallucination at the source.
- **Scope: static domains only.** `Values` here is @aborruso's *codelist* — the
  set fixed by the binary (output formats, compression types). His *constraint*
  case — values that are only valid given the current input (e.g. a field name
  for `-g`) — is data-dependent and out of scope for the static catalog; that
  belongs to `mlr describe` (PR6), which reads the input schema. Keep the line
  clean: PR3 enums are binary-fixed, never data-derived.

**Emitter.** Prefer `Options` when non-nil; otherwise fall back to `usage_text`.
Agents always get *something*; no big-bang migration. Optionally render each
verb's `UsageFunc` *from* `Options` so prose and JSON stay in sync.
(Done post-migration: `WriteVerbOptions` in `aaa_verb_usage.go` renders each
usage message's "Options:" block from the specs; all 70 verbs migrated.)

**Migration tracking.** Add a `VerbOptionsNilCheck` mirroring
`FLAG_TABLE.NilCheck()` (`flag_types.go:310`) wired through a `mlr help`
entrypoint (`entry.go`) and asserted in a regression test: report which verbs
still have `Options == nil`. Migrate verbs incrementally here and in follow-ups.

---

## PR 4 — Structured errors: `--errors-json`

**Goal.** Agents branch on error *kind* instead of regex-matching English; the
catalog becomes the dictionary errors resolve against. (Biggest operability win
per the issue.)

- `--errors-json` emits `{error, kind, verb, position, hint, did_you_mean[]}`.
- **`did_you_mean`:** Levenshtein nearest-match over verb/flag/function/keyword
  names from the PR1 catalog — closing the self-correction loop the catalog
  enables.
- **`hint` and `did_you_mean` are copy-pasteable corrected command lines**, not
  prose (e.g. `mlr cut -f x,y -- file.csv`) — agents recover from a command far
  faster than from a sentence describing the fix.
- Identify Miller's central CLI/DSL error-emission points and route them through
  a structured-error type when the flag (or the `MLR_HELP_JSON`-style global) is
  set.

---

## PR 5 — DSL `--explain` / validate dry-run  *(landed)*

**Goal.** Validate/type-check a DSL expression *before* spending a full input
pass (a big context saver for agents).

- `mlr put --explain '...'` (and `mlr filter --explain`) parse + type-check the
  DSL, report errors (ideally via the PR4 structured-error path), and exit
  **without consuming the full input stream**.
- Reuse the existing DSL parse/CST build path; gate it before the record loop.

**Landed.** `--explain` added to put/filter (`put_or_filter.go`): after the
existing `cstRootNode.Build` (which already does parse → ValidateAST → CST build
→ Resolve), a valid expression prints `mlr {put,filter}: DSL expression is
valid.` and exits 0; an invalid one returns the build error up the normal path,
so `--errors-json` yields a structured document. The gate sits in the pass-two
constructor, before any input file is opened, so no input is read. DSL parser
messages (`parse error: ...`) now categorize as `dsl-parse-error` rather than
`generic` (`climain/errors_json.go`). Tests: `dsl-explain/0001-0004` regression
cases (valid put/filter, invalid plain, invalid `--errors-json`) plus categorize
unit tests. Note: the older `-X` ("exit after parsing") still exits 0 even on a
parse error — a pre-existing quirk left as-is since `--explain` is the correct
validation path.

---

## PR 6 — `mlr describe` schema/shape introspection  *(landed)*

**Goal.** Let an agent learn the *data's* shape, complementing the catalog's
*tool* shape.

- `mlr describe` (verb or terminal) reports field names, inferred types, and
  cardinality over the input stream, with an `--as-json` form.
- Leverage Miller's existing type-inference (`pkg/mlrval`) and field-collection
  machinery; likely a new verb in `pkg/transformers/`.

**Landed.** New verb `describe` (`pkg/transformers/describe.go`), registered in
`TRANSFORMER_LOOKUP_TABLE` with Tier-2 `Options` so it appears structured in
the PR1 catalog and PR2 index automatically. One output record per input
field: `field_name`, `types` (type-name → occurrence-count map, via
`GetTypeName` type inference), `count`, `null_count`, `distinct_count`,
`min`/`max`, and — for fields whose cardinality is within `-n`/`--max-values`
(default 20; 0 suppresses) — a `values` array listing every distinct value in
first-seen order. The `values` list is the data-derived *constraint* domain
deferred out of PR3: an agent copies real values for `-g`, DSL comparisons,
etc. instead of guessing. The JSON form is Miller-native — `mlr --ojson
describe` — with `types`/`values` nesting in JSON and auto-flattening in
tabular formats, so no verb-level `--as-json` flag was needed; `describe` is
positioned relative to `summary` as schema-shape vs. summary-statistics.
Distinctness is on string representations, matching `summary`'s
`distinct_count`; null semantics (empty or JSON null) match `summary`'s
`null_count`. Tests: `test/cases/verb-describe/` (JSON, pprint-flattened,
heterogeneous input, `-n` cap, `-n 0`, null-vs-empty, bad-option); docs:
`## describe` in `reference-verbs.md.in`.

---

## PR 7 — MCP server + Agent Skill (the loop)

**Goal.** Package the above so an agent gets both the *surface* and the *loop*.

- Thin MCP tool-server wrapping the catalog and friends: `list_capabilities`
  (PR1/PR2), `validate_dsl` (PR5), `describe_data` (PR6), `run`.
  `list_capabilities` caches the dump keyed on `mlr_version` (PR1).
- **Ship an Agent Skill / playbook** encoding the discover → constrain →
  validate → run loop — the recipe is what makes a CLI "shine when driven by an
  agent," beyond the raw tool surface.

**Design (worked out; ready to build against).**

- **Transport: stdio, not HTTP.** MCP's standard transport for local tools is
  JSON-RPC 2.0 over stdin/stdout: the client (Claude Code, Claude Desktop,
  Cursor, …) spawns the server as a subprocess. No localhost port, no auth
  story, no firewall prompts, works offline, dies with the client. MCP's
  streamable-HTTP transport exists for *remote* servers only; if ever wanted it
  could be a later opt-in flag, but nothing here needs it.
- **Entry point: a new terminal in the existing binary — `mlr mcp`** — not a
  separate executable. Registration is `claude mcp add miller -- mlr mcp`
  (or the equivalent JSON config). Shipping inside `mlr` means zero extra
  install and the server is version-locked to the binary, which is exactly
  what the `mlr_version` + `catalog_schema_version` cache keying assumes.
- **Dependency decision (the main open call before starting):** official
  `modelcontextprotocol/go-sdk`, vs. hand-rolling the small protocol subset
  needed (`initialize`, `tools/list`, `tools/call`, `ping`) over
  `encoding/json` — a few hundred lines, stable wire format. SDK leans toward
  spec conformance as MCP evolves; hand-rolling fits Miller's
  near-stdlib-only ethos.
- **Tools** — thin wrappers over PR1–PR6, nothing new underneath:
  - `list_capabilities` — PR1 catalog / PR2 index, with kind/name filters so
    an agent fetches one verb entry cheaply.
  - `which` — intent → ranked verbs (PR2).
  - `validate_dsl` — the PR5 `--explain` path; failures return the PR4
    structured-error document.
  - `describe_data` — the PR6 verb with `--ojson`.
  - `run` — execute an mlr command line; returns stdout (size-capped, with a
    truncation note), stderr, exit code, and the parsed `--errors-json`
    document when one fired.
- **Execution model: in-process for pure lookups, subprocess for execution.**
  Catalog and `which` are pure functions over compiled-in registries — serve
  in-process. `validate_dsl`, `describe_data`, and `run` shell out to the same
  binary via `os.Executable()`: the CLI paths call `os.Exit`, mutate global
  option state, and can panic, so subprocess isolation is simpler and
  guarantees the agent sees byte-identical behavior to a terminal.
- **Safety wrinkle — `run` is arbitrary-code-execution by design.** The DSL
  has `system()` and `exec()`, and verbs like `tee`/`split` write files. MCP
  clients prompt per tool call, but the server should still enforce: a
  timeout, an output cap, and — as a small prerequisite piece of this PR — a
  new `--no-shell`-style flag in Miller itself that the server sets by default
  (with explicit opt-out), so `system`/`exec` fail cleanly unless the user
  asks. The `run` tool carries the MCP `destructiveHint` annotation.
- **Skill half, single-sourced twice.** A playbook encoding *describe →
  index/which → fetch entry → validate → run*, branching on structured error
  `kind` and `did_you_mean`. Ship the same content as an in-repo Agent Skill
  (SKILL.md) and as an MCP prompt/resource exposed by the server itself, so
  agents that only see the server still get the loop.
- **Tests.** Golden-transcript tests that spawn `mlr mcp` and drive it over
  stdio (initialize → tools/list → each tool), plus unit tests on the tool
  handlers. The `MLR_AGENT` open question below lands here: the server makes
  it mostly moot (it sets flags explicitly), but it's worth resolving for the
  skill-without-server case.

---

## Open questions (carry into the relevant PR; not blocking the roadmap)

- **Env-var scope:** `MLR_HELP_JSON` flips help/catalog output. Should the same
  (or a broader `MLR_AGENT`) env var also flip `--errors-json` on, so an agent
  sets one variable for both? (Decide when PR4 lands.)
- **`mlr help schema` alias** for the full dump, in addition to the `--as-json`
  flag? (Distinct from publishing a JSON Schema *describing* the catalog
  document, which the exported Go DTOs already serve as a de-facto version of.)

Resolved: the per-call flag is `--as-json` (with `MLR_HELP_JSON` as the env-var
equivalent); function signatures are emitted structurally from the function-info
API (PR1), not parsed from prose.
