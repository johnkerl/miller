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
# The Miller MCP server

Miller ships with a built-in [Model Context Protocol](https://modelcontextprotocol.io)
server, so AI agents (Claude Code, Claude Desktop, Cursor, and other MCP
clients) can discover and drive Miller without scraping help text or guessing
at flags. (For the overview of Miller's whole AI feature set -- with or
without MCP -- see [Miller and AI](ai.md).)

The server speaks JSON-RPC over stdin/stdout (MCP's **stdio** transport): the MCP client spawns `mlr
mcp` as a subprocess. No network port is opened, and the server exits when the client disconnects.
Example registration, for Claude Code:

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

## Transparency: what actually happens

The tool descriptions above are accurate but abstract. Concretely, here is
the full trace of what happens on your machine when an agent with the
Miller MCP server registered decides to answer a request like "describe the
file `test/input/abixy`":

The agent does not run a shell command itself. It calls the `describe_data`
tool with a small JSON argument, typically just `{"files":
["test/input/abixy"]}` -- no `input_format`, `max_values`, or `extra_args`,
since none were implied by the request. The `describeDataHandler` in Miller's
source turns that into an argv: `--ojson describe test/input/abixy`. There's
no `-i` flag because `input_format` was empty, so Miller falls back to its
own default input format (DKVP); there's no `-n` because `max_values` was
zero, so the `describe` verb's own default of 20 distinct values applies.

That argv is run as a subprocess of the *same* `mlr` binary that is serving
the MCP session (Go's `os.Executable()`), not whatever `mlr` happens to be
first on `$PATH` -- so behavior is pinned to the version you registered,
even if another Miller build is installed elsewhere. Two environment
variables are layered onto the subprocess's otherwise-inherited environment:
`MLR_ERRORS_JSON=1` always, so a failure comes back as a structured JSON
document on stderr instead of prose that the agent would have to
guess-parse; and `MLR_NO_SHELL=1` unless the server was started with
`--allow-shell`, disabling `system`/`exec` and piped redirects in the child
(see [Sandboxing](#sandboxing---no-shell) above). Every other environment
variable in the server's own environment is passed through unchanged. The
subprocess's stdin is always a fresh empty reader -- it never sees the MCP
session's own stdin, which is busy carrying the JSON-RPC transport. The call
is wrapped in a wall-clock timeout, 60 seconds by default (`mlr mcp
--timeout`); `run` lets an agent override this per call, but `describe_data`
and `validate_dsl` always use the server's default. Captured stdout and
stderr are each capped at 1 MiB by default (`mlr mcp --max-output-bytes`);
past that, further output is silently dropped and the result is flagged
`truncated` rather than the subprocess being killed.

If the subprocess exits 0, the handler parses its captured stdout as JSON --
an array of per-field objects with name, inferred type, count, cardinality,
min/max, and (for low-cardinality fields) the distinct values themselves --
and hands that structured data back to the agent as the tool's result, not
as raw text. The agent then composes whatever natural-language answer it
gives you from that structured data, not from re-reading Miller's terminal
output. The same shape of trace -- typed input, an argv built from it,
subprocess exec against the registered `mlr` binary, capped and
timeout-bounded I/O, structured output back to the agent -- applies to
`validate_dsl` and `run` as well; `run` is the one tool where the argv is
whatever the agent supplies rather than one Miller assembles for it.

The short version: your words never touch a shell. The agent calls a typed
tool; Miller turns that into an exact `mlr` command line, runs it as a
subprocess of the very binary you registered, with `MLR_ERRORS_JSON=1` and
(by default) `MLR_NO_SHELL=1` layered onto your inherited environment; and
hands back structured data rather than terminal text.

## The agent playbook

The server also exposes a playbook -- as MCP prompt `miller-playbook` and MCP
resource `miller://playbook` -- encoding the loop that makes an agent
effective with Miller: **discover** capabilities from the catalog rather than
inventing them, **constrain** to the data's actual fields and values via
`describe_data`, **validate** DSL before running it, and **run** with
structured-error recovery. The same text lives in the Miller source tree at
[pkg/terminals/mcp/SKILL.md](https://github.com/johnkerl/miller/blob/main/pkg/terminals/mcp/SKILL.md)
in Agent Skill format.
