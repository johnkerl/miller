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
# asdfsadssf

## Using a skill file

<pre class="pre-highlight-in-pair">
<b>$ mlr skill print | head -n 15</b>
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

<pre class="pre-highlight-in-pair">
<b>$ mlr skill install /tmp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Wrote /tmp/SKILL.md
</pre>

<pre class="pre-highlight-in-pair">
<b>$ mlr skill install</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Wrote .claude/skills/miller/SKILL.md
</pre>

<pre class="pre-highlight-in-pair">
<b>$ mlr skill install ~/.claude/skills/miller</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Wrote /Users/kerl/.claude/skills/miller/SKILL.md
</pre>

## Why AI support

Miller treats AI agents as first-class users. When an agent drives a
command-line tool, it can fail in predictable ways: it invents flags that don't
exist, guesses values that aren't in the data, misreads error prose, and
burns whole runs discovering a typo. Miller closes off each of those failure
modes with the following structure:

* Miller's entire surface -- verbs, DSL functions, flags, keywords -- is
  available as **machine-readable JSON**, so agents ground themselves in what
  actually exists.
* Options with fixed domains carry their **complete value sets**, and input
  data can be **profiled in one pass**, so that agents copy real values instead
  of inventing them.
* DSL expressions can be **validated before running**, without reading any
  input.
* **Errors are structured** -- kind, hint, did-you-mean -- so agents branch
  on data rather than parsing English.
* A **sandbox flag** removes external-command execution, so an
  agent-constructed command line is just data processing.

Everything on this page is an ordinary command-line feature: it works from
any agent harness, system prompt, or script -- and it's equally useful for
plain shell tooling like `jq`.

## The essentials

**To get the AI features:** install Miller 6.20 or newer ([Installing
Miller](installing-miller.md)). That's it. Everything on this page ships inside the ordinary `mlr`
binary -- there are no plugins, no separate installs, no API keys, and nothing here makes network
calls.

To get your AI to see these features, pick whichever matches your setup.

### If your agent speaks MCP

For Claude Code, Claude Desktop, Cursor, etc.: register the "Miller MCP server", which is simply
having the AI run the `mlr` executable to ask it questions. For Claude Code, that's

<pre class="pre-highlight-in-pair">
<b>claude mcp add miller -- mlr mcp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Added stdio MCP server miller with command: mlr mcp to local config
File modified: /Users/kerl/.claude.json [project: /Users/kerl/git/johnkerl/miller]
</pre>

The MCP tools describe themselves, and the `mlr` binary ships its own instructions and playbook, so
you usually don't need to say anything special; if the agent doesn't reach for them, a nudge like
"use the Miller tools" suffices. Details in [The MCP server](mcp-server.md).

What happens to your system when you run this? Only that Claude Code will remember to run the `mlr`
binary -- the same one you use a the command line -- with command-line options that help Claude talk
to it.  No webserver is installed.

To uninstall, you can do

<pre class="pre-highlight-in-pair">
<b>claude mcp remove miller</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Removed MCP server "miller" from local config
File modified: /Users/kerl/.claude.json [project: /Users/kerl/git/johnkerl/miller]
</pre>

What happens to your system when you run this? It tells Claude Code to forget about running the `mlr`
binary to get how-to instructions.

### If your agent just runs shell commands

A system prompt, a `CLAUDE.md`, Cursor rules, a script harness): paste this standing instruction
into its context:

<pre class="pre-non-highlight-non-pair">
Miller (mlr) is installed for processing CSV/TSV/JSON/etc. data. When
constructing mlr commands:
1. Discover: `mlr help --as-json --index` lists every verb/function/flag;
   `mlr which "&lt;intent&gt;"` routes a goal to the right one; `mlr help
   verb &lt;name&gt; --as-json` gives full details. Never invent flag or
   function names.
2. Constrain: `mlr --icsv --ojson describe &lt;file&gt;` (or --ijson etc.)
   shows the data's fields, types, and values. Copy names and values from it
   rather than guessing them.
3. Validate: check DSL expressions with `mlr put --explain '&lt;expr&gt;'`
   before using them.
4. Run with `--errors-json`; on failure, correct using the error's kind,
   hint, and did_you_mean rather than re-guessing.
