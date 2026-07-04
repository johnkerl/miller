<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
<div>
<span class="quicklinks">
Quick links:
&nbsp;
<a class="quicklink" href="../reference-main-flag-list/index.html">Flags</a>
&nbsp;
<a class="quicklink" href="../reference-verbs/index.html">Verbs</a>
&nbsp;
<a class="quicklink" href="../reference-dsl-builtin-functions/index.html">Functions</a>
&nbsp;
<a class="quicklink" href="../glossary/index.html">Glossary</a>
&nbsp;
<a class="quicklink" href="../release-docs/index.html">Release docs</a>
</span>
</div>
# The Miller Agent Skill

As of Miller version 6.20, released in July 2026, there are two main ways to get your AI to know
about a software tool (Miller, or others): **agent skills**, and [**MCP**](mcp-server.md).  (See
[Miller and AI](ai.md) for an introduction.)

Miller ships a built-in [Agent Skill](https://www.anthropic.com/news/skills) -- a single `SKILL.md`
file -- inside the `mlr` executable, so agents that read skills directly from disk (Claude Code,
and other tools that support the Agent Skills format) can discover and drive Miller without
scraping help text or guessing at flags.

The skill is plain markdown with a YAML frontmatter header, placed where your agent already looks
for skills. The agent reads it into context once, the same way it reads any other instructions, and
from then on it runs `mlr` commands via whatever shell-executing tool it already has.

Here's what the skill file looks like:

<pre class="pre-highlight-in-pair">
<b>mlr skill print | head -n 15</b>
</pre>
<pre class="pre-non-highlight-in-pair">
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

</pre>

For more background on the `mlr` commands the agent runs on your behalf, please see
[Miller AI internals](ai-support.md).

## Setup

Write the skill file to Claude Code's personal skills directory (do this before starting your
`claude` session):

<pre class="pre-highlight-in-pair">
<b>mlr skill install ~/.claude/skills/miller</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Wrote /Users/kerl/.claude/skills/miller/SKILL.md
</pre>

For Codex and Gemini:

<pre class="pre-highlight-non-pair">
<b>mlr skill install ~/.agents/skills/miller</b>
</pre>

With no argument, `install` writes to `.claude/skills/miller/SKILL.md` under the current directory
instead. This is handy for a project-scoped skill checked into that project's repo rather than one
installed for every project on your machine:

<pre class="pre-highlight-in-pair">
<b>mlr skill install</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Wrote .claude/skills/miller/SKILL.md
</pre>

There's no "uninstall" subcommand, since `install` only ever writes one plain file. Removing it is
an ordinary file operation:

<pre class="pre-highlight-non-pair">
<b>rm -rf ~/.claude/skills/miller</b>
</pre>

Then -- just interact with your agent as always! When you say something like `describe the data file example.csv`,
the agent will already know how to use Miller to help answer that question.

## What the Miller skill maps to

You don't have to type `skill` or anything else special in your agent session: rather you've
empowered the agent to discover things about Miller for itself. But if you're curious what's
actually placed in front of it:

<pre class="pre-highlight-in-pair">
<b>mlr skill --help</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Usage: mlr skill {print|install} [options]
Puts the Miller Agent Skill (SKILL.md) where a coding agent can find it.
This is the same playbook mlr mcp serves as its "miller-playbook"
prompt/resource, packaged for agents that read Agent Skills from disk.

Subcommands:
  print          Write the skill content to stdout.
  install [DIR]  Write DIR/SKILL.md, creating DIR if needed.
                 Default DIR is .claude/skills/miller

 -h or --help   Show this message.
</pre>

And here's the file itself -- the whole thing, not an excerpt, since this and nothing else is what
the agent has to go on:

<pre class="pre-highlight-in-pair">
<b>mlr skill print</b>
</pre>
<pre class="pre-non-highlight-in-pair">
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
</pre>

That playbook is prose, not named tools, but it rests on the Miller features documented
in the [Miller AI internals](ai-support.md) page.

## What using the Miller skill looks like in practice

There's no server status to check and no tool list to browse -- the skill is just text the agent
already has -- so "in practice" mostly looks like an ordinary conversation. Say you're looking at
[example.csv](example.csv) for the first time:

> **You:** In example.csv, show me the red rows.

Without the skill, a plausible guess for the DSL is `$color == "Red"` -- and Miller silently
returns nothing for it, since the real values are lowercase. With the skill installed, the agent
runs `mlr --icsv --ojson describe example.csv` on your behalf first, sees the real value set for
`color` (`yellow`, `red`, `purple`), and only then answers:

> **Agent:** Four rows have color = red: rows 2, 3, 4, and 6.

The full worked version of this example, including the exact commands run at each step, is in
[Miller and AI](ai.md#before-and-after-a-first-session-with-the-skill-installed).

## A note on sandboxing

The [MCP server](mcp-server.md) enforces a sandbox by construction: subprocesses it spawns run with
`MLR_NO_SHELL=1` unless you start it with `--allow-shell`, so an agent-constructed command line
can't execute external commands even if the agent wanted it to.

The skill file has no equivalent enforcement. It's advisory text, not a wrapper around subprocess
execution -- nothing stops an agent from running `mlr put 'end{print system("whatever")}'` with
your full shell permissions if it decides to. If you want that guarantee with the skill alone,
set the `MLR_NO_SHELL` [environment variable](reference-main-env-vars.md) yourself (or pass
`--no-shell` explicitly), rather than relying on the playbook text for isolation. If you want the
enforced version, register the [MCP server](mcp-server.md) instead of, or alongside, the skill.
