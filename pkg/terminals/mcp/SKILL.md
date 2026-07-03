---
name: miller
description: >
  Drive Miller (mlr) to process CSV/TSV/JSON/etc. data. Use when constructing
  mlr command lines: discover capabilities from the catalog rather than
  guessing, learn the data's shape before writing expressions, validate DSL
  before running, and recover from failures via structured errors.
---

# Miller agent playbook

Miller (`mlr`) is a command-line data processor for CSV, TSV, JSON, JSON
Lines, and other tabular/record formats, with SQL-like verbs (`cut`, `sort`,
`join`, `stats1`, ...) and an awk-like DSL (`put`, `filter`).

Work this loop. Each step exists to prevent a specific, common failure.

## 1. Discover — never invent names

Everything valid is in the catalog; anything not in the catalog does not
exist. Hallucinated flag/function names are the top failure mode.

- Route an intent: `which` with e.g. `"join two files on a key"` → ranked
  candidates. `confident: true` means a name matched; trust the top hit.
- Browse cheaply: `list_capabilities` with `index: true` → every
  verb/function/flag/keyword with one-line summaries.
- Drill in: `list_capabilities` with `kind: "verb", names: ["join"]` → the
  full entry. Prefer the structured `options` list (flag, arg, type, enum
  `values`) when present; `usage_text` is the prose fallback.
- The whole catalog is cacheable against `(mlr_version,
  catalog_schema_version)` — re-fetch only when either changes.

## 2. Constrain — learn the data before touching it

Call `describe_data` on the input first. It returns, per field: name, types
seen with counts, occurrence count, null count, cardinality, min/max, and —
for low-cardinality fields — every distinct value.

- Copy field names exactly from `describe_data`; never guess casing or
  spelling.
- For flags like `-g` (group-by) and DSL comparisons, use values from the
  `values` array, not values you expect to exist.
- Fields whose `count` is less than other fields' are absent in some records:
  guard DSL with `is_present($field)`.

## 3. Validate — check DSL before spending a run

Before any `run` that includes `put` or `filter`, call `validate_dsl` with the
expression. Cost: parse-only, no data read. On `valid: false`, the `error`
document has `kind`, `hint`, and `did_you_mean` — apply the hint, don't
re-guess syntax.

## 4. Run — and read errors structurally

Call `run` with argv as a list, one element per shell word (no shell quoting):

    {"args": ["--icsv", "--ojson", "cat", "data.csv"]}

Command-line shape rules that prevent most argv errors:

- Main flags (I/O formats etc.) come **before** the verb: `mlr --icsv sort -f name f.csv`.
- Format shorthands: `--icsv --ojson` (separate in/out), `--csv`/`--c2j` etc. (combined).
- Chain verbs with `then`: `["--icsv", "sort", "-f", "k", "then", "head", "-n", "3", "f.csv"]`.
- If a field value being compared in `filter` might collide with a verb flag,
  end verb flags with `--` before filenames.
- Inline data goes in `stdin_text`; files go at the end of `args`.

On failure, `exit_code` is nonzero and `error` (when present) carries `kind`,
`hint`, and `did_you_mean` — `hint` is often a corrected command line; prefer
executing it over reasoning from the message. `stdout_truncated: true` means
the output exceeded the server's cap: narrow the query (e.g. `head`, `cut`)
rather than re-running the same command.

## Notes

- `run` cannot execute external commands (DSL `system`/`exec`, piped
  redirects, `--prepipe`) unless the server was started with `--allow-shell`;
  such calls fail cleanly. It **can** write files via `tee`, `split`, and DSL
  output redirects — treat it as a write-capable tool.
- Long inputs: prefer `describe_data` + targeted verbs over dumping whole
  files through `run`.
- One record format in, another out: Miller is format-to-format; there is no
  separate conversion step.
