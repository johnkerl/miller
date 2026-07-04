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
# Miller and AI

Miller treats AI agents as first-class users. When an agent drives a
command-line tool, it fails in predictable ways: it invents flags that don't
exist, guesses values that aren't in the data, misreads error prose, and
burns whole runs discovering a typo. Miller closes off each of those failure
modes with structure:

* Miller's entire surface -- verbs, DSL functions, flags, keywords -- is
  available as **machine-readable JSON**, so agents ground themselves in what
  actually exists.
* Options with fixed domains carry their **complete value sets**, and input
  data can be **profiled in one pass** -- so agents copy real values instead
  of inventing them.
* DSL expressions can be **validated before running**, without reading any
  input.
* **Errors are structured** -- kind, hint, did-you-mean -- so agents branch
  on data rather than parsing English.
* A **sandbox flag** removes external-command execution, so an
  agent-constructed command line is just data processing.

Everything on this page is an ordinary command-line feature: it works from
any agent harness, system prompt, or script -- and it's equally useful for
plain shell tooling like `jq`. The [MCP server](#plug-it-in-the-mcp-server)
at the end packages it all up for MCP-speaking agents.

## The bare minimum

**To get the AI features:** install Miller 6.20 or newer ([Installing
Miller](installing-miller.md)). That's all. Everything on this page ships
inside the ordinary `mlr` binary -- there are no plugins, no separate
installs, no API keys, and nothing here makes network calls.

**To get your AI to use them,** pick whichever matches your setup:

* **If your agent speaks MCP** (Claude Code, Claude Desktop, Cursor, ...):
  register the server -- for Claude Code that's `claude mcp add miller -- mlr
  mcp` -- and you're done. The tools describe themselves, and the server
  ships its own instructions and playbook, so you usually don't need to say
  anything special; if the agent doesn't reach for them, a nudge like "use
  the Miller tools" suffices. Details in [The MCP server](mcp-server.md).

* **If your agent just runs shell commands** (a system prompt, a
  `CLAUDE.md`, Cursor rules, a script harness): paste this standing
  instruction into its context:

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