</pre>

  A fuller, ready-made version of that lesson ships in the Miller source
  tree at
  [pkg/terminals/mcp/SKILL.md](https://github.com/johnkerl/miller/blob/main/pkg/terminals/mcp/SKILL.md),
  in Agent Skill format -- suitable for dropping into e.g. a
  `.claude/skills/miller/` directory as-is.

The rest of this page is what those instructions rest on, feature by
feature.

## Discover: the machine-readable catalog

`mlr help --as-json` emits Miller's entire help catalog as one JSON document.
The `--index` form is the cheap first call -- every capability with a
one-line summary (here trimmed, and then counted, using Miller itself):

<pre class="pre-highlight-in-pair">
<b>mlr help --as-json --index | mlr --json head -n 2</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
{
  "kind": "verb",
  "name": "altkv",
  "summary": "Given fields with values of the form a,b,c,d,e,f emits a=b,c=d,e=f pairs."
},
{
  "kind": "verb",
  "name": "bar",
  "summary": "Replaces a numeric field with a number of asterisks, allowing for cheesy"
}
]
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr help --as-json --index | mlr --json count</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
{
  "count": 661
}
]
</pre>

From the index, an agent drills into full entries one at a time: `mlr help
verb sort --as-json`, `mlr help function splitax --as-json`, `mlr help flag
--ifs --as-json`, `mlr help keyword ENV --as-json` -- each accepting one or
more names. A verb entry carries a structured option list -- flag, argument
placeholder, type -- alongside the familiar usage text:

<pre class="pre-highlight-in-pair">
<b>mlr help verb decimate --as-json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
  {
    "name": "decimate",
    "summary": "Passes through one of every n records, optionally by category.",
    "ignores_input": false,
    "options": [
      {
        "flag": "-b",
        "type": "bool",
        "desc": "Decimate by printing first of every n."
      },
      {
        "flag": "-e",
        "type": "bool",
        "desc": "Decimate by printing last of every n (default)."
      },
      {
        "flag": "-g",
        "arg": "{a,b,c}",
        "type": "csv-list",
        "desc": "Optional group-by-field names for decimate counts, e.g. a,b,c."
      },
      {
        "flag": "-n",
        "arg": "{n}",
        "type": "int",
        "desc": "Decimation factor (default 10)."
      }
    ],
    "usage_text": "Usage: mlr decimate [options]\nPasses through one of every n records, optionally by category.\nOptions:\n-b         Decimate by printing first of every n.\n-e         Decimate by printing last of every n (default).\n-g {a,b,c} Optional group-by-field names for decimate counts, e.g. a,b,c.\n-n {n}     Decimation factor (default 10).\n-h|--help  Show this message."
  }
]
</pre>

Note that `usage_text` -- what `mlr decimate --help` prints -- is rendered
*from* the same structured options, so the human help and the machine help
cannot drift apart. Function entries carry name, class, arity, help, and
examples; the examples across the whole catalog are exercised by Miller's
test suite, so they never rot.

Three properties make the catalog cheap to use:

* **It's a perfect cache key.** Every document carries `mlr_version` and
  `catalog_schema_version`. Miller is a static binary, so the catalog changes
  only when the binary does: fetch once, cache forever, re-fetch on a version
  bump. No TTLs.
* **It's deterministic.** One document per invocation, sorted entries, no
  colorization -- stable for diffing and for prompt caches.
* **It's opt-in twice over.** Per-call via `--as-json`, or set-once via a
  truthy `MLR_HELP_JSON` environment variable.

For routing an *intent* to a capability -- the reverse of browsing -- `mlr
which` returns ranked candidates:

<pre class="pre-highlight-in-pair">
<b>mlr which "join two files on a key" | mlr --json head -n 2</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
{
  "kind": "verb",
  "name": "join",
  "score": 25,
  "summary": "Joins records from specified left file name with records from all file names"
},
{
  "kind": "function",
  "name": "joink",
  "score": 25,
  "summary": "Makes string from map/array keys. First argument is map/array; second is separator string."
}
]
</pre>

Its exit code signals confidence -- 0 when a query word matched a
capability's name, 2 when it didn't -- so a harness can branch on status
without parsing anything.

## Constrain: the tool's shape, and the data's shape

Agents don't just hallucinate flags; they hallucinate *values*. Miller
attacks that from both sides.

