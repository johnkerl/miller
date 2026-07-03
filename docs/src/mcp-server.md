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
# AI agents and the MCP server

Miller ships with a built-in [Model Context Protocol](https://modelcontextprotocol.io)
server, so AI agents (Claude Code, Claude Desktop, Cursor, and other MCP
clients) can discover and drive Miller without scraping help text or guessing
at flags. (For the quick-start overview of Miller's whole AI-agent feature
set -- with or without MCP -- see [Miller and AI agents](ai-agents.md).)

The server speaks JSON-RPC over stdin/stdout (MCP's "stdio" transport): the
MCP client spawns `mlr mcp` as a subprocess. No network port is opened, and
the server exits when the client disconnects. Example registration, for
Claude Code:

<pre class="pre-highlight-in-pair">
<b>claude mcp add miller -- mlr mcp</b>
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr mcp --help</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Usage: mlr mcp [options]
Runs a Model Context Protocol (MCP) server exposing Miller to AI agents.

The server speaks JSON-RPC over stdin/stdout (MCP "stdio" transport); it is
meant to be spawned by an MCP client rather than run interactively. Example
registration, for Claude Code:

  claude mcp add miller -- mlr mcp

Tools exposed:
  list_capabilities  the mlr help --as-json catalog/index, filterable
  which              intent -> ranked matching verbs/functions/flags/keywords
  validate_dsl       parse/type-check a put/filter expression without running it
  describe_data      field names/types/cardinality/values of input data
  run                run an mlr command line

Also exposed: an agent playbook, as MCP prompt "miller-playbook" and MCP
resource "miller://playbook", encoding the discover -> constrain -> validate
-> run loop.

Commands started by the run/validate_dsl/describe_data tools are run with
MLR_NO_SHELL=1 -- the DSL system/exec functions, piped redirects, and
--prepipe/--prepipex fail cleanly -- and with MLR_ERRORS_JSON=1 so errors
come back structured.

Options:
 --allow-shell           Do not set MLR_NO_SHELL=1 on subprocesses, re-enabling
                         the DSL system/exec functions for agent-run commands.
 --timeout {seconds}     Default wall-clock limit for the run tool (default 60).
 --max-output-bytes {n}  Cap captured stdout/stderr per command (default 1048576).
 -h or --help            Show this message.
</pre>

## What the tools map to

Each MCP tool is a thin wrapper over a Miller feature you can also use
directly from the command line:

* `list_capabilities` is [`mlr help --as-json`](online-help.md) -- the
  machine-readable catalog of verbs, DSL functions, flags, and keywords.
* `which` is `mlr which` -- natural-language intent to ranked capabilities.
* `validate_dsl` is `mlr put --explain` / `mlr filter --explain` -- parse and
  type-check a DSL expression without reading any input.
* `describe_data` is [`mlr describe`](reference-verbs.md#describe) -- field
  names, types, cardinality, and value domains for input data.
* `run` executes an `mlr` command line and reports exit code, output, and --
  on failure -- the structured error document from `mlr --errors-json`.

The catalog tools are answered in-process; the others run this same `mlr`
binary as a subprocess, so agents see exactly what a terminal user sees.

## Sandboxing: --no-shell

Miller's DSL includes [`system` and `exec`](shell-commands.md), and
`--prepipe`/piped redirects also run external commands. So that an
agent-constructed command line doesn't imply arbitrary command execution,
subprocesses started by the MCP server run with `MLR_NO_SHELL=1`: those
features fail cleanly instead of executing. Start the server with
`mlr mcp --allow-shell` to turn that off.

The same gate is available outside the MCP server: pass `--no-shell` to any
`mlr` invocation, or set the `MLR_NO_SHELL` environment variable to a truthy
value. Note that Miller can still write files when asked to (`tee`, `split`,
DSL output redirection) -- the gate is specifically about executing external
commands.

## The agent playbook

The server also exposes a playbook -- as MCP prompt `miller-playbook` and MCP
resource `miller://playbook` -- encoding the loop that makes an agent
effective with Miller: **discover** capabilities from the catalog rather than
inventing them, **constrain** to the data's actual fields and values via
`describe_data`, **validate** DSL before running it, and **run** with
structured-error recovery. The same text lives in the Miller source tree at
`pkg/terminals/mcp/SKILL.md` in Agent Skill format.
