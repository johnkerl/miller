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
# Miller and AI agents

Miller isn't just friendly at the command line -- it's built to be driven by
AI agents. Everything an agent needs to use Miller *well* is a first-class
feature of the tool itself:

* Miller's entire surface -- verbs, DSL functions, flags, keywords -- is
  available as **machine-readable JSON**, so agents ground themselves in what
  actually exists instead of hallucinating flags.
* **Errors are structured** -- kind, hint, did-you-mean -- so agents branch on
  data rather than parsing English.
* DSL expressions can be **validated before running**, without reading any
  input.
* Input data can be **profiled in one pass** -- field names, types,
  cardinality, value domains -- so agents copy real names and values instead
  of guessing.
* A built-in **MCP server** packages all of the above, plus a playbook that
  teaches the agent how to use it.

## Getting started: the MCP fast path

If your agent speaks the [Model Context
Protocol](https://modelcontextprotocol.io) -- Claude Code, Claude Desktop,
Cursor, and many others -- one line connects it to Miller. For Claude Code:

<pre class="pre-highlight-in-pair">
<b>claude mcp add miller -- mlr mcp</b>
</pre>

That's the whole setup. Your agent now has five tools --
`list_capabilities`, `which`, `validate_dsl`, `describe_data`, and `run` --
plus a playbook teaching it the discover &rarr; constrain &rarr; validate
&rarr; run loop. Commands the agent runs are sandboxed against
external-command execution (see [AI agents and the MCP
server](mcp-server.md) about `--no-shell`).

Then just talk to your agent about your data:

* "Which fields in `data.csv` have missing values?"
* "Convert this CSV to JSON, keeping only rows where status is active."
* "Join `a.csv` and `b.csv` on id, and give me the mean rate per group."

See [AI agents and the MCP server](mcp-server.md) for the full tool
reference and server options.

## Getting started without MCP

Any agent that can run shell commands can use the same surface directly --
these are all ordinary Miller features. If you're writing a system prompt, an
agent skill, or a script harness, these are the pieces to teach it. (The
ready-made version of that lesson ships in the Miller source tree as
`pkg/terminals/mcp/SKILL.md`, in Agent Skill format.)

### Route an intent to the right tool

`mlr which` turns "what I want" into ranked candidates -- here trimmed to the
top two using Miller itself:

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

The exit code signals confidence: 0 when a query word matched a capability's
name, 2 when it didn't -- so a script can branch without reading the output.

### Read the catalog, not the prose

Every help topic has an `--as-json` form:

<pre class="pre-highlight-in-pair">
<b>mlr help function splitax --as-json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
  {
    "name": "splitax",
    "class": "conversion",
    "arity": "2",
    "help": "Splits string into array without type inference. First argument is string to split; second is the separator to split on.",
    "examples": [
      "splitax(\"3,4,5\", \",\") = [\"3\",\"4\",\"5\"]"
    ]
  }
]
</pre>

`mlr help --as-json` emits the entire catalog in one document, and `mlr help
--as-json --index` gives just names and one-line summaries -- the cheap first
call. The document carries `mlr_version` and `catalog_schema_version`, which
together make a perfect cache key. Setting the `MLR_HELP_JSON` environment
variable to a truthy value makes all help output JSON without the flag.

### Learn the data before writing a command

The [describe](reference-verbs.md#describe) verb profiles input in one pass:
types, counts, cardinality, and -- for low-cardinality fields -- every
distinct value, so an agent can copy real values into `-g` flags and DSL
comparisons:

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

### Validate DSL before spending a run

`mlr put --explain` (likewise `mlr filter --explain`) parses and type-checks
an expression and exits without reading any input:

<pre class="pre-highlight-in-pair">
<b>mlr put --explain '$z = $x + $y'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
mlr put: DSL expression is valid.
</pre>

### Get errors as data, not prose

With `--errors-json` (or a truthy `MLR_ERRORS_JSON` environment variable),
errors arrive as a structured document -- and `did_you_mean` closes the
self-correction loop:

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

### Keep the agent in a sandbox

`--no-shell` (or a truthy `MLR_NO_SHELL` environment variable) disables
Miller's ability to run external commands -- the DSL `system` and `exec`
functions, piped redirects, and `--prepipe` all fail cleanly -- so an
agent-constructed command line can't shell out. A typical agent profile sets
all three environment variables once:

<pre class="pre-non-highlight-non-pair">
export MLR_HELP_JSON=1    # help/catalog output as JSON
export MLR_ERRORS_JSON=1  # errors as structured JSON
export MLR_NO_SHELL=1     # no external-command execution
</pre>

## The loop

Whichever path you use, effective agents follow the same four steps -- each
one exists to prevent a specific failure mode:

1. **Discover** capabilities from the catalog (`which`, `mlr help --as-json`)
   -- never invent flag or function names.
2. **Constrain** to the data's actual shape (`describe`) -- copy field names
   and values, don't guess them.
3. **Validate** DSL expressions before running them (`--explain`).
4. **Run**, and on failure branch on the structured error's `kind`, `hint`,
   and `did_you_mean` rather than re-guessing.