Where an option's domain is fixed by the binary, the catalog says so:
`type` is `enum` and `values` is the complete list. Here's one option of the
[summary](reference-verbs.md#summary) verb, extracted from the catalog --
using Miller to query Miller:

<pre class="pre-highlight-in-pair">
<b>mlr help verb summary --as-json | mlr --json put -q 'emit $options[1]'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
{
  "flag": "-a",
  "arg": "{mean,sum,etc.}",
  "type": "enum",
  "desc": "Use only the specified summarizers.",
  "values": ["field_type", "count", "null_count", "distinct_count", "mode", "sum", "mean", "stddev", "var", "skewness", "minlen", "maxlen", "min", "p25", "median", "p75", "max", "iqr", "lof", "lif", "uif", "uof"]
}
]
</pre>

Where the domain depends on your *data* -- which fields exist, what values
`filter` could compare against, what to pass to `-g` -- the
[describe](reference-verbs.md#describe) verb profiles the input in one pass:
per field, the types seen, counts, cardinality, min/max, and (for
low-cardinality fields) every distinct value:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson describe then head -n 2 example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
{
  "field_name": "color",
  "types": {
    "string": 10
  },
  "count": 10,
  "null_count": 0,
  "distinct_count": 3,
  "min": "purple",
  "max": "yellow",
  "values": ["yellow", "red", "purple"]
},
{
  "field_name": "shape",
  "types": {
    "string": 10
  },
  "count": 10,
  "null_count": 0,
  "distinct_count": 3,
  "min": "circle",
  "max": "triangle",
  "values": ["triangle", "square", "circle"]
}
]
</pre>

The catalog is the *tool's* shape; `describe` is the *data's* shape. An
agent that consults both has nothing left to guess.

## Validate: check DSL before spending a run

`mlr put --explain` (likewise `mlr filter --explain`) parses and type-checks
an expression, then exits -- without opening any input at all:

<pre class="pre-highlight-in-pair">
<b>mlr put --explain '$z = $x + $y'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
mlr put: DSL expression is valid.
</pre>

## Run and recover: errors as data

With `--errors-json` (or a truthy `MLR_ERRORS_JSON` environment variable),
errors arrive as a structured document. The `kind` field gives an agent
something to branch on; `hint` is a runnable next step, not a sentence; and
`did_you_mean` is computed against the same catalog the agent discovered
from, closing the self-correction loop:

<pre class="pre-highlight-in-pair">
<b>mlr --errors-json --icsv sortt -f shape example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "error": "mlr: verb \"sortt\" not found. Please use \"mlr -l\" for a list.",
  "kind": "unknown-verb",
  "token": "sortt",
  "hint": "Run 'mlr -l' for a list of verbs, or 'mlr help verb \u003cname\u003e' for details.",
  "did_you_mean": [
    "sort"
  ]
}
</pre>

And since Miller's DSL includes [system and exec](shell-commands.md), there's
a sandbox: `--no-shell` (or a truthy `MLR_NO_SHELL` environment variable)
disables all external-command execution -- the DSL `system` and `exec`
functions, piped redirects, and `--prepipe` fail cleanly:

<pre class="pre-highlight-in-pair">
<b>mlr --no-shell -n put 'end{print system("hostname")}'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
(error)
</pre>

A typical agent profile sets all three environment variables once:

<pre class="pre-non-highlight-non-pair">
export MLR_HELP_JSON=1    # help/catalog output as JSON
export MLR_ERRORS_JSON=1  # errors as structured JSON
export MLR_NO_SHELL=1     # no external-command execution
</pre>

Put together, the sections above are a loop -- discover, constrain,
validate, run -- where each step feeds the next and failures route back with
structure instead of prose.

## Plug it in: the MCP server

If your agent speaks the [Model Context
Protocol](https://modelcontextprotocol.io) -- Claude Code, Claude Desktop,
Cursor, and many others -- everything above is one line away. For Claude
Code:

<pre class="pre-highlight-in-pair">
<b>claude mcp add miller -- mlr mcp</b>
</pre>

That's the whole setup. The server's five tools are exactly the features on
this page -- `list_capabilities` and `which` for discovery, `describe_data`
to constrain, `validate_dsl` to validate, and `run` (sandboxed with
`--no-shell` by default) to execute -- plus a shipped playbook, as MCP prompt
and resource, teaching the agent the loop. Then just talk to your agent
about your data:

* "Which fields in `data.csv` have missing values?"
* "Convert this CSV to JSON, keeping only rows where status is active."
* "Join `a.csv` and `b.csv` on id, and give me the mean rate per group."

See [The MCP server](mcp-server.md) for the full tool reference and server
options.

## Before and after: a first session with the skill installed

Everything above assumes you already know why each piece matters. If you're new to Miller, or
you've used Miller before but this is your first time on 6.20 or newer, here's the same idea as one
worked example: install the skill, then watch what changes about talking to your AI assistant.

One thing to be clear on before the example: you never type `mlr` yourself in this section. You
type plain English to your agent, same as always. Every `mlr` command shown below is the agent's
*own* work -- what it runs on your behalf, in the background, to answer you. They're printed here
so you can see exactly what changes, not because you'd type them.

### Install the skill

One command. It writes a single file; nothing runs in the background, and nothing here makes a
network call:

<pre class="pre-highlight-in-pair">
<b>$ mlr skill install ~/.claude/skills/miller</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Wrote /Users/kerl/.claude/skills/miller/SKILL.md
</pre>

For Claude Code, that's the whole setup -- it picks up skills under `~/.claude/skills/`
automatically. Other agents that read Agent Skills from disk look in their own equivalent
directory; check your agent's docs for where that is.

The file itself is short -- it's instructions, not code, so it's worth reading once:

<pre class="pre-highlight-in-pair">
<b>mlr skill print | head -n 8</b>
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
</pre>

### Before: an agent guessing at your data

Say you're looking at [example.csv](example.csv) for the first time. You type this, and nothing
else, to your AI assistant:

> **You:** In example.csv, show me the red rows.

Without the skill, a reasonable-sounding guess for the DSL is `$color == "Red"` -- English
capitalizes color names, so why not? Here's that guess, run exactly as the agent would run it,
behind the scenes, on your machine:

<pre class="pre-highlight-non-pair">
<b>mlr --csv filter '$color == "Red"' example.csv</b>
</pre>

Nothing comes back -- no error, no rows, no warning, exit code 0. And here's the trap: the agent
still owes you an answer, so it turns that silence into a sentence:

> **Agent:** I checked example.csv and there aren't any rows where color is red.

That's wrong -- there are four -- but it *reads* like a fact, because a wrong guess about *your
data*, unlike a wrong flag or function name, doesn't look like a failure on Miller's end. It looks
like an empty result, which could just as easily have been true.

### After: an agent that checks first

Same question, word for word, with the skill installed:

> **You:** In example.csv, show me the red rows.

The skill's playbook puts a step between your question and any guesswork: "constrain -- learn the
data before touching it." So before writing any comparison, the agent runs a `describe`, again
invisibly, on your behalf:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson describe example.csv | mlr --json filter '$field_name == "color"'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
{
  "field_name": "color",
  "types": {
    "string": 10
  },
  "count": 10,
  "null_count": 0,
  "distinct_count": 3,
  "min": "purple",
  "max": "yellow",
  "values": ["yellow", "red", "purple"]
}
]
</pre>

That's the real value set for `color`: lowercase, three values, no `"Red"` anywhere. The agent
corrects its own comparison before running anything user-facing:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint filter '$color == "red"' example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color shape  flag  k index quantity rate
red   square true  2 15    79.2778  0.0130
red   circle true  3 16    13.8103  2.9010
red   square false 4 48    77.5542  7.4670
red   square false 6 64    77.1991  9.5310
</pre>

-- and only now replies to you:

> **Agent:** Four rows have color = red: rows 2, 3, 4, and 6 (all square or circle, quantity
> ranging from about 13.8 to 79.3).

Same question, same data, same underlying `mlr` binary -- the only thing that changed is that the
agent looked before it leapt, and you never saw the intermediate `describe` unless you asked to.
That one habit, *check the data before writing a comparison*, is the skill in miniature; the rest
of the playbook applies the same idea to verb and function names (discover), DSL syntax (validate),
and error messages (recover) -- see the sections above for each.

### Do you need the MCP server too?

Not to get this benefit. The skill alone gets an agent through the whole loop using plain `mlr`
commands over whatever shell-executing tool your agent already has -- nothing above required MCP.
If your agent also speaks [MCP](https://modelcontextprotocol.io), registering the [MCP
server](mcp-server.md) upgrades the *mechanism*: `describe_data`, `validate_dsl`, and the rest
become typed tool calls returning structured JSON, instead of the agent reading `mlr`'s text output
itself. But the *loop* -- discover, constrain, validate, run -- is identical either way. Day one,
the skill alone is enough; add the MCP server later if your agent supports it and you want the
sturdier plumbing.
